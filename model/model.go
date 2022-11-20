package model

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func Conn(dbName string) error {
	if db != nil {
		return nil
	}

	if dbName == "" {
		log.Println("No database specified, defaulting to /tmp/test.db")
		dbName = "/tmp/test.db"
	}

	var err error
	db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&User{}, &Note{}); err != nil {
		return err
	}
	return nil
}

type User struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	PasswordHash string
	Name         string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	IsAdmin bool

	AvatarURL string
	Bio       string

	ViewableNotes []Note `gorm:"many2many:viewable_note"`
}

type Note struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	Title    string
	Contents string
	AuthorID string

	CreatedAt time.Time
	UpdatedAt time.Time

	Author  User   `gorm:"foreignKey:AuthorID;references:ID"`
	Viewers []User `gorm:"many2many:viewable_note"`
}
