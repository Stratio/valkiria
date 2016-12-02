package proc

var Sessions []Session

const (
	CHAOS = iota
	SHOOTER
)

type Session struct {
	Id          int
	Start       int
	Finish      int
	SessionType int
	Daemon      []daemon
	Docker      []docker
	Service     []service
}

func IsSessionLock() (res bool) {
	if len(Sessions) != 0 {
		if Sessions[len(Sessions)-1].Finish == 0 {
			res = true
		}
	}
	return
}
