package proc

import (
	"syscall"
	"time"
)

type Service struct {
	Pid            uint64
	Name           string
	KillName       string
	FrameWorkName  string
	Ppid           int64
	Executor       bool
	ChaosTimeStamp int64
}

func (s *Service) Kill() (err error) {
	err = syscall.Kill(int(s.Pid), 9)
	if err != nil {
		s.ChaosTimeStamp = time.Now().UTC().UnixNano()
	}
	return
}






