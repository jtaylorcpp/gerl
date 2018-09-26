package gerl

type GenServer interface {
	Start()
	CallHandler()
	CastHandler()
	Terminate()
}

type GenServerClient interface {
	Init()
	Call()
	Cast()
	Terminate()
}

type GServer struct {
	Pid   ProcessID
	State interface{}
}

func NewGenServer(interface{}) Pid {
	pid := Pid{}
	return pid
}
