package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"log"
	"strconv"
	"zrpc/rpc/zrpc"
)

type Test struct {
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
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// QueryUser 用于测试用户查询的方法
func (t *Test) QueryUser(ctx context.Context, arg Args, reply *any, error *error) error {
	user := make(map[int64]User)
	// 假数据
	user[0] = User{"zs", 20}
	user[1] = User{"ls", 21}
	user[2] = User{"ww", 22}
	//log.Printf("v1.QueryUser ===================== v1.QueryUser ID：%v, X: %v", arg.Id, arg.X)
	// 模拟查询用户
	uid := arg.Id
	if u, ok := user[uid]; ok {
		//rep, _ := json.Marshal(u)
		//*reply = string(rep)
		*reply = u
	} else {
		return errors.New("No data found. uid = " + strconv.FormatInt(uid, 10))
	}
	return nil
}

// QueryInt 用于测试用户查询的方法
func (t *Test) QueryInt(ctx context.Context, arg map[string]any, reply *any, error *error) error {
	log.Printf("v1.QueryInt ===================== v1.QueryInt ID：%v", arg["id"])
	//time.Sleep(8 * time.Second)

	if arg["id"].(uint16) == 10000 {
		*reply = "111222333" + "::" + arg["msg"].(string)
	} else {
		*error = errors.New("QueryInt" + "===::===" + arg["msg"].(string))
	}
	return nil
}

type QueryIntC struct {
	Id      int      `json:"id"`
	Msg     string   `json:"msg"`
	Address []string `json:"address"`
}

// QueryIntC 用于测试用户查询的方法
func (t *Test) QueryIntC(c *zrpc.Context) error {
	arg := QueryIntC{}
	err := mapstructure.Decode(c.Args, &arg)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Printf("v1.QueryIntC ===================== v1.QueryIntC ID：%v", arg.Id)
	if arg.Id == 10000 {
		*c.Reply = "QueryIntC" + "000::000" + arg.Msg
	} else {
		*c.Error = errors.New("QueryIntC" + "===::===" + arg.Msg)
	}
	return nil
}
