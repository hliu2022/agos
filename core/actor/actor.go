package actor

type IActor interface {
	OnStart(server *Server) error
	OnStop(reason string) error
	OnCall(request interface{}) (interface{}, error)
	OnCast(request interface{})
}

type Server struct {
	Actor    IActor
	Id       string //ActorId
	ServerId string //服务Id
	ActiveAt int64  //开启时间
	//邮箱
	//定时器
}

func (ins *Server) Init(args []interface{}) (err error) {
	return nil
}
