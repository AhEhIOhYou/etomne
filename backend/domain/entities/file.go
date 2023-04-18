package entities

import (
	"html"
	"strings"
	"time"
)

type File struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	OwnerId   uint64    `json:"owner_id"`
	Url       string    `json:"url"`
	Extension string    `json:"extension"`
	CreatedAt time.Time `json:"created_at"`
}

type ModelFile struct {
	FileId  uint64 `gorm:"not null" json:"file_id"`
	ModelId uint64 `gorm:"not null" json:"model_id"`
}

type UserPhoto struct {
	FileId uint64 `gorm:"not null" json:"file_id"`
	UserId uint64 `gorm:"not null" json:"user_id"`
	Size   uint64 `gorm:"not null" json:"size"`
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
