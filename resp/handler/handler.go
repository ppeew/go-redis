package handler

import (
	"bufio"
	"context"
	databaseface "go-redis/interface/database"
	"go-redis/lib/logger"
	"go-redis/lib/sync/atomic"
	"go-redis/lib/sync/wait"
	"io"
	"net"
	"sync"
	"time"
)

type RespHandler struct {
	activeConn sync.Map
	db         databaseface.Database
	closing    atomic.Boolean
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (hander *EchoHandler) Handler(ctx context.Context, conn net.Conn) {
	if hander.closing.Get() {
		conn.Close()
	}
	client := &EchoClient{
		Conn: conn,
	}
	hander.activeConn.Store(client, struct{}{})
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				hander.activeConn.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Waiting.Add(1)
		b := []byte(msg)
		conn.Write(b)
		client.Waiting.Done()
	}
}

func (hander *EchoHandler) Close() error {
	logger.Info("handler shutting down")
	hander.closing.Set(true)
	hander.activeConn.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		client.Close()
		return true
	})

	return nil
}

type EchoClient struct {
	Conn    net.Conn
	Waiting wait.Wait
}

func (e *EchoClient) Close() error {
	e.Waiting.WaitWithTimeout(10 * time.Second)
	e.Conn.Close()
	return nil
}
