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
	Daemon      []Daemon
	Docker      []Docker
	Service     []Service
}
