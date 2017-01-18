package plugin

type Process interface {
	Kill() (err error)
}

type Plugin interface {
	GetDaemons() func() ([]Process, error)
	GetServices() func() ([]Process, error)
	GetDocker() func() ([]Process, error)
	FindAndKill() func([]Process, string, string) ([]Process, []error)
}
