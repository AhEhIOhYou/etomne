package entities

import (
	"html"
	"strings"
	"time"

	"github.com/AhEhIOhYou/etomne/backend/constants"
)

type Model struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ModelRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ModelData struct {
	Model Model       `json:"model"`
	User  PublicUser  `json:"author"`
	Files SortedFiles `json:"files"`
}

func (modelReq *ModelRequest) NewModel() *Model {
	return &Model{
		Title:       modelReq.Title,
		Description: modelReq.Description,
	}
}

func (m *Model) BeforeUpdate() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
	m.Description = html.EscapeString(strings.TrimSpace(m.Description))
	m.UpdatedAt = time.Now()
}

func (m *Model) Prepare() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
	m.Description = html.EscapeString(strings.TrimSpace(m.Description))
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
}

func (model *ModelRequest) Validate() string {
	if model.Title == "" {
		return constants.ModelTitleCantBeEmpty
	}
	return ""
}

func (model *Model) Validate() string {
	if model.Title == "" {
		return constants.ModelTitleCantBeEmpty
	}
	if model.UserID == 0 {
		return constants.UserIDInvalid
	}
	return ""
}
