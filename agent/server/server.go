package server

import "lemna/agent"

type Server struct {
	target agent.Target
	info   *Info
}
