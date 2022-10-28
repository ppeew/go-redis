package handler

import (
	"context"
	databaseface "go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/resp/connection"
	"net"
	"sync"
)

type RespHandler struct {
	activeConn sync.Map
	db         databaseface.Database
	closing    atomic.Boolean
}

func (handler *RespHandler) Handler(ctx context.Context, conn net.Conn) {
	panic("not implemented") // TODO: Implement
}

func (handler *RespHandler) Close() error {
	logger.Info("handler shutting down")
	handler.closing.Set(true)
	handler.activeConn.Range(func(key, value any) bool {
		client := key.(*connection.Connection)
		client.Close()
		return true
	})
	handler.db.Close()
	return nil
}
