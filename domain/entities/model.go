package entities

import (
	"html"
	"strings"
)

type Model struct {
	Id          uint64 `json:"id"`
	Title       string `json:"name"`
	CreateDate  string `json:"create_date"`
	Description string `json:"description"`
	File        string `json:"file_model"`
}

func (m *Model) BeforeSave() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
}

func (m *Model) Prepare() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
	m.Description = html.EscapeString(strings.TrimSpace(m.Description))
}

func (m *Model) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if m.Title == "" || m.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if m.Description == "" || m.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	default:
		if m.Title == "" || m.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if m.Description == "" || m.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	}
	return errorMessages
}
