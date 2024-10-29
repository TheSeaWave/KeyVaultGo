package storage

import (
	"strconv"

	"go.uber.org/zap"
)

type Kind string

// Добавьте интерфейс StorageInterface в storage.go
type StorageInterface interface {
	Set(key, value string)
	Get(key string) *string
}

// Теперь измените `Storage`, чтобы он реализовал интерфейс
var _ StorageInterface = (*Storage)(nil) // Проверка на реализацию интерфейса

const (
	KindInt       Kind = "D"
	KindString    Kind = "S"
	KindUndefined Kind = "U"
)

type Storage struct {
	innerString map[string]string
	logger      *zap.Logger
}

func NewStorage() (Storage, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return Storage{}, err
	}

	defer logger.Sync()
	logger.Info("created new storage")

	return Storage{
		innerString: make(map[string]string),
		logger:      logger,
	}, nil
}

func (r Storage) Set(key, value string) {
	r.innerString[key] = value
	r.logger.Info("key set", zap.String("key", key), zap.String("value", value))
	r.logger.Sync()
}

func (r Storage) Get(key string) *string {
	res, ok := r.innerString[key]
	if !ok {
		return nil
	}
	return &res
}

func getType(value string) Kind {
	var val any

	val, err := strconv.Atoi(value)
	if err != nil {
		val = value
	}
	if val == "" {
		return KindUndefined
	}
	if val == "undefined" {
		return KindUndefined
	}
	switch val.(type) {
	case int:
		return KindInt
	case string:
		return KindString
	default:
		return KindUndefined
	}
}
