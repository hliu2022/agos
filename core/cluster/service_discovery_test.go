package cluster

import (
	"log"
	"testing"
	"time"
)

func TestServiceDiscovery(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser, err := NewServiceRegister(endpoints, "/web/node1", "abc", 5)
	if err != nil {
		log.Fatalln(err)
	}
	//监听续租相应chan
	go ser.ListenLeaseRespChan()
	//服务发现
	go serviceDiscovery()
	for {
		select {
		case <-time.After(10 * time.Second):
			ser.Close()
		}
	}
}

func serviceDiscovery() {
	var endpoints = []string{"localhost:2379"}
	ser := NewServiceDiscovery(endpoints)
	defer ser.Close()
	ser.WatchService("/web/")
	ser.WatchService("/gRPC/")
	for {
		select {
		case <-time.Tick(time.Second):
			log.Println(ser.GetServices())
		}
	}
}
