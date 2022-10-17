package balancer

import (
	"github.com/xnnssj1002/load_balance/server"
)

// LB 包全局变量
var LB *LoadBalancer

// ServerIndices 供各种加权算法使用，服务的权重有多大，在该切片中就存在多少个对应服务的索引值
// eg: server1 weight = 3; server2 weight = 2; server3 weight = 1
// server1在HttpServers索引为0，server2索引为1， server3索引为2
// ServerIndices 的元素为[0, 0, 0, 1, 1, 2]
var ServerIndices []int

// SumWeight 总权重数字
var SumWeight int

func init() {
	LB = NewLoadBalancer()
	LB.AddServer(server.NewHttpServer("http://127.0.0.1:90910", 1))
	LB.AddServer(server.NewHttpServer("http://127.0.0.1:90920", 2))
	LB.AddServer(server.NewHttpServer("http://127.0.0.1:90930", 3))

	for index, serverInstance := range LB.servers {
		for i := 0; i < serverInstance.Weight; i++ {
			ServerIndices = append(ServerIndices, index)
		}
		SumWeight = SumWeight + serverInstance.Weight
	}
}

type HttpServers []*server.HttpServer

func (x HttpServers) Len() int           { return len(x) }
func (x HttpServers) Less(i, j int) bool { return x[i].CWeight > x[j].CWeight } // 倒序，即从大到小
func (x HttpServers) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// LoadBalancer 负载均衡器
type LoadBalancer struct {
	servers HttpServers

	// curIndex 当使用轮训算法时，该字段表示已经轮训到的服务在 servers 中的索引值
	// 当使用加权轮训算法是，该字段表示已经轮训到的服务在 ServerIndices 中的索引值
	curIndex int
}

// NewLoadBalancer create a LoadBalancer instance
func NewLoadBalancer() *LoadBalancer {
	return &LoadBalancer{
		servers: make(HttpServers, 0),
	}
}

// AddServer Add  Server in LoadBalances
func (lb *LoadBalancer) AddServer(server *server.HttpServer) {
	lb.servers = append(lb.servers, server)
}
