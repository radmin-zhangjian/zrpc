package v1

import (
	"context"
	"errors"
	"strconv"
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
func (t *Test) QueryUser(ctx context.Context, arg Args, reply *any) error {
	user := make(map[int64]User)
	// 假数据
	user[0] = User{"zs", 20}
	user[1] = User{"ls", 21}
	user[2] = User{"ww", 22}
	//log.Printf("v1.QueryUser ===================== v1.QueryUser ID：%v, X: %v", arg.Id, arg.X)
	// 模拟查询用户
	uid := arg.Id
	if u, ok := user[uid]; ok {
		//var rep []byte
		//rep, _ = json.Marshal(u)
		//*reply = rep
		*reply = u
	} else {
		return errors.New("No data found. uid = " + strconv.FormatInt(uid, 10))
	}
	return nil
}

// QueryInt 用于测试用户查询的方法
func (t *Test) QueryInt(ctx context.Context, arg map[string]any, reply *any) error {
	//log.Printf("v1.QueryInt ===================== v1.QueryInt ID：%v", arg["Id"])
	time.Sleep(5 * time.Second)
	*reply = "111222333" + "::" + arg["msg"].(string)
	return nil
}
