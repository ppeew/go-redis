package database

import (
	"go-redis/interface/resp"
	"go-redis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}

func (d *EchoDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (d *EchoDatabase) Close() {
	panic("not implemented") // TODO: Implement
}

func (d *EchoDatabase) AfterClientClose(c resp.Connection) {
	panic("not implemented") // TODO: Implement
}
