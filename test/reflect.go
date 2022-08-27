package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// 通过反射获取方法的结构体参数并赋值
func main() {

	service := new(Test)
	rType := reflect.TypeOf(service)
	rValue := reflect.ValueOf(service)

	method, _ := rType.MethodByName("QueryUser")

	// 遍历解析客户端传来的参数,放切片里
	mNumIn := method.Type.NumIn()
	if mNumIn != 4 {
		log.Println("Err NumIn Lens Not Equal")
		return
	}

	// map
	argMap := map[string]any{
		"Id": 1,
		"X":  2,
		"Y":  3,
		"Z":  "aassdd",
	}

	// Args reflect.Type类型
	t := method.Type.In(2)
	// Args reflect.Value类型
	v := reflect.New(t)
	//var inArgs []any
	if t.Kind() == reflect.Struct {
		argNum := t.NumField()
		//inArgs = make([]any, 0, argNum)
		for i := 0; i < argNum; i++ {
			//inArgs = append(inArgs, t.Field(i))
			// 获取指定字段的反射值
			f := v.Elem().Field(i)
			// 获取struct的指定字段
			stf := t.Field(i)
			// 获取tag
			name := strings.Split(stf.Tag.Get("json"), ",")[0]
			// 判断是否为忽略字段
			if name == "-" {
				continue
			}
			// 判断是否为空，若为空则使用字段本身的名称获取value值
			if name == "" {
				name = stf.Name
			}
			// 获取value值
			val, ok := argMap[name]
			if !ok {
				continue
			}
			// 获取指定字段的类型
			kind := v.Elem().Field(i).Kind()
			// 若字段为指针类型
			if kind == reflect.Ptr {
				// 获取对应字段的kind
				kind = f.Type().Elem().Kind()
			}
			switch kind {
			case reflect.Int:
			case reflect.Int64:
				res, _ := strconv.ParseInt(fmt.Sprint(val), 10, 64)
				v.Elem().Field(i).SetInt(res)
			case reflect.String:
				v.Elem().Field(i).SetString(fmt.Sprint(val))
			}
		}
	}
	log.Println(t)
	log.Printf("%#v", v)              // 指针
	log.Printf("Args: %#v", v.Elem()) // 值
	//return
	//mapArgs := GetArgsRValue(method, rpcData.MapArgs)
	mapArgs := v.Elem()
	inArgs := make([]reflect.Value, 0, 4)
	inArgs = append(inArgs, rValue)
	inArgs = append(inArgs, reflect.ValueOf(context.Background()))
	inArgs = append(inArgs, mapArgs)
	var Reply any
	inArgs = append(inArgs, reflect.ValueOf(reflect.ValueOf(&Reply).Interface()))
	// 反射调用方法
	// 返回Value类型，用于给客户端传递返回结果,out是所有的返回结果
	out := method.Func.Call(inArgs)
	log.Printf("out: %#v", out)
	return
}

type Test struct {
}

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
func (t *Test) QueryUser(ctx context.Context, arg Args, reply *any) any {
	user := make(map[int64]User)
	// 假数据
	user[0] = User{"zs", 20}
	user[1] = User{"ls", 21}
	user[2] = User{"ww", 22}
	log.Println("QueryUser ===================== QueryUser ID：", arg.Id)
	// 模拟查询用户
	uid := arg.Id
	var rep []byte
	if u, ok := user[uid]; ok {
		rep, _ = json.Marshal(u)
		*reply = u
	} else {
		a := map[string]any{"Name": "qqq", "Age": 23}
		rep, _ = json.Marshal(a)
		*reply = a
	}
	return rep
}
