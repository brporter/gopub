package storage

import (
	"time"

	"github.com/brporter/gopub/models"
	"github.com/google/uuid"
)

type PostFactoryFunc func() (IPostRepo, error)

type IPostRepo interface {
	Open() error
	Close() error
	Save(p *models.Post) error
	FetchOne(id *uuid.UUID) (*models.Post, error)
	FetchMany(publishDate time.Time, pageSize int) ([]*models.Post, error)
	Remove(id *uuid.UUID) error
}
