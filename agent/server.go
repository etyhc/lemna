package agent

import (
	"lemna/agent/rpc"
)

type Server interface {
	Stream() rpc.Server_ForwardClient
}
