package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	Id       primitive.ObjectID `json:"id,"`
	Name     string             `json:"name"`
	Email    string             `json:"email" validate:"required"`
	Password string             `json:"password" validate:"required"`
}

// func (user *User) SetPassword(password string) {
// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
// 	user.Password = hashedPassword
// }
// func (user *User) ComparePassword(password string) error {
// 	return bcrypt.CompareHashAndPassword(user.Password, []byte(password))
// }

func SetPassword2(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hashedPassword)
}

func ComparePassword2(password string, passwordHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
}
