package grpc_server

import (
	"context"
	"net"
	"os"
	"sync"

	"github.com/gookit/slog"
	auto_generate "github.com/rlapenok/test/internal/library/auto_generate/tg/proto"
	"google.golang.org/grpc"
)

type TgServerInterface interface {
	ChechChatID(int64) bool
	AddChatId(int64, chan<- *auto_generate.PrintLogRequest)
	DeleteChatId(int64)
	Run()
}

type GrpcServer struct {
	port             string
	rw               *sync.RWMutex
	storage_channels map[int64]chan<- *auto_generate.PrintLogRequest
}

func NewServer(port string) *GrpcServer {
	//Создание RWMutex для доступа к storage_channels
	rw := sync.RWMutex{}
	//Создание map для хранения sender для канала, по его chat_id
	storage_channels := make(map[int64]chan<- *auto_generate.PrintLogRequest)
	my_server := GrpcServer{rw: &rw, storage_channels: storage_channels, port: port}
	return &my_server
}

// Проверяет есть ли канал для связи с chat_id
func (srv *GrpcServer) ChechChatId(chat_id int64) bool {
	srv.rw.RLock()
	defer srv.rw.RUnlock()
	_, exist := srv.storage_channels[chat_id]
	if exist {
		return true
	} else {
		return false
	}
}

func (srv *GrpcServer) AddChatId(chat_id int64, sender chan<- *auto_generate.PrintLogRequest) {
	srv.rw.Lock()
	defer srv.rw.Unlock()
	srv.storage_channels[chat_id] = sender

}
func (srv *GrpcServer) DeleteChatId(chat_id int64) {
	srv.rw.Lock()
	defer srv.rw.Unlock()
	sender, exist := srv.storage_channels[chat_id]
	if exist {
		close(sender)
		delete(srv.storage_channels, chat_id)
	}
}

// Импелментация для создания gRPC сервера из авто сгенереного кода
func (srv *GrpcServer) PrintLog(ctx context.Context, incoming *auto_generate.PrintLogRequest) (*auto_generate.Null, error) {
	slog.Info("gRPC Server:New PrintLogRequest")
	srv.rw.RLock()
	defer srv.rw.RUnlock()
	wg := sync.WaitGroup{}
	for _, sender := range srv.storage_channels {
		wg.Add(1)
		go func(sender chan<- *auto_generate.PrintLogRequest) {
			sender <- incoming
			wg.Done()
		}(sender)
	}
	wg.Wait()
	return &auto_generate.Null{}, nil
}

// Запуск grpcServer
func (srv *GrpcServer) Run() {
	lis, err := net.Listen("tcp", ":"+srv.port)
	if err != nil {
		slog.Fatal(err.Error())
		os.Exit(1)
	}
	server := grpc.NewServer()
	auto_generate.RegisterTgBotServiceServer(server, srv)
	slog.Info("Starting gRPC server...")
	server.Serve(lis)
}
