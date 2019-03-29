package agent

type Balancer interface {
	GetServer(target int32, clientid int32) (Server, bool)
}
