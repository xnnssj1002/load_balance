package balancer

import (
	"github.com/xnnssj1002/load_balance/server"
	"sort"
)

// *********************** 轮训算法 *********************** //
// 每次的请求到达时，对每个服务器都轮询访问，保证每个服务器命中概率相同，实现简单但无法解决不同服务器之间性能差异问题

// RoundRotation 轮训算法 - 冗余的写法
func (lb *LoadBalancer) RoundRotation() *server.HttpServer {
	serverInstance := lb.servers[lb.curIndex]
	lb.curIndex++
	if lb.curIndex >= len(lb.servers) {
		lb.curIndex = 0
	}
	return serverInstance
}

// RoundRotationOptimized 轮训算法 - 优化版本
func (lb *LoadBalancer) RoundRotationOptimized() *server.HttpServer {
	serverInstance := lb.servers[lb.curIndex]
	lb.curIndex = (lb.curIndex + 1) % len(lb.servers)
	return serverInstance
}

// *********************** 加权轮训算法 *********************** //
// 权重高的服务器请求命中的概率更高，根据不同服务器的性能调整权重比可以降低服务器性能差异带来的问题

// WeightRotationByIndex 加权轮训算法 - 使用索引
// 算法实现上可以将所有的服务器连接对象放到一个list中，按权重比例放不同数量的连接对象到list
// 比如有三台服务器权重比是，1:2:3，list中的连接对象数量数可以是1,2,3，轮询访问这个list即可
func (lb *LoadBalancer) WeightRotationByIndex() *server.HttpServer {
	serverInstance := lb.servers[ServerIndices[lb.curIndex]]
	lb.curIndex = (lb.curIndex + 1) % len(ServerIndices)
	return serverInstance
}

// WeightRotationByGap 加权轮训算法 - 使用区间 TODO 未完成
func (lb *LoadBalancer) WeightRotationByGap() *server.HttpServer {
	serverInstance := lb.servers[0]
	sum := 0
	sumList := make([]int, len(lb.servers)) // 1:2:3 [0, 1), [1, 3), [3, 6)
	for i := 0; i < len(lb.servers); i++ {
		sum += lb.servers[i].Weight
		sumList[i] = sum
	}

	// lb.curIndex 表示已经轮训到的小于sum的值
	for index, value := range sumList {
		stop := false
		for i := lb.curIndex; i < sum; i++ {
			if lb.curIndex < value {
				stop = true
				serverInstance = lb.servers[index]
			}
			lb.curIndex++
			if lb.curIndex >= sum {
				lb.curIndex = 0
			}
			if stop {
				break
			}
		}
		if stop {
			break
		}

	}
	return serverInstance
}

// WeightRotationByGapOptimized 加权轮训算法 - 使用区间优化版本
func (lb *LoadBalancer) WeightRotationByGapOptimized() *server.HttpServer {
	serverInstance := lb.servers[0]
	sum := 0
	// 1:2:3 [0, 1), [1, 3), [3, 6)
	for i := 0; i < len(lb.servers); i++ {
		sum += lb.servers[i].Weight // 第一次是1，[0, 1), [1, 3), [3, 6)
		if lb.curIndex < sum {
			serverInstance = lb.servers[i]
			if lb.curIndex == sum-1 && i != len(lb.servers)-1 {
				lb.curIndex++
			} else {
				lb.curIndex = (lb.curIndex + 1) % sum
			}
			break
		}
	}
	return serverInstance
}

// *********************** 平滑加权轮训算法 *********************** //
// 思路：
// 1、初始化权重{s1:3, s2:1, s3:1}，总权重为5，当前权重{s1:0, s2:0, s3:0}
// 2、每次循环前将 原始权重 加到 当前权重上
// 3、每次命中权重最大的，将其返回。然后把【命中节点】的当前权重数减去总权重数
// 以上做法的目的是：在五次循环后，能将各自节点的当前权重全部改为0

// 思路表格：初始化权重{s1:3, s2:1, s3:1}，总权重为5，当前权重{s1:0, s2:0, s3:0}
// 当前权重{s1:0, s2:0, s3:0}		命中		命中后的权重
// {s1:3, s2:1, s3:1}都加上原始权重	s1		{s1:-2, s2:1, s3:1}
// {s1:1, s2:2, s3:2}				s2		{s1:1, s2:-3, s3:2}
// {s1:4, s2:-2, s3:3)				s1		{s1:-1, s2:-2, s3:3}
// {s1:2, s2:-1, s3:4}				s3		{s1:2, s2:-1, s3:-1}
// {s1:5, s2:0, s3:0}				s1		{s1:0, s2:0, s3:0}

// WeightRotationSmooth 平滑加权轮训算法
func (lb *LoadBalancer) WeightRotationSmooth() *server.HttpServer {
	// 将当前权重，加上原始权重
	for _, s := range lb.servers {
		s.CWeight = s.Weight + s.CWeight
	}
	// 对 LoadBalance 里面的 servers 进行倒序排列
	sort.Sort(lb.servers)

	// 权重最大的服务，作为命中服务
	maxServer := lb.servers[0]

	// 命中服务的当前权重，需要减去原始总权重
	maxServer.CWeight = maxServer.CWeight - SumWeight

	//str := ""
	//for _, httpServer := range lb.servers {
	//	str = str + "," + strconv.Itoa(httpServer.CWeight)
	//}
	//fmt.Println(str)

	return maxServer
}
