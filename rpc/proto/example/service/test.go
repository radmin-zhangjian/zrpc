package service

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/anypb"
	"log"
	"time"
	pd "zrpc/rpc/proto"
)

type Test struct {
}

type Reply struct {
}

// Args 参数
type Args struct {
	Id int64
	X  int64
	Y  int64
	Z  string
}

// User 定义用户对象
type User struct {
	Name string
	Age  int
}

// QueryUser 用于测试用户查询的方法
func (t *Test) QueryUser(ctx context.Context, arg Args, reply *any) error {
	user := make(map[int64]User)
	// 假数据
	user[0] = User{"zs", 20}
	user[1] = User{"ls", 21}
	user[2] = User{"ww", 22}
	//log.Printf("service.QueryUser ===================== service.QueryUser ID：%v, X: %v", arg.Id, arg.X)
	// 模拟查询用户
	uid := arg.Id
	if u, ok := user[uid]; ok {
		//var rep []byte
		//rep, _ = json.Marshal(u)
		//*reply = rep
		*reply = u
	} else {
		return errors.New("No data found")
	}
	return nil
}

// QueryInt 用于测试用户查询的方法
func (t *Test) QueryInt(ctx context.Context, arg map[string]any, reply *any) error {
	log.Printf("service.QueryInt ===================== service.QueryInt ID：%v", arg["Id"])
	time.Sleep(1 * time.Second)
	*reply = 111222333
	return nil
}

// QueryProto 用于测试用户查询的方法
func (t *Test) QueryProto(ctx context.Context, arg *anypb.Any, reply *anypb.Any) (*anypb.Any, error) {
	log.Printf("service.QueryProto ===================== service.QueryProto：%v", arg)
	args := arg
	unmarshal := &pd.Args{}
	ptypes.UnmarshalAny(args, unmarshal)
	log.Printf("args: %v", unmarshal)
	time.Sleep(1 * time.Second)
	reply, _ = ptypes.MarshalAny(&pd.Reply{Result: "aaasss"})
	//*reply = *inArgsAny
	return reply, nil
}
