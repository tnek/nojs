package model

import (
	"encoding/base64"
	"errors"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrorUserExists = errors.New("user already exists")
)

func UserExists(username string) (bool, error) {
	return userExists(db, username)
}

func userExists(tx *gorm.DB, username string) (bool, error) {
	var count int64
	if result := tx.Model(&User{}).Where("name = ?", username).Count(&count); result.Error != nil {
		return true, result.Error
	}
	return count > 0, nil
}

func NewUser(username string, password string, isAdmin bool) (string, error) {
	if username == "" {
		return "", errors.New("username can't be empty")
	}

	if password == "" {
		return "", errors.New("password can't be empty")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	var u *User

	if err := db.Transaction(func(tx *gorm.DB) error {
		exists, err := userExists(tx, username)
		if err != nil {
			return err
		}

		if exists {
			return ErrorUserExists
		}
		u = &User{
			ID:           uuid.New().String(),
			Name:         username,
			PasswordHash: base64.RawURLEncoding.EncodeToString(passwordHash),
			IsAdmin:      isAdmin,
		}

		return tx.Model(&User{}).Omit(clause.Associations).Create(u).Error

	}); err != nil {
		return "", err
	}

	log.Printf("User created: %v\n", u)
	return u.ID, nil
}

func userByUUID(tx *gorm.DB, uuid string) (*User, error) {
	var user User

	if result := tx.Model(&User{}).Where("id = ?", uuid).First(&user); result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func UserByUUID(uuid string) (*User, error) {
	return userByUUID(db, uuid)
}

func GetUser(username string) (*User, error) {
	var user User
	result := db.Model(&User{}).Where("name = ?", username).First(&user)
	return &user, result.Error
}

func Login(username string, pw string) (string, error) {
	var user User

	result := db.Model(&User{}).Where("name = ?", username).First(&user)
	if result == nil {
		return "", errors.New("unknown user")
	}

	if result.Error != nil {
		return "", result.Error
	}

	// Check password
	hash, err := base64.RawURLEncoding.DecodeString(user.PasswordHash)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)); err != nil {
		return "", errors.New("incorrect password")
	}
	return user.ID, nil
}
