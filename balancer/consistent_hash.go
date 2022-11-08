package balancer

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 一致性哈希

type UInt32Slice []uint32

func (s UInt32Slice) Len() int           { return len(s) }
func (s UInt32Slice) Less(i, j int) bool { return s[i] < s[j] }
func (s UInt32Slice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Hash 定义了函数类型 Hash
// 采取依赖注入的方式，允许用于替换成自定义的 Hash 函数，也方便测试时替换.
// 默认为 crc32.ChecksumIEEE 算法。
type Hash func(data []byte) uint32

// Map 是一致性哈希算法的主数据结构，包含4个成员变量：哈希函数；虚拟结点倍数；哈希环；虚拟节点与真实节点的映射表
type Map struct {
	hash     Hash              // hash函数
	replicas int               // 虚拟节点倍数 解决数据倾斜问题
	keys     UInt32Slice       // 哈希环。已排序的节点哈希切片
	hashMap  map[uint32]string // 虚拟节点与真实节点的映射表 hashMap，键是虚拟节点的哈希值，值是真实节点的名称。
}

// New 允许自定义虚拟节点倍数和 Hash 函数。
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[uint32]string),
	}

	// 默认使用CRC32算法
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// Add 方法用来添加缓存节点，参数为节点key，比如使用IP
// Add 函数允许传入 0 或 多个真实节点的名称。
// 对每一个真实节点 key，对应创建 m.replicas 个虚拟节点，虚拟节点的名称是：strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点。
// 使用 m.hash() 计算虚拟节点的哈希值，使用 append(m.keys, hash) 添加到环上。
// 在 hashMap 中增加虚拟节点和真实节点的映射关系。
// 最后一步，环上的哈希值排序。
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 结合复制因子计算所有虚拟节点的hash值，并存入m.keys中，同时在m.hashMap中保存哈希值和key的映射
		for i := 0; i < m.replicas; i++ {
			hash := m.hash([]byte(strconv.Itoa(i) + key))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	// 对所有虚拟节点的哈希值进行排序，方便之后进行二分查找
	sort.Sort(m.keys)
}

// Get 方法根据给定的对象获取最靠近它的那个节点key
// 第一步，计算 key 的哈希值。
// 第二步，顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值。如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
// 第三步，通过 hashMap 映射得到真实的节点。
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}

	hash := m.hash([]byte(key))

	// 通过二分查找获取最优节点，第一个节点hash值大于对象hash值的就是最优节点
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// 如果查找结果大于节点哈希数组的最大索引，表示此时该对象哈希值位于最后一个节点之后，那么放入第一个节点中
	if idx == len(m.keys) {
		idx = 0
	}

	return m.hashMap[m.keys[idx]]
}
