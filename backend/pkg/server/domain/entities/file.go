package entities

import (
	"html"
	"strings"
	"time"

	"github.com/AhEhIOhYou/etomne/pkg/server/constants"
)

type File struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	OwnerId   uint64    `json:"owner_id"`
	Url       string    `json:"url"`
	Extension string    `json:"extension"`
	CreatedAt time.Time `json:"created_at"`
}

type Files []File

type FileRequest struct {
	Title   string `json:"title"`
	OwnerId uint64 `json:"owner_id"`
}

type ModelFile struct {
	FileId  uint64 `json:"file_id"`
	ModelId uint64 `json:"model_id"`
}

type SortedFiles struct {
	GLB   Files `json:"glb"`
	IMG   Files `json:"img"`
	Video Files `json:"video"`
	Other Files `json:"other"`
}

func (fileReq *FileRequest) NewFile() *File {
	return &File{
		Title:   fileReq.Title,
		OwnerId: fileReq.OwnerId,
	}
}

func (f *File) Prepare() {
	f.Title = html.EscapeString(strings.TrimSpace(f.Title))
	f.CreatedAt = time.Now()
}

func (f *FileRequest) Validate() string {
	if f.Title == "" {
		return constants.FileTitleCantBeEmpty
	}
	if f.OwnerId == 0 {
		return constants.UserIDInvalid
	}
	return ""
}

func (f *File) Validate() string {
	if f.Title == "" {
		return constants.FileTitleCantBeEmpty
	}
	if f.OwnerId == 0 {
		return constants.UserIDInvalid
	}
	return ""
}
