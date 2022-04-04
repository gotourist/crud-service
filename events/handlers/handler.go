package handlers

import (
	configPkg "github.com/iman_task/crud-service/config"
	loggerPkg "github.com/iman_task/crud-service/pkg/logger"
	storage "github.com/iman_task/crud-service/storage"
)

const (
	PostAddTopic    = "post.add"
	PostChangeTopic = "post.change"
)

type EventHandler struct {
	conf    *configPkg.Config
	storage storage.Storage
	logger  loggerPkg.Logger
}

func NewEventHandler(storage storage.Storage, logger loggerPkg.Logger, conf configPkg.Config) *EventHandler {
	return &EventHandler{
		storage: storage,
		conf:    &conf,
		logger:  logger,
	}
}

func (e *EventHandler) Handle(topic string, value []byte) (msg string, err error) {
	switch topic {
	case PostAddTopic:
		msg, err = e.CreatePost(value)
		if err != nil {
			return "", err
		}
	}

	return msg, nil
}
