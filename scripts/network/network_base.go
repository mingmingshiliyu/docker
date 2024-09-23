package network

import (
	"net"

	"github.com/vishvananda/netlink"
)

/*
需要三个包：
	net:
		net.IP ip数据结构*(172.168.0.1)，ParseIP和String方法将字符串与ip转换
		net.IPNet 定义ip数据结构,通过parseCIDR和String方法转换网段(172.168.0.0/16)和字符串
	github.com/vishvananda/netlink:
		操作网络接口和路由表配置的库
		netlink.LinkAdd(xxxlink)
		netlink.LinkDel(xxxlink)
	github.com/vishvananda/netns
		在namespace中执行代码
		netns.Set(container_netns)

*/

// 网络模型
// 通过网络和网络中别的容器互连，包含网络的容器地址段,网络操作调用的网络驱动
type Network struct {
	Name    string
	IpRange *net.IPNet //网段
	Driver  string     //网络驱动名
}

//网络端点:连接容器和网络的，保证容器内部和网络的通信，包含连接到网络的一些信息,如地址,veth设备，端口映射,连接的容器和网络
//一端连接到容器内部，一端挂载到bridge。信息传输需要两个组件配合：网络驱动和IPAM

type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`
	IPAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	PortMapping []string         `json:"portmapping"`
	Network     *Network
}

// 网络驱动
type NetworkDriver interface {
	//驱动名
	Name() string
	//创建网络
	Create(subnet string, name string) (*Network, error)
	//删除网络
	Delete(network Network) error
	//连接容器网络端点到网络
	Connect(network *Network, endpoint *Endpoint) error
	//从网络删除容器网络端点
	Disconnect(network Network, endpoint *Endpoint) error
}

//IPAM
//用于网络ip地址分配和释放,包括容器ip地址和网络网关ip地址
//IPAM.Allocate(subnet *net.IPNet)从指定subnet网段中分配ip地址
//IPAM.Release(subnet net.IPNet,ipaddr net.IP) 从指定subnet网段中释放指定ip地址

//调用关系
//mydocker network create --subnet 192.168.0.0/24 --driver bridge testbridgenet
//通过bridge网络驱动创建一个网络，网段是192.168.0.0/24，网络驱动是Bridge
//流程:
//创建网络-》docker-》从ipam获取ip和gatewayIR-》ipam返回网段配置-》docker发送给network driver创建网络-》networkdriver配置网络设备-》返回配置好的网络信息
//IPAM:传入的ip网段分配一个可用ip地址给容器和网络的网关,比如网络网段192.168.0.0/16，通过ipam获取这个网段的容器地址分配给容器的连接端点，保证ip不会冲突
//Network Driver:网络管理,创建网络时用于网络初始化以及网络端点配置，比如bridge驱动的动作就是创建bridge和挂载veth设备。

// 创建网络
func CreateNetwork(driver, subnet, name string) error {
	//网段字符串转换成net.IPNet对象
	cidr, ipNet, err := net.ParseCIDR(subnet)
	//通过ipam分配网关ip,获取网段中第一个ip作为网关ip
	//调用指定网络驱动创建网络,driver字典是网络驱动的实例字典,

	//保存网络信息，网络信息存入文件系统,以便查询和在网络上连接网络端点

}

func (nw *Network) dump(dumpPath string) error {
	//检查目录是否存在，不存在创建

	//保存的文件名是网络名
	//打开文件，文件权限分别是清空内容，只写入，不存在则创建
	//序列化成字节并写入文件
}

func (nw *Network) load(dumpPath string) error {
	//打开配置文件
	//从配置文件读取网络配置的json
	//通过json字符串反序列化成网络结构
}
