package entities

import (
	"html"
	"strings"
	"time"
)

type File struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:100;not null;" json:"title"`
	OwnerId   uint64    `gorm:"size:100;not null" json:"owner_id"`
	Url       string    `gorm:"size:255;not null;" json:"url"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

type ModelFile struct {
	FileId  uint64 `gorm:"not null" json:"file_id"`
	ModelId uint64 `gorm:"not null" json:"model_id"`
}

func (f *File) Prepare() {
	f.Title = html.EscapeString(strings.TrimSpace(f.Title))
	f.CreatedAt = time.Now()
}

type Files []File

func (f *File) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if f.Title == "" || f.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
	default:
		if f.Title == "" || f.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
	}
	return errorMessages
}
