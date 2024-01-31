package usecase

import (
	"errors"
	"fmt"
	"mygrpc/internal/domain/model"
	"mygrpc/internal/domain/repository"
)

type TodoUseCase interface {
	Insert(title string) (int, error)
	GetAllByQuery(query string) ([]*model.Todo, error)
	Update(id int, title string) error
	Delete(id int) error
}

type todoUseCase struct {
	todoRepository repository.TodoRepository
}

func NewTodoUseCase(tr repository.TodoRepository) TodoUseCase {
	return &todoUseCase{
		todoRepository: tr,
	}
}

// titleを指定してデータベースにTodoを追加し、追加したTodoのIDを返す
func (tu todoUseCase) Insert(title string) (int, error) {
	if title == "" {
		return 0, errors.New("title is empty")
	}
	// SQL文を実行する
	id, err := tu.todoRepository.Insert(title)
	if err != nil {
		return 0, fmt.Errorf("failed to insert todo in Insert: %w", err)
	}
	return id, nil
}

// クエリを元にTODOを取得する
func (tu todoUseCase) GetAllByQuery(query string) ([]*model.Todo, error) {
	// クエリが空の場合は全件取得
	if query == "" {
		return tu.todoRepository.GetAll()
	} else {
		// クエリが空でない場合はクエリに部分一致するTodoを取得
		return tu.todoRepository.GetAllByTitle(query)
	}
}

// idを指定してデータベースのTodoを更新する
func (tu todoUseCase) Update(id int, title string) error {
	if err := tu.todoRepository.Update(id, title); err != nil {
		return fmt.Errorf("failed to update todo in Update: %w", err)
	}
	return nil
}

// idを指定してデータベースのTodoを削除する
func (tu todoUseCase) Delete(id int) error {
	if err := tu.todoRepository.Delete(id); err != nil {
		return fmt.Errorf("failed to delete todo in Delete: %w", err)
	}
	return nil
}
