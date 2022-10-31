package entities

import (
	"html"
	"strings"
	"time"
)

type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	AuthorId  uint64    `gorm:"size:100;not null;" json:"author_id"`
	ModelId   uint64    `gorm:"size:100;not null;" json:"model_id"`
	Message   string    `gorm:"text;notnull" json:"message"`
	ReplyId   uint64    `gorm:"size:100" json:"reply_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (c *Comment) BeforeUpdate() {
	c.Message = html.EscapeString(strings.TrimSpace(c.Message))
	c.UpdatedAt = time.Now()
}

func (c *Comment) Prepare() {
	c.Message = html.EscapeString(strings.TrimSpace(c.Message))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

func (c *Comment) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if c.Message == "" || c.Message == "null" {
			errorMessages["error"] = "message is required"
		}
	default:
		if c.Message == "" || c.Message == "null" {
			errorMessages["error"] = "message is required"
		}
	}
	return errorMessages
}
