package cache

import (
	"Sber/app/internal/model"
)

type Cache struct {
	Task map[int64]*model.Task
}

func NewCache() *Cache {
	return &Cache{
		Task: make(map[int64]*model.Task),
	}
}
