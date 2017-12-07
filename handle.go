package main

import (
	"util/logs"

	"core/net/dispatcher/pb"

	"share/handler"
	"share/msg"
	"share/pipe"
	"share/rpc"

	"github.com/golang/protobuf/proto"
)

var _ = logs.Debug

//
func handleMsgs(f *pb.PbFrame) {
	handler.HandleFrame(f)
}

//
func regFunc(msgId msg.EMsg, h func(f *pb.PbFrame)) {
	handler.RegFunc(int32(msgId), h)
}

//
func Call(msgId msg.EUserMsg, in, out proto.Message) error {
	return rpc.CallPbMsg(int32(msgId), in, out)
}

//
func init() {
	regFunc(msg.EMsg_ID_CSEnterWorld, handleEnterWorld)
}

//
func handleEnterWorld(f *pb.PbFrame) {
	//
	accId := f.GetAccId()

	// parse
	var m msg.CSEnterWorld
	e := handler.ParseMsgData(f.MsgRaw, &m)
	if e != nil {
		logs.Error("invalid msg! accId:%v, fromUrl:%v, error:%v", accId, f.GetSrcUrl(), e)
		return
	}
	logs.Info("user enter world: accId=%v, msg=%v", accId, m.String())

	// process
	resp, e := loadUser(accId)
	logs.Info("load user:%v, error:%v", resp.String(), e)

	// feedback// to do
	fb := &msg.SCEnterWorld{}

	pipe.SendMsg(accId, pipe.SrvUrl(), f.GetSrcUrl(), msg.EMsg_ID_SCEnterWorld, fb)
}

func loadUser(userId int64) (msg.LoadUserResp, error) {
	m := msg.LoadUserReq{UserId: proto.Int64(userId)}
	r := msg.LoadUserResp{}
	e := Call(msg.EUserMsg_ID_LoadUser, &m, &r)

	return r, e
}
