package ping

type service struct {
}

func (s service) Ping() *pong {
	return &pong{Message: "pong"}
}
