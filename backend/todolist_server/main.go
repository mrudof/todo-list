package main

import (
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/julienschmidt/httprouter"
	"github.com/mrudof/todo-list/backend/todolist"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
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

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("logger failed: %v", err)
	}

	l, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err)
	}

	mux := runtime.NewServeMux()
	s := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_zap.UnaryServerInterceptor(logger),
		)),
	)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(s)
	go s.Serve(l)
	todolist.RegisterTodoListServer(s, &server{})
	opts := []grpc.DialOption{grpc.WithInsecure()}
	if err := todolist.RegisterTodoListHandlerFromEndpoint(ctx, mux, "localhost:50051", opts); err != nil {
		logger.Error(err.Error())
	}

	r := httprouter.New()
	r.Handler("GET", "/metrics", promhttp.Handler())
	r.Handler("GET", "/api/todo/list", mux)
	logger.Info("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}
