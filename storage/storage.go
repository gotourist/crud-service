package storage

import (
	"database/sql"

	"github.com/iman_task/crud-service/storage/postgres"
	"github.com/iman_task/crud-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoRows = sql.ErrNoRows
)

type Storage interface {
	Post() repo.PostStorage
}
type storagePg struct {
	db       *sqlx.DB
	postRepo repo.PostStorage
}

func NewStoragePg(db *sqlx.DB) Storage {
	return &storagePg{
		db:       db,
		postRepo: postgres.NewPostRepo(db),
	}
}

func (s storagePg) Post() repo.PostStorage {
	return s.postRepo
}
