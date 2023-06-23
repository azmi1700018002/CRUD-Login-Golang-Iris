package main

import (
	"crud-golang-iris/app/handlers"
	"crud-golang-iris/app/repositories"
	"crud-golang-iris/infrastructure/database"
	"log"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
)

func main() {
	dsn := "host=localhost user=postgres password=postgres dbname=golang_iris port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := database.NewGormDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = database.AutoMigrate(db)
	if err != nil {
		log.Fatal(err)
	}

	app := iris.New()
	app.Use(recover.New())
	app.Use(logger.New())

	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.UserHandler{Repo: *userRepo}

	app.Post("/login", userHandler.Login)
	app.Post("/users", userHandler.CreateUser)
	// Menggunakan middleware untuk autentikasi token pada grup rute "/users"
	users := app.Party("/secure")
	users.Use(userHandler.AuthenticateToken)
	{
		users.Get("/users", userHandler.GetUsers)
		users.Post("/users", userHandler.CreateUser)
		users.Get("/{id:int}", userHandler.GetUserID)
		users.Put("/{id:int}", userHandler.UpdateUser)
		users.Delete("/{id:int}", userHandler.DeleteUser)
	}

	err = app.Run(iris.Addr(":8080"))
	if err != nil {
		log.Fatal(err)
	}
}
