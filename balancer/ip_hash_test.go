package balancer

import (
	"fmt"
	"testing"
)

func TestLoadBalancer_IpHash(t *testing.T) {
	server := LB.IpHash("127.0.0.1")
	fmt.Println(server)
	for i := 0; i < len(LB.servers); i++ {
		instServer := LB.IpHash("127.0.0.1")
		fmt.Println(instServer)
		if instServer.Host != server.Host {
			t.Errorf("hope server host is %s, got host is %s", server.Host, instServer.Host)
		}
	}
}
