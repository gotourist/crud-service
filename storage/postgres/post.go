package postgres

import (
	"github.com/iman_task/crud-service/domain/entities"
	"github.com/iman_task/crud-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) repo.PostStorage {
	return &postRepo{
		db: db,
	}
}

func (p *postRepo) CreatePost(post *entities.Post) error {

	query := `
		INSERT INTO 
		    post(
		         "id", 
		         "title", 
		         "body",
		         "crated_at"
		         )
		values ($1, $2, $3, $4)`

	_, err := p.db.Exec(
		query,
		post.Id,
		post.Title,
		post.Body,
		post.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (p *postRepo) DetailPost(id int64) (*entities.Post, error) {
	var post entities.Post

	query := `
		SELECT 
		       "id", 
		       "title", 
		       "body"
		FROM post 
		WHERE "id"=$1 
		  AND "is_deleted"=false`

	err := p.db.Get(&post, query, id)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *postRepo) ListPosts(request *entities.ListPostsRequest) ([]*entities.Post, int64, error) {
	var (
		posts []*entities.Post
		count int64
	)

	query := `
		SELECT 
		       "id", 
		       "title", 
		       "body"
	FROM post
	WHERE "is_deleted"=false 
	ORDER BY "created_at" DESC 
	LIMIT $1 OFFSET $2`

	err := p.db.Select(&posts, query, request.Limit, request.Offset)
	if err != nil {
		return nil, 0, err
	}

	err = p.db.Get(&count, `SELECT count(1) FROM post WHERE "is_deleted=false"`)

	return posts, count, nil
}

func (p *postRepo) UpdatePost(request *entities.UpdatePostRequest) (tx *sqlx.Tx, err error) {
	tx, err = p.db.Beginx()
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE post 
		SET "title" = $1, 
		    "body" = $2 
		WHERE "id" = $3`

	_, err = tx.Exec(
		query,
		request.Title,
		request.Body,
		request.Id,
	)

	return tx, err
}

func (p *postRepo) DeletePost(id int64) (tx *sqlx.Tx, err error) {
	tx, err = p.db.Beginx()
	if err != nil {
		return nil, err
	}

	query := `
		UPDATE post 
		SET "is_deleted" = true 
		WHERE "id" = $1`

	_, err = tx.Exec(
		query,
		id,
	)

	return tx, err
}
