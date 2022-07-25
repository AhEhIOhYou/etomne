package repository

import "github.com/AhEhIOhYou/etomne/backend/domain/entities"

type CommentRepository interface {
	SaveComment(*entities.Comment) (*entities.Comment, map[string]string)
	GetComment(uint64) (*entities.Comment, error)
	GetCommentsByModel(uint64, uint64) ([]entities.Comment, error)
	UpdateComment(*entities.Comment) (*entities.Comment, map[string]string)
	DeleteComment(uint64) error
	DeleteCommentsByModel(uint64) error

	GetReplies(uint64, uint64) ([]entities.Comment, error)
}
