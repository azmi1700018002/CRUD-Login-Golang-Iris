package handlers

import (
	"crud-golang-iris/app/repositories"
	"crud-golang-iris/domain"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	Repo repositories.UserRepository
}

func (h *UserHandler) GenerateToken(user *domain.User) (string, error) {
	// Buat payload JWT dengan data yang relevan
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		// Anda dapat menambahkan data lain ke dalam payload sesuai kebutuhan
	}

	// Buat token dengan signing key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("golang-iris"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (h *UserHandler) AuthenticateToken(ctx iris.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString("Missing token")
		return
	}

	// Verifikasi token
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// Menggunakan signing key yang sama seperti saat pembuatan token
		return []byte("golang-iris"), nil
	})

	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString("Invalid token")
		return
	}

	// Token valid, lanjutkan ke handler berikutnya
	ctx.Next()
}

func (h *UserHandler) Login(ctx iris.Context) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := ctx.ReadJSON(&credentials)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	user, err := h.Repo.GetUserByUsername(credentials.Username)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString("Invalid username or password")
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.WriteString("Invalid username or password")
		return
	}

	// Jika login berhasil, generate token
	token, err := h.GenerateToken(user)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("Failed to generate token")
		return
	}

	// Kirim token ke klien
	ctx.JSON(iris.Map{"token": token})
}

// Handler untuk mendapatkan daftar semua pengguna
func (h *UserHandler) GetUsers(ctx iris.Context) {
	users, err := h.Repo.GetUsers()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.JSON(users)
}

// Handler untuk membuat pengguna baru
func (h *UserHandler) CreateUser(ctx iris.Context) {
	var user domain.User
	err := ctx.ReadJSON(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	err = h.Repo.CreateUser(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(user)
}

// Handler untuk mendapatkan pengguna berdasarkan ID
func (h *UserHandler) GetUserID(ctx iris.Context) {

	id, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	user, err := h.Repo.GetUserByID(id)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString(err.Error())
		return
	}

	ctx.JSON(user)

}

// Handler untuk memperbarui pengguna berdasarkan ID
func (h *UserHandler) UpdateUser(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	var user domain.User
	err = ctx.ReadJSON(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	// Hash password sebelum menyimpannya
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}
	user.Password = string(hashedPassword)

	err = h.Repo.UpdateUser(id, &user)
	if err != nil {
		ctx.StatusCode(iris.StatusNotFound)
		ctx.WriteString(err.Error())
		return
	}

	ctx.JSON(user)
}

// Handler untuk menghapus pengguna berdasarkan ID
func (h *UserHandler) DeleteUser(ctx iris.Context) {
	id, err := ctx.Params().GetInt64("id")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	err = h.Repo.DeleteUser(id)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	ctx.WriteString("User deleted")
}
