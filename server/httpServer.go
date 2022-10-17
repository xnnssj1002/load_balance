package server

// this file is http server

// Status 为 HttpServer 的状态值
type Status int

// server status StatusDown is closed. StatusUp is opened
const (
	StatusDown Status = iota
	StatusUp
)

// HttpServer this is a httpServer
type HttpServer struct {
	Host    string
	Weight  int    // 权重
	CWeight int    // 当前权重
	Status  Status // 当前服务的状态值
}

// NewHttpServer create a HttpServer instance
func NewHttpServer(host string, weight int) *HttpServer {
	return &HttpServer{
		Host:    host,
		Weight:  weight,
		CWeight: 0,
		Status:  StatusUp,
	}
}
