package agent

import (
	fmt "fmt"
	"io"
	"log"
)

type client struct {
	id     int32
	stream Client_ForwardServer
	login  bool
}

func (c *client) session() string {
	return fmt.Sprintf("%d", c.id)
}

func (c *client) run() error {
	for {
		msg, err := c.stream.Recv()
		if err != nil && err != io.EOF {
			return err
		}

		//转发指令
		if server, ok := serverMap[msg.ServerType]; ok {
			log.Println("c to s ", msg.ServerType)
			server.stream.Send(&ServerFwdMsg{ClientID: c.id, Msg: msg.Msg})
		} else {
			return fmt.Errorf("invaild servertype %d", msg.ServerType)
		}
	}
}
