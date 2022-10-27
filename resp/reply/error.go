package reply

type UnkownErrReply struct {
}

var unkownErrReply = []byte("-Err unknown\r\n")

func (r *UnkownErrReply) Error() string {
	return "unkown err"
}

func (r *UnkownErrReply) ToBytes() []byte {
	return unkownErrReply
}

type ArgNumErrReply struct {
	Cmd string
}

func (r *ArgNumErrReply) Error() string {
	return "-Err wrong number of arguments for '" + r.Cmd + "'command"
}

func (r *ArgNumErrReply) ToBytes() []byte {
	return []byte("-Err wrong number of arguments for '" + r.Cmd + "'command\r\n")
}
func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{
		Cmd: cmd,
	}
}

type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")
var theSyntaxErrReply = &SyntaxErrReply{}

// MakeSyntaxErrReply creates syntax error
func MakeSyntaxErrReply() *SyntaxErrReply {
	return theSyntaxErrReply
}
func (r *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}
func (r *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE 0peration against a key holding the wrong kind of values\r\n")

// ToBytes marshals redis.Reply
func (r *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

func (r *WrongTypeErrReply) Error() string {
	return "WRONGTYPEOperation against a key holding the wrong kind of value"
}

type ProtocolErrReply struct {
	Msg string
}

func (r *ProtocolErrReply) ToBytes() []byte {
	return []byte("-ERR Protocol error: '" + r.Msg + "' \r\n")
}

func (r *ProtocolErrReply) Error() string {
	return "ERR Protocol error : '" + r.Msg
}
