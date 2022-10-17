package balancer

import (
	"fmt"
	"testing"
)

func TestLoadBalancer_Rand(t *testing.T) {
	for i := 0; i < len(LB.servers); i++ {
		server := LB.Rand()
		fmt.Println(server.Host)
	}
}

func TestLoadBalancer_WeightRandByIndex(t *testing.T) {
	for i := 0; i < len(ServerIndices); i++ {
		server := LB.WeightRandByIndex()
		fmt.Println(server)
	}
}

func TestLoadBalancer_WeightRandByGap(t *testing.T) {
	for i := 0; i < len(LB.servers); i++ {
		instServer := LB.WeightRandByGap()
		fmt.Println(instServer)
	}
}
