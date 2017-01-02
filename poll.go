package gadu

import (
	"log"
	"syscall"
)

func (session GGSession) poller() {
	rd := new(syscall.FdSet)
	wr := new(syscall.FdSet)

	for {
		if session.session == nil {
			break
		}
		fdZero(rd)
		fdZero(wr)

		fd := (int)(session.session.fd)

		if session.session.check&ggCheckRead != 0 {
			fdSet(rd, fd)
		}
		if session.session.check&ggCheckWrite != 0 {
			fdSet(wr, fd)
		}
		n, err := syscallSelect(fd+1, rd, wr, nil, &syscall.Timeval{Sec: 1, Usec: 0})
		if err != nil {
			log.Fatalf("Unable to select(): %s", err)
			break
		}

		if n > 0 && (fdIsset(rd, fd) || fdIsset(wr, fd) || (session.session.timeout == 0 && session.session.soft_timeout > 0)) {
			e := session.watchFd()
			if e == nil {
				session.Close()
				break
			}
			session.events <- e
		}
	}
}

func fdSet(p *syscall.FdSet, i int) {
	p.Bits[i/64] |= 1 << uint(i) % 64
}

func fdIsset(p *syscall.FdSet, i int) bool {
	return (p.Bits[i/64] & (1 << uint(i) % 64)) != 0
}

func fdZero(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
}
