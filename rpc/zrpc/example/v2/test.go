package v2

import (
	"context"
	"errors"
	"log"
	"time"
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
	Name string
	Age  int
}

// QueryUser 用于测试用户查询的方法
func (t *Test) QueryUser(ctx context.Context, arg Args, reply *any, error *error) error {
	user := make(map[int64]User)
	// 假数据
	user[0] = User{"zs", 20}
	user[1] = User{"ls", 21}
	user[2] = User{"ww", 22}
	log.Printf("v2.QueryUser ===================== v2.QueryUser ID：%v, X: %v, Z: %v", arg.Id, arg.X, arg.Z)
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
func (t *Test) QueryInt(ctx context.Context, arg map[string]any, reply *any, error *error) error {
	log.Printf("v2.QueryInt ===================== v2.QueryInt ID：%v", arg["id"])
	time.Sleep(1 * time.Second)
	*reply = 111222333
	return nil
}
