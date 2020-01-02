// Code generated by go proto generate tool. DO NOT EDIT.
package protocol

import "gogame/protocol/pb"

const (
	MsgID_MessageSyncServer   = 4613 //服务器间同步
	MsgID_MessageGroupMsgNoti = 4614 //ss之间同步群组消息
	MsgID_MessagePing         = 4097 //ping消息
	MsgID_MessagePong         = 4098 //ping消息回
	MsgID_MessageCharLoginReq = 4099 //登录
	MsgID_MessageCharLoginRes = 4100 //登录结果
	MsgID_MessageErrorNoti    = 4101 //错误消息
)

func init() {
	processor = NewProcessor()
	processor.Register(MsgID_MessageSyncServer, (*pb.MessageSyncServer)(nil))
	processor.Register(MsgID_MessageGroupMsgNoti, (*pb.MessageGroupMsgNoti)(nil))
	processor.Register(MsgID_MessagePing, (*pb.MessagePing)(nil))
	processor.Register(MsgID_MessagePong, (*pb.MessagePong)(nil))
	processor.Register(MsgID_MessageCharLoginReq, (*pb.MessageCharLoginReq)(nil))
	processor.Register(MsgID_MessageCharLoginRes, (*pb.MessageCharLoginRes)(nil))
	processor.Register(MsgID_MessageErrorNoti, (*pb.MessageErrorNoti)(nil))
}
