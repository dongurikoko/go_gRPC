package persistence

import (
	"database/sql"
	"fmt"
	"mygrpc/internal/domain/model"
	"mygrpc/internal/domain/repository"
)

type todoPersistence struct {
	Conn *sql.DB
}

func NewTodoPersistence(conn *sql.DB) repository.TodoRepository {
	return &todoPersistence{
		Conn: conn,
	}
}

// titleを指定してデータベースにTodoを追加し、追加したTodoのIDを返す
func (tp todoPersistence) Insert(title string) (int, error) {
	// SQL文を実行する
	result, err := tp.Conn.Exec("INSERT INTO todo (title) VALUES (?)", title)
	if err != nil {
		return 0, fmt.Errorf("failed to insert todo in Insert: %w", err)
	}

	// 追加したTodoのIDを取得
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert id in Insert: %w", err)
	}
	return int(lastInsertID), nil
}

// 作成したTodoの一覧を取得する
func (tp todoPersistence) GetAll() ([]*model.Todo, error) {
	rows, err := tp.Conn.Query("SELECT * FROM todo")
	if err != nil {
		return nil, fmt.Errorf("failed to select todo in GetAll: %w", err)
	}
	defer rows.Close()

	return convertToTodo(rows)
}

// 引数のstringに部分一致しているtitleを持つTodoの一覧を取得する
func (tp todoPersistence) GetAllByTitle(title string) ([]*model.Todo, error) {
	rows, err := tp.Conn.Query("SELECT * FROM todo WHERE title LIKE ?", "%"+title+"%")
	if err != nil {
		return nil, fmt.Errorf("failed to select todo in GetAllByTitle: %w", err)
	}
	defer rows.Close()

	return convertToTodo(rows)
}

// rows型をTodo型に変換する
func convertToTodo(rows *sql.Rows) ([]*model.Todo, error) {
	var todos []*model.Todo
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Created_at, &todo.Updated_at); err != nil {
			return nil, fmt.Errorf("failed to scan todo in convertToTodo: %w", err)
		}
		todos = append(todos, &todo)
	}
	return todos, nil
}

// idを指定してデータベースのTodoを更新する
func (tp todoPersistence) Update(id int, title string) error {
	result, err := tp.Conn.Exec("UPDATE todo SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return fmt.Errorf("failed to update todo in Update: %w", err)
	}

	// 更新した行数を取得
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected in Update: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// idを指定してデータベースのTodoを削除する
func (tp todoPersistence) Delete(id int) error {
	result, err := tp.Conn.Exec("DELETE FROM todo WHERE id = ?", id)
	if err != nil {
		return err
	}

	// 削除した行数を取得
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected in Delete: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
