package example

import (
	"context"
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	test "zrpc/rpc/proto/proto/test"
)

type Test struct {
}

// QueryProto 用于测试用户查询的方法
func (t *Test) QueryProto(ctx context.Context, arg *anypb.Any) (reply *anypb.Any, err error) {
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", arg)
	unmarshal := &test.Args{}
	anypb.UnmarshalTo(arg, unmarshal, proto.UnmarshalOptions{})
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", unmarshal)
	//time.Sleep(1 * time.Second)
	reply, err = anypb.New(&test.Reply{Code: 200, Message: unmarshal.Param})
	return
}

// QueryProto2 用于测试用户查询的方法
func (t *Test) QueryProto2(ctx context.Context, arg *anypb.Any) (reply *anypb.Any, err error) {
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", arg)
	unmarshal := &test.Args2{}
	anypb.UnmarshalTo(arg, unmarshal, proto.UnmarshalOptions{})
	//log.Printf("service.QueryProto ===================== service.QueryProto：%v", unmarshal)
	//time.Sleep(1 * time.Second)

	res := map[string]interface{}{
		"a": "a111",
		"b": 2,
		"c": 3,
	}
	res2 := map[string]interface{}{
		"a": "a222",
		"b": 22,
		"c": 33,
	}
	resArr := make([]map[string]interface{}, 0)
	resArr = append(resArr, res, res2)
	strbyte, _ := json.Marshal(resArr)
	list := anypb.Any{Value: strbyte}
	data := map[string]*anypb.Any{"a": {TypeUrl: "a", Value: []byte("1")}, "b": {Value: []byte("app")}}
	reply, err = anypb.New(&test.Reply{Code: 200, Message: unmarshal.Param, Data: data, List: &list})
	return
}
