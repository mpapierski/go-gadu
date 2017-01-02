package gadu

import (
	"log"
	"syscall"
)

func (session GGSession) poller() {
	rd := new(syscall.FdSet)
	wr := new(syscall.FdSet)
	fd := (int)(session.session.fd)
	for {
		FD_ZERO(rd)
		FD_ZERO(wr)
		if session.session.check&GG_CHECK_READ != 0 {
			FD_SET(rd, fd)
		}
		if session.session.check&GG_CHECK_WRITE != 0 {
			FD_SET(wr, fd)
		}
		_, err := syscall.Select(fd+1, rd, wr, nil, &syscall.Timeval{Sec: 1, Usec: 0})
		if err != nil {
			log.Fatalf("Unable to select(): %s", err)
			break
		}

		if FD_ISSET(rd, fd) || FD_ISSET(wr, fd) {
			e := session.watchFd()
			if e == nil {
				break
			}
			session.events <- e
		}
	}
}

func FD_SET(p *syscall.FdSet, i int) {
	p.Bits[i/64] |= 1 << uint(i) % 64
}

func FD_ISSET(p *syscall.FdSet, i int) bool {
	return (p.Bits[i/64] & (1 << uint(i) % 64)) != 0
}

func FD_ZERO(p *syscall.FdSet) {
	for i := range p.Bits {
		p.Bits[i] = 0
	}
}
