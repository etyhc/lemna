package agent

import (
	fmt "fmt"
	"io"
	"lemna/agent/rpc"
)

type client struct {
	id     int32
	stream rpc.Client_ForwardServer
	login  bool
}

func (c *client) session() string {
	return fmt.Sprintf("%d", c.id)
}

func (c *client) run() error {
	for {
		cfmsg, err := c.stream.Recv()
		if err != nil && err != io.EOF {
			return err
		}

		//转发指令
		if server, ok := serverMap[cfmsg.Target]; ok {
			//log.Println("c to s ", cfmsg.ServerType)
			cfmsg.Target = c.id
			server.stream.Send(cfmsg)
		} else {
			return fmt.Errorf("invaild servertype %d", cfmsg.Target)
		}
	}
}
