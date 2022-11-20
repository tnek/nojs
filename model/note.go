package model

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrorUnauthorized = errors.New("Unauthorized")
)

const (
	kUntitledNote         = "Untitled note"
	kUntitledNoteContents = "Write something..."
)

func ShareNote(from *User, noteID string, toName string) (*User, error) {
	var to User
	if err := db.Transaction(func(tx *gorm.DB) error {
		var note Note

		if err := tx.Model(&Note{}).Where("id = ? and author_id = ?", noteID, from.ID).First(&note).Error; err != nil {
			log.Printf("note doesn't exist %v as user %v: %w", noteID, from.Name, err)
			return fmt.Errorf("note id %v doesn't exist or you don't have permission", noteID)
		}

		if err := tx.Model(&User{}).Where("name = ?", toName).First(&to).Error; err != nil {
			log.Printf("share failed to fetch existing user %v: %w", toName, err)
			return fmt.Errorf("user %v doesn't exist", toName)
		}

		if err := tx.Model(&to).Association("ViewableNotes").Append(&note); err != nil {
			return fmt.Errorf("share failed to update association: %w", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &to, nil
}

func FetchNotes(u *User) ([]Note, error) {
	var notes []Note
	if err := db.Model(u).Preload("Author").Association("ViewableNotes").Find(&notes); err != nil {
		return nil, err
	}
	sort.SliceStable(notes, func(i, j int) bool {
		return notes[i].UpdatedAt.After(notes[j].UpdatedAt)
	})
	return notes, nil
}

func NewNote(u *User, title string, contents string) (string, error) {
	log.Printf("Creating note: %v %v", u.Name, title)
	noteID := uuid.New().String()
	if title == "" {
		title = kUntitledNote
	}
	if contents == "" {
		contents = kUntitledNoteContents
	}

	note := &Note{
		ID:       noteID,
		Title:    title,
		Contents: contents,
	}
	if err := db.Model(u).Association("ViewableNotes").Append(note); err != nil {
		return "", fmt.Errorf("failed to associate note: %w", err)
	}
	if err := db.Model(note).Association("Author").Replace(u); err != nil {
		return "", fmt.Errorf("failed to associate author: %w", err)
	}
	return noteID, nil
}

func DeleteNote(u *User, noteID string) (bool, error) {
	log.Printf("%v requested deleteNote id %v\n", u.Name, noteID)
	if err := db.Transaction(func(tx *gorm.DB) error {
		var note Note

		if err := tx.Model(&Note{}).Where("id = ? and author_id = ?", noteID, u.ID).First(&note).Error; err != nil {
			log.Printf("note doesn't exist %v as user %v: %w", noteID, u.Name, err)
			return fmt.Errorf("note id %v doesn't exist or you don't have permission", noteID)
		}
		if err := tx.Unscoped().Delete(&Note{}, "id = ? and author_id = ?", noteID, u.ID).Error; err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}
		return nil
	}); err != nil {
		return false, err
	}
	return true, nil
}
