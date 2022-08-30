package example

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/anypb"
	pd "zrpc/rpc/proto"
)

type Test struct {
}

// QueryProto 用于测试用户查询的方法
func (t *Test) QueryProto(ctx context.Context, arg *anypb.Any) (reply *anypb.Any, err error) {
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", arg)
	unmarshal := &pd.Args{}
	ptypes.UnmarshalAny(arg, unmarshal)
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", unmarshal)
	//time.Sleep(1 * time.Second)
	reply, err = ptypes.MarshalAny(&pd.Reply{Code: 200, Message: unmarshal.Param})
	return
}

// QueryProto2 用于测试用户查询的方法
func (t *Test) QueryProto2(ctx context.Context, arg *anypb.Any) (reply *anypb.Any, err error) {
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", arg)
	unmarshal := &pd.Args2{}
	ptypes.UnmarshalAny(arg, unmarshal)
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", unmarshal)
	//time.Sleep(1 * time.Second)
	reply, err = ptypes.MarshalAny(&pd.Reply{Code: 200, Message: unmarshal.Param})
	return
}
