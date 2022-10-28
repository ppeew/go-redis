package handler

import (
	"context"
	"go-redis/database"
	databaseface "go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/resp/connection"
	"go-redis/resp/parser"
	"go-redis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map
	db         databaseface.Database
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	var db databaseface.Database
	//TODO 实现Database
	db = database.NewEchoDatabase()
	return &RespHandler{
		db: db,
	}
}

func (handler *RespHandler) closeClient(client *connection.Connection) {
	client.Close()
	handler.db.AfterClientClose(client)
	handler.activeConn.Delete(client)
}

// 核心业务
func (handler *RespHandler) Handler(ctx context.Context, conn net.Conn) {
	if handler.closing.Get() {
		conn.Close()
	}
	client := connection.NewConn(conn)
	handler.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		//错误
		if payload.Err != nil {
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF || strings.Contains(payload.Err.Error(), "use of closed network connection") {
				handler.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			//protocol error
			errReply := reply.MakeErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				handler.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		//exec
		if payload.Data == nil {
			continue
		}
		reply, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := handler.db.Exec(client, reply.Args)
		if result != nil {
			client.Write(result.ToBytes())
		} else {
			client.Write(unknownErrReplyBytes)
		}
	}
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
