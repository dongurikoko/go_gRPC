package handler

import (
	"context"
	"database/sql"
	"errors"
	"mygrpc/internal/usecase"
	todopb "mygrpc/pkg/grpc"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type todoServer struct {
	todoUseCase usecase.TodoUseCase
	todopb.UnimplementedTodoServiceServer
}

func NewTodoServer(tu usecase.TodoUseCase) *todoServer {
	return &todoServer{
		todoUseCase: tu,
	}
}

// TODOの新規作成
func (s *todoServer) Create(ctx context.Context, req *todopb.CreateTodoRequest) (*todopb.CreateTodoResponse, error) {
	// gRPCリクエストからtitleを取得
	title := req.GetTitle()

	// titleを元にTodoを作成
	id, err := s.todoUseCase.Insert(title)
	if err != nil {
		// エラーが発生した場合、gRPCステータスとメッセージを設定
		return nil, status.Errorf(codes.Internal, "Service Unavailable: %v", err)
	}

	// gRPCレスポンスを返却
	return &todopb.CreateTodoResponse{
		Id: int64(id),
	}, nil

}

// TODOの取得
func (s *todoServer) Get(ctx context.Context, req *todopb.GetTodosRequest) (*todopb.GetTodosResponse, error) {
	// gRPCリクエストから"query" を取得
	query := req.GetQuery()

	// クエリを元にTodoを取得
	todos, err := s.todoUseCase.GetAllByQuery(query)
	if err != nil {
		// エラーが発生した場合、gRPCステータスとメッセージを設定
		return nil, status.Errorf(codes.Internal, "Service Unavailable: %v", err)
	}

	// 取得したドメインモデルをレスポンスモデルに変換
	todosResponse := make([]*todopb.Todo, 0, len(todos))
	for _, todo := range todos {
		todosResponse = append(todosResponse, &todopb.Todo{
			Id:        int64(todo.ID),
			Title:     todo.Title,
			CreatedAt: timestamppb.New(todo.Created_at),
			UpdatedAt: timestamppb.New(todo.Updated_at),
		})
	}

	// gRPCレスポンスを返却
	return &todopb.GetTodosResponse{Todos: todosResponse}, nil
}

// TODOの更新
func (s *todoServer) Update(ctx context.Context, req *todopb.UpdateTodoRequest) (*todopb.UpdateTodoResponse, error) {
	//  gRPCリクエストからid,titleを取得
	id := req.GetId()
	title := req.GetTitle()

	// idを元にTodoを更新
	if err := s.todoUseCase.Update(int(id), title); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 指定されたTODOが存在しなかった場合
			return nil, status.Errorf(codes.NotFound, "Todo not found: %v", err)
		} else {
			// それ以外のエラーの場合
			return nil, status.Errorf(codes.Internal, "Failed to update Todo: %v", err)
		}
	}

	return &todopb.UpdateTodoResponse{Id: id}, nil
}

// TODOの削除
func (s *todoServer) Delete(ctx context.Context, req *todopb.DeleteTodoRequest) (*todopb.DeleteTodoResponse, error) {
	//  gRPCリクエストからidを取得
	id := req.GetId()

	// idを元にTodoを削除
	if err := s.todoUseCase.Delete(int(id)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// 指定されたTODOが存在しなかった場合
			return nil, status.Errorf(codes.NotFound, "Todo not found: %v", err)
		} else {
			// それ以外のエラーの場合
			return nil, status.Errorf(codes.Internal, "Failed to delete Todo: %v", err)
		}
	}

	return &todopb.DeleteTodoResponse{}, nil
}
