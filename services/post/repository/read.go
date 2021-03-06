package repository

import (
	"microServiceBoilerplate/services/post/models"

	"github.com/mreza0100/golog"
	"gorm.io/gorm"
)

type read struct {
	lgr *golog.Core
	db  *gorm.DB
}

func (r *read) GetPost(postIds []uint64) ([]*models.Post, error) {
	const query = `SELECT * FROM posts WHERE id IN(?)`
	params := []interface{}{postIds}

	tx := r.db.Raw(query, params...)

	if tx.Error != nil {
		return nil, tx.Error
	}

	result := make([]*models.Post, 0, len(postIds))
	tx.Scan(&result)

	return result, nil
}

func (r *read) IsExists(postIds []uint64) ([]uint64, error) {
	const query = `SELECT id FROM posts WHERE id IN (?)`
	params := []interface{}{postIds}

	tx := r.db.Raw(query, params...)
	if tx.Error != nil {
		return nil, tx.Error
	}

	result := make([]uint64, 0)
	tx.Scan(&result)

	return result, nil
}
