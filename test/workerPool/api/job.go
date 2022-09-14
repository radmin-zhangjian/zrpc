package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
	"zrpc/utils/workerPool"
)

type Payload struct {
	num int64
	id  int64
}

func Response(code string, msg string, data any) any {
	result := make(map[string]any)
	result["code"] = code
	result["msg"] = msg
	result["data"] = data
	return result
}

var i int64 = 0

func To(c *gin.Context) {
	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")

	i++
	payload := Payload{num: i, id: 2}
	work := workerPool.Job{Payload: payload.Task}
	workerPool.JobQueue <- work

	c.JSON(http.StatusOK, Response("200", "success", nil))
}

// Task 耗时任务
func (p Payload) Task() (err error) {

	//log.Printf("id: %v", p.id)
	time.Sleep(10 * time.Second)
	log.Printf("num: %v", p.num)

	return
}
