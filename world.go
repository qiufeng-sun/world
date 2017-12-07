package main

import (
	"util/logs"

	"core/server"

	"share/pipe"
	"share/rpc"
)

var _ = logs.Debug

//
type World struct {
	server.Server
}

//
func NewWorld() *World {
	return &World{}
}

//
func (this *World) Init() bool {
	// config
	if !LoadConfig("conf/") {
		return false
	}

	// recv/send msg among servers
	pipe.Init(Cfg.LanCfg, Cfg.EtcdCfg, handleMsgs)

	// init rpc client
	rpc.InitClient(Cfg.UserPoolCfg, Cfg.EtcdCfg.EtcdAddrs, Cfg.UserPath)

	return true
}

//
func (this World) String() string {
	return "World"
}
