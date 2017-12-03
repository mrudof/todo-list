package main

import (
	"log"
	"net"

	// "golang.org/x/net/context"
	"github.com/mrudof/todo-list/backend/todolist"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) ListTodos(todo *todolist.Todo, stream todolist.TodoList_ListTodosServer) error {
	stream.Send(&todolist.Todo{Id: 1, Title: "First", DueDate: "asdf", Owner: "mrudof"})
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	todolist.RegisterTodoListServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	s.Serve(lis)
}
