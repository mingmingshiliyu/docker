#!/usr/bin/env sh
# veth pair create
#创建两个网络ns
ip netns add ns1
ip netns add ns2
#创建一对veth
ip link add veth0 type veth peer name veth1
#连接
ip link set veth0 netns ns1
ip link set veth1 netns ns2
#验证,查看ns内veth
ip netns exec ns1 ip link

#配置每个veth网络地址和namespace的路由
#ns中配置ip
ip netns exec ns1 ifconfig veth0 172.18.0.2/24 up
ip netns exec ns2 ifconfig veth1 172.18.0.3/24 up
#ns中配置路由
ip netns exec ns1 route add default dev veth0
ip netns exec ns2 route add default dev veth1
#发包给veth
ip netns exec ns1 ping -c 172.18.0.3
ip netns exec ns2 ping -c 172.18.0.2