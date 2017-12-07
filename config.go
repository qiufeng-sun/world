package main

import (
	"path/filepath"

	"github.com/astaxie/beego/config"

	"util/etcd"
	"util/logs"
	"util/logs/scribe"

	"core/net/lan"
	"core/net/lan/rpc"
)

var _ = logs.Debug

//
type Config struct {
	// scribe
	Scribed    bool
	ScribeAddr string

	// server
	LanCfg  *lan.LanCfg
	EtcdCfg *etcd.SrvCfg

	// rpc
	UserPoolCfg *rpc.PoolConfig
	UserPath    string
}

func (this *Config) init(fileName string) bool {
	confd, e := config.NewConfig("ini", fileName)
	if e != nil {
		logs.Panicln("load config file failed! file:", fileName, "error:", e)
	}

	//[scribe]
	scribe.Init("user", confd)

	//[server]
	srvName := confd.String("server::name")
	srvAddr := confd.String("server::addr")
	this.LanCfg = lan.NewLanCfg(srvName, srvAddr)

	//[etcd]
	this.EtcdCfg = &etcd.SrvCfg{}
	this.EtcdCfg.EtcdAddrs = confd.Strings("etcd::addrs")
	this.EtcdCfg.SrvAddr = srvAddr
	this.EtcdCfg.SrvRegPath = confd.String("etcd::reg_path")
	this.EtcdCfg.SrvRegUpTick = confd.DefaultInt64("etcd::reg_uptick", 2000)

	this.EtcdCfg.WatchPaths = confd.Strings("etcd::watch_path")

	//[rpc_user]
	this.UserPoolCfg = &rpc.PoolConfig{
		Name: confd.String("rpc_user::name"),
		InitNum: confd.DefaultInt("rpc_user::init_num",1),
		IdleNum: confd.DefaultInt("rpc_user::idle_num",1),
		MaxNum: confd.DefaultInt("rpc_user::max_num",1),
	}
	this.UserPath = confd.String("rpc_user::path")

	// echo
	logs.Info("user config:%+v", *this)

	return true
}

//
var Cfg = &Config{}

//
func LoadConfig(confPath string) bool {
	// config
	confFile := filepath.Clean(confPath + "/self.ini")

	return Cfg.init(confFile)
}

//
func SrvId() string {
	return Cfg.LanCfg.ServerId()
}

//
func SrvName() string {
	return Cfg.LanCfg.Name
}

// to do add check func
