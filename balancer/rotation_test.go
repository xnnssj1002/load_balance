package balancer

import (
	"fmt"
	"testing"
)

func TestLoadBalancer_RoundRotation(t *testing.T) {
	for i := 0; i < 2*len(LB.servers); i++ {
		serverInstance := LB.RoundRotation()
		fmt.Println(serverInstance)
	}
}
func TestLoadBalancer_RoundRotationOptimized(t *testing.T) {
	for i := 0; i < 2*len(LB.servers); i++ {
		serverInstance := LB.RoundRotationOptimized()
		fmt.Println(serverInstance)
	}
}

func TestLoadBalancer_WeightRotationByIndex(t *testing.T) {
	for i := 0; i < 10; i++ {
		serverInstance := LB.WeightRotationByIndex()
		fmt.Println(serverInstance)
	}
}

func TestLoadBalancer_WeightRotationByGap(t *testing.T) {
	for i := 0; i < 10; i++ {
		serverInstance := LB.WeightRotationByGap()
		fmt.Println(serverInstance)
	}
}

func TestLoadBalancer_WeightRotationByGapOptimized(t *testing.T) {
	for i := 0; i < 10; i++ {
		serverInstance := LB.WeightRotationByGapOptimized()
		fmt.Println(serverInstance)
	}
}

func TestLoadBalancer_WeightRotationSmooth(t *testing.T) {
	for i := 0; i < 10; i++ {
		serverInstance := LB.WeightRotationSmooth()
		fmt.Println(serverInstance)
	}
}
