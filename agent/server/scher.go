package server

import (
	"container/ring"
)

type scher interface {
	sche() Server
	add(Server)
	del(Server)
}

type nilScher struct {
	ss map[Server]struct{}
}

func newNilScher() *nilScher {
	return &nilScher{ss: make(map[Server]struct{})}
}

func (ns *nilScher) sche() Server {
	for s := range ns.ss {
		return s
	}
	return nil
}

func (ns *nilScher) add(s Server) {
	ns.ss[s] = struct{}{}
}

func (ns *nilScher) del(s Server) {
	delete(ns.ss, s)
}

type roundScher struct {
	ring *ring.Ring
}

func newRoundScher() *roundScher {
	return &roundScher{}
}

func (rs *roundScher) sche() Server {
	if rs.ring != nil {
		rs.ring = rs.ring.Next()
		//logger.Debugf("sche server %v", rs.ring.Value.(Server))
		return rs.ring.Value.(Server)
	}
	return nil
}

func (rs *roundScher) add(s Server) {
	newring := ring.New(1)
	newring.Value = s
	rs.ring = newring.Link(rs.ring).Next()
	//logger.Debugf("add server %v", rs.ring.Value.(Server))
}

func (rs *roundScher) del(s Server) {
	if rs.ring == nil {
		return
	}
	n := rs.ring.Len()
	if n == 1 {
		if rs.ring.Value == s {
			rs.ring = nil
		}
		return
	}
	del := rs.ring
	for i := 0; i < n; i++ {
		if del.Value == s {
			if del == rs.ring {
				rs.ring = rs.ring.Next()
			}
			del.Prev().Unlink(1)
			return
		}
		del = del.Next()
	}
}

type loadScher struct {
	ss []Server
}

func newLoadScher() *loadScher {
	return &loadScher{ss: make([]Server, 0)}
}

func (ls *loadScher) sche() Server {
	n := len(ls.ss)
	if n == 0 {
		return nil
	}
	load, ret := ls.ss[0].info().Load, ls.ss[0]
	for i := 1; i < n; i++ {
		if load > ls.ss[i].info().Load {
			load, ret = ls.ss[i].info().Load, ls.ss[i]
		}
	}
	return ret
}

func (ls *loadScher) add(s Server) {
	ls.ss = append(ls.ss, s)
}

func (ls *loadScher) del(s Server) {
	for i := 0; i < len(ls.ss); i++ {
		if ls.ss[i] == s {
			ls.ss = append(ls.ss[:i], ls.ss[i+1:]...)
			return
		}
	}
}
