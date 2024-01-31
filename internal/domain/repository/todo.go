package repository

import (
	"mygrpc/internal/domain/model"
)

// infra層、usecase層がこのinterfaceに依存する
type TodoRepository interface {
	Insert(title string) (int, error)
	GetAll() ([]*model.Todo, error)
	GetAllByTitle(title string) ([]*model.Todo, error)
	Update(id int, title string) error
	Delete(id int) error
}
