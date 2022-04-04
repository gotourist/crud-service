package service

import (
	configPkg "github.com/iman_task/crud-service/config"
	"github.com/iman_task/crud-service/domain/entities"
	"github.com/iman_task/crud-service/events/handlers"
	brokerpb "github.com/iman_task/crud-service/genproto/broker/post"
	loggerPkg "github.com/iman_task/crud-service/pkg/logger"
	broker "github.com/iman_task/crud-service/pkg/messagebroker"
	"github.com/iman_task/crud-service/storage"
	"github.com/jmoiron/sqlx"
)

type PostService struct {
	storage   storage.Storage
	logger    loggerPkg.Logger
	config    configPkg.Config
	publisher map[string]broker.Producer
}

func NewPostService(db *sqlx.DB, logger loggerPkg.Logger, config configPkg.Config, publisher map[string]broker.Producer) *PostService {
	return &PostService{
		storage:   storage.NewStoragePg(db),
		logger:    logger,
		config:    config,
		publisher: publisher,
	}
}

func (s *PostService) publishUpdatePostMessage(post *entities.Post) error {

	var postPb brokerpb.Post

	postPb.Id = post.Id
	postPb.Title = post.Title
	postPb.Body = post.Body
	postPb.IsDeleted = post.IsDeleted

	data, err := postPb.Marshal()
	if err != nil {
		return err
	}

	logBody := postPb.String()

	err = s.publisher[handlers.PostChangeTopic].Publish([]byte("post_change"), data, logBody)
	if err != nil {
		return err
	}

	return nil
}
