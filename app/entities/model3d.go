package entities

type Model3d struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CreateDate  string `json:"create_date"`
	Description string `json:"description"`
	FileId      int64  `json:"fileId"`
	FilePath    string `json:"file_path"`
}
