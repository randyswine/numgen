package control

type Cmd string

const (
	Start   Cmd = "start"
	Stop    Cmd = "stop"
	Destroy Cmd = "destroy"
)

type Signal string

const (
	Success Signal = "success"
)
