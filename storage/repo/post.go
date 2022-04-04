package repo

import (
	"github.com/iman_task/crud-service/domain/entities"
	"github.com/jmoiron/sqlx"
)

type PostStorage interface {
	CreatePost(post *entities.Post) error
	DetailPost(id int64) (*entities.Post, error)
	ListPosts(request *entities.ListPostsRequest) ([]*entities.Post, int64, error)
	UpdatePost(request *entities.UpdatePostRequest) (*sqlx.Tx, error)
	DeletePost(id int64) (*sqlx.Tx, error)
}
