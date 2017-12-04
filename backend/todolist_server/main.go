package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mrudof/todo-list/backend/todolist"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) ListTodos(todo *todolist.Todo, stream todolist.TodoList_ListTodosServer) error {
	if err := stream.Send(&todolist.Todo{Id: 1, Title: "First", DueDate: "asdf", Owner: "mrudof", State: 0}); err != nil {
		panic(err)
	}
	return nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	logger := log.NewLogfmtLogger(os.Stdout)
	logger = level.NewFilter(logger, level.AllowInfo())
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := todolist.RegisterTodoListHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		level.Error(logger).Log("event", "server failed to start", err)
	}

	level.Info(logger).Log("event", "server starting")
	http.ListenAndServe(":8080", mux)
}
