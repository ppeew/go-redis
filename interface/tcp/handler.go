package tcp

import (
	"context"
	"net"
)

// 业务逻辑处理接口
type Handler interface {
	Handler(ctx context.Context, conn net.Conn)
	Close() error
}
