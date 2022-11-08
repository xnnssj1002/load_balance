package balancer

import (
	"github.com/xnnssj1002/load_balance/server"
	"hash/crc32"
)

// IP Hash 算法 即hash取余
// 对请求的ip地址用hash算法映射到服务器上，保证一个客户端的所有请求都命中到一台服务器上。
// 适合服务端保存客户端的状态，开启session会话的情况。
// 但是不能跨服务器会话，如果服务器有新上线，下线，重启等导致服务器序号发生改变时会导致此种策略异常

// IpHash IP Hash 算法
func (lb *LoadBalancer) IpHash(ip string) *server.HttpServer {
	index := int(crc32.ChecksumIEEE([]byte(ip))) % len(lb.servers)
	return lb.servers[index]
}
