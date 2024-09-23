package network

import (
	"encoding/json"
	"os"
	"path"
)

const (
	defaultIpamAllocatorPath = "/var/run/mydocker/network/ipam/subnet.json"
)

type IPAM struct {
	//分配文件存放位置
	SubnetAllocatorPath string
	//网段和位图算法map，key网段，value分配的位图数组
	//value表示网段中某些ip是否已分配
	Subnets *map[string][]uint64
}

var ipAllocator = &IPAM{
	SubnetAllocatorPath: defaultIpamAllocatorPath,
}

func (p *IPAM) load() error {
	_, err := os.Stat(p.SubnetAllocatorPath)
	if err != nil {
		return err
	}
	file, err := os.Open(p.SubnetAllocatorPath)
	if err != nil {
		return err
	}
	var content []byte
	n, err := file.Read(content)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content[:n], p.Subnets)
	if err != nil {
		return err
	}
	return nil
}

func (p *IPAM) dump() error {
	dir, _ := path.Split(p.SubnetAllocatorPath)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(dir, 0644)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	openFile, err := os.OpenFile(p.SubnetAllocatorPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	marshal, err := json.Marshal(p.Subnets)
	if err != nil {
		return err
	}
	_, err = openFile.Write(marshal)
	if err != nil {
		return err
	}
	return nil
}

type Bitmap struct {
	bits []uint64
}

type Bitmapper interface {
	Insert(num uint) //num位插入元素
	Delete(num uint)
	Check(num uint) (b bool)
	All() (nums []uint) //返回所有存储的元素下标
	Clear()
}

// https://zhuanlan.zhihu.com/p/423253110
func New() (bm *Bitmap) {
	return &Bitmap{bits: make([]uint64, 0, 0)}
}
