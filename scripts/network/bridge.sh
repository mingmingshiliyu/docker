#bridge虚拟设备是用于桥接的网络设备，相当于交换机，可以连接不同网络设备
#请求到达bridge，通过报文的mac地址进行广播和转发

#创建bridge设备，连接ns中网络设备和宿主机网络
  #创建veth并连入ns
ip netns add ns1
ip link add veth0 type veth peer name veth1
ip link set veth1 netns ns1
  #创建bridge并将veth连入
brctl addbr br0
brctl addif br0 veth0 #连入veth0
brctl addif br0 eth0 #连入网卡

#veth一端连入ns1，一端连入bridge，bridge连入eth0网卡

#linux路由表
#添加路由表
#启动veth0和br0虚拟设备
ip link set veth0 up
ip link set br0 up
#设置ns1路由
ip netns exec ns1 ifconfig veth1 172.18.0.2/24 up #veth0设置ip
ip netns exec ns1 route add default dev veth1 # default代表0.0.0.0/0，ns1所有流量都从veth1设备流出
#设置宿主机路由,宿主机将172.18.0.0/24网段的请求路由到br0
route add -net 172.18.0.0/24 dev br0
#验证
ifconfig eth0
  #从ns1访问宿主机
ip netns exec ns1 ping -c xxx
  #从宿主机访问ns1
ping -c 172.18.0.2


#iptables是对netfilter模块进行操作和展示的工具
#两种策略： MASQUERADE DNAT 用于容器和宿主机外部网络通信
#MASQUERADE:将请求包中源地址转换成一个网络设备地址,例如在宿主机内可以路由到某个ip，但到宿主机外部不知道如何路由到这个ip.如果需要请求外部地址，需要先通过MASQUERADE策略将ip转换成宿主机出口网卡的ip
  #打开ip转发
sysctl -w net.ipv4.conf.all.forwarding=1
  #对namespace中发出的包添加网络地址转换
iptables -t nat -A POSTROUTING -s 172.18.0.0/24 -o eth0 -j MASQUERADE
#在namespace中请求宿主机外部地址时,将namespace中源地址转换成宿主机地址作为源地址,就可以在namespace中访问宿主机外的网络

#DNAT:更换目标地址,用于将内部网络地址的端口映射到外部,namespace需要提供服务给宿主机之外的应用去请求,外部应用找不到172.18.0.2，需要dnat
iptables -t nat -A PREROUTING -p tcp -m tcp --dport 80 -j DNAT --to-destination 172.18.0.2:80
#这里将宿主机80端口的tcp请求转发到namespace中地址172.18.0.2:80，实现外部调用

