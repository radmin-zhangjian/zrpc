package example

import (
	"context"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	pb "zrpc/rpc/proto/proto/user"
)

// GetUserList 用于测试用户查询的方法
func (t *Test) GetUserList(ctx context.Context, arg *anypb.Any) (reply *anypb.Any, err error) {
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", arg)
	unmarshal := &pb.Args2{}
	anypb.UnmarshalTo(arg, unmarshal, proto.UnmarshalOptions{})
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", unmarshal)
	//time.Sleep(1 * time.Second)
	user := pb.Reply{
		Code:    1001,
		Message: unmarshal.Param,
		Gender:  pb.Reply_FEMALE,
		List: []*pb.Reply_Data{
			{Model: "iPhone12", Brand: "Apple.Inc"},
			{Model: "Mate40", Brand: "Huawei"},
			{Model: "S21", Brand: "Samsung"},
		},
		Detail: map[string]string{
			"身高": "180CM",
			"体重": "75KG",
			"爱好": "无",
		},
	}
	reply, err = anypb.New(&user)
	return
}
