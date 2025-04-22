package server

import (
	"sync"
	"time"
)

type Address struct {
	serviceName string
	addr        string
	expireAt    int64
}

type nameStore struct {
	data         map[string]map[string]*Address
	dataLocker   sync.RWMutex
	expireAtData map[int64][]*Address
	expireLocker sync.RWMutex
}

var serviceNameData *nameStore
var expireTs = time.Second * 10

func init() {
	serviceNameData = &nameStore{
		data:         map[string]map[string]*Address{},
		expireAtData: map[int64][]*Address{},
	}

	// 异步删除过期的注册信息
	go serviceNameData.clearData()
}

func GetAllData() *nameStore {
	return serviceNameData
}

// 将服务名与地址注册到集合
func Register(serviceName, address string) {
	ns := serviceNameData
	addr := &Address{
		serviceName: serviceName,
		addr:        address,
		expireAt:    time.Now().Add(expireTs).Unix(),
	}
	ns.expireLocker.Lock()
	expireAtKey := addr.expireAt
	_, ok := ns.expireAtData[expireAtKey]
	if !ok {
		ns.expireAtData[expireAtKey] = make([]*Address, 0)
	}
	ns.expireAtData[expireAtKey] = append(ns.expireAtData[expireAtKey], addr)
	ns.expireLocker.Unlock()

	ns.dataLocker.Lock()
	_, ok = ns.data[serviceName]
	if !ok {
		ns.data[serviceName] = make(map[string]*Address, 0)
	}
	ns.data[serviceName][address] = addr
	ns.dataLocker.Unlock()
}

// 删除注册信息
func Delete(serviceName, address string) {
	ns := serviceNameData
	ns.dataLocker.Lock()
	ns.deleteNotLock(serviceName, address)
	ns.dataLocker.Unlock()
}

func (ns nameStore) deleteNotLock(serviceName, address string) {
	_, ok := ns.data[serviceName]
	if ok {
		delete(ns.data[serviceName], address)
	}
}

// 修改注册信息有效时间
func Keepalive(serviceName, address string) {
	ns := serviceNameData
	_, ok := ns.data[serviceName]
	if !ok {
		Register(serviceName, address)
		//return
	}
	_, ok = ns.data[serviceName][address]
	if !ok {
		Register(serviceName, address)
		//return
	}
	ns.dataLocker.Lock()
	addr := ns.data[serviceName][address]
	addr.expireAt = time.Now().Add(expireTs).Unix()
	ns.data[serviceName][address] = addr
	ns.dataLocker.Unlock()

	ns.expireLocker.Lock()
	expireAtKey := addr.expireAt
	_, ok = ns.expireAtData[expireAtKey]
	if !ok {
		ns.expireAtData[expireAtKey] = make([]*Address, 0)
	}
	ns.expireAtData[expireAtKey] = append(ns.expireAtData[expireAtKey], addr)
	ns.expireLocker.Unlock()
}

// 根据服务名称获取地址信息
func GetByServiceName(serviceName string) []string {
	ns := serviceNameData
	ns.dataLocker.RLock()
	defer ns.dataLocker.RUnlock()
	mp, ok := ns.data[serviceName]
	if !ok {
		return []string{}
	}
	if len(mp) == 0 {
		return []string{}
	}
	list := make([]string, 0)
	for _, addr := range mp {
		if addr.expireAt <= time.Now().Unix() {
			continue
		}
		list = append(list, addr.addr)
	}
	return list
}

func (ns nameStore) clearData() {
	timeTicker := time.NewTicker(3 * time.Second)
	defer timeTicker.Stop()

	for {
		select {
		case <-timeTicker.C:
			list := make([]*Address, 0)
			ns.expireLocker.Lock()
			for key, items := range ns.expireAtData {
				if key <= time.Now().Unix()-10 {
					list = append(list, items...)
					delete(ns.expireAtData, key)
				}
			}
			ns.expireLocker.Unlock()
			ns.dataLocker.Lock()
			for _, item := range list {
				if item.expireAt <= time.Now().Unix()-10 {
					ns.deleteNotLock(item.serviceName, item.addr)
				}
			}
			ns.dataLocker.Unlock()
		}
	}
}
