package domain

type User struct {
	ID       int64 `gorm:"primaryKey"`
	Username string
	Password string `gorm:"not null"`
	Name     string `gorm:"not null"`
	Age      int    `gorm:"not null"`
}
