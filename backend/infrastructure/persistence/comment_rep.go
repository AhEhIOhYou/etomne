package persistence

import (
	"errors"
	"github.com/AhEhIOhYou/etomne/backend/domain/entities"
	"github.com/AhEhIOhYou/etomne/backend/domain/repository"
	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

var _ repository.CommentRepository = &CommentRepo{}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

func (r *CommentRepo) SaveComment(comment *entities.Comment) (*entities.Comment, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&comment).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return comment, nil
}

func (r *CommentRepo) GetComment(id uint64) (*entities.Comment, error) {
	var comment entities.Comment
	err := r.db.Debug().Where("id = ?", id).Take(&comment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("comment not found")
	}
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	return &comment, nil
}

func (r *CommentRepo) GetCommentsByModel(id, count uint64) ([]entities.Comment, error) {
	var comments []entities.Comment
	err := r.db.Debug().Table("comments").Order("created_at asc").
		Where("model_id = ? AND reply_id is null", id).Limit(int(count)).Find(&comments).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("comment not found")
	}
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepo) UpdateComment(comment *entities.Comment) (*entities.Comment, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&comment).Error
	if err != nil {
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return comment, nil
}

func (r *CommentRepo) DeleteComment(id uint64) error {
	var comment entities.Comment
	err := r.db.Debug().Where("id = ?", id).Delete(&comment).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

func (r *CommentRepo) DeleteCommentsByModel(id uint64) error {
	var comment entities.Comment
	err := r.db.Debug().Table("comments").Where("model_id = ?", id).Delete(&comment).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}

func (r *CommentRepo) GetReplies(parentId, count uint64) ([]entities.Comment, error) {
	var comments []entities.Comment
	err := r.db.Debug().Table("comments").Order("created_at asc").
		Where("reply_id = ?", parentId).Limit(int(count)).Find(&comments).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("comment not found")
	}
	if err != nil {
		return nil, err
	}
	return comments, nil
}
