package server

import "lemna/agent"

type Server interface {
	agent.Target
	info() *Info
}
