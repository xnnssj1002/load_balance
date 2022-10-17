package balancer

import (
	"github.com/xnnssj1002/load_balance/server"
	"math/rand"
	"time"
)

// *********************** 随机算法 *********************** //

// Rand 所有服务器随机访问，实现简单，服务器的命中概率取决于随机算法，无法解决不同服务器之前性能差异问题
func (lb *LoadBalancer) Rand() *server.HttpServer {
	rand.Seed(time.Now().UnixNano()) // 添加随机种子
	index := rand.Intn(len(lb.servers))
	return lb.servers[index]
}

// *********************** 加权随机算法 *********************** //
// 实现上可以参照加权轮询，生成的随机数作为list列表的索引值，也可以降低服务器性能差异带来的问题

// WeightRandByIndex 加权随机算法 - 使用索引
func (lb *LoadBalancer) WeightRandByIndex() *server.HttpServer {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(ServerIndices))
	return lb.servers[ServerIndices[index]]
}

// WeightRandByGap 加权随机算法 - 使用区间
// 1、假设A:B:C = 5:2:1
// 2、遍历，获取总权重：5 + 2 + 1 = 8
// 3、得到三个区间，分别为[0, 5), [5, 7), [7, 8)。都是前闭后开区间
// 4、从区间[0, 8)随机获取一个数字，落到哪个区间，就返回哪个数字
func (lb *LoadBalancer) WeightRandByGap() *server.HttpServer {
	sumList := make([]int, len(lb.servers))
	sum := 0
	for i := 0; i < len(lb.servers); i++ {
		sum = sum + lb.servers[i].Weight
		sumList[i] = sum
	}
	// 获取[0, sum)中的一个随机数字
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(sum)
	for index, value := range sumList {
		if randNum < value {
			return lb.servers[index]
		}
	}
	return lb.servers[0]
}
