package application

import (
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/domain/repository"
)

type commentApp struct {
	cm repository.CommentRepository
}

var _ CommentAppInterface = &commentApp{}

type CommentAppInterface interface {
	SaveComment(*entities.Comment) (*entities.Comment, map[string]string)
	GetComment(uint64) (*entities.Comment, error)
	GetCommentsByModel(uint64, uint64) ([]entities.Comment, error)
	UpdateComment(*entities.Comment) (*entities.Comment, map[string]string)
	DeleteComment(uint64) error
	DeleteCommentsByModel(uint64) error

	GetReplies(uint64, uint64) ([]entities.Comment, error)
}

func (c *commentApp) SaveComment(comment *entities.Comment) (*entities.Comment, map[string]string) {
	return c.cm.SaveComment(comment)
}

func (c *commentApp) GetComment(commentId uint64) (*entities.Comment, error) {
	return c.cm.GetComment(commentId)
}

func (c *commentApp) GetCommentsByModel(modelId, count uint64) ([]entities.Comment, error) {
	return c.cm.GetCommentsByModel(modelId, count)
}

func (c *commentApp) UpdateComment(comment *entities.Comment) (*entities.Comment, map[string]string) {
	return c.cm.UpdateComment(comment)
}

func (c *commentApp) DeleteComment(commentId uint64) error {
	return c.cm.DeleteComment(commentId)
}

func (c *commentApp) DeleteCommentsByModel(modelId uint64) error {
	return c.cm.DeleteCommentsByModel(modelId)
}

func (c *commentApp) GetReplies(parentId, count uint64) ([]entities.Comment, error) {
	return c.cm.GetReplies(parentId, count)
}
