package handler

import (
	"context"
	"encoding/json"
	auth "github.com/micro-in-cn/tutorials/microservice-in-micro/part3/auth/proto/auth"
	invS "github.com/micro-in-cn/tutorials/microservice-in-micro/part3/inventory-srv/proto/inventory"
	orders "github.com/micro-in-cn/tutorials/microservice-in-micro/part3/orders-srv/proto/orders"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part3/plugins/session"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/client"
	"net/http"
	"strconv"
	"time"
)

var (
	serviceClient orders.OrdersService
	authClient    auth.Service
	invClient     invS.InventoryService
)

// Error 错误结构体
type Error struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func Init() {
	serviceClient = orders.NewOrdersService("mu.micro.book.srv.orders", client.DefaultClient)
	authClient = auth.NewService("mu.micro.book.srv.auth", client.DefaultClient)
}

// New 新增订单入口
func New(w http.ResponseWriter, r *http.Request) {

	// 只接受POST请求
	if r.Method != "POST" {
		log.Logf("非法请求")
		http.Error(w, "非法请求", 400)
		return
	}

	r.ParseForm()

	bookId, _ := strconv.ParseInt(r.Form.Get("userName"), 64, 10)

	// 调用后台服务
	rsp, err := serviceClient.New(context.TODO(), &orders.Request{
		BookId: bookId,
		UserId: session.GetSession(w, r).Values["userId"].(int64),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 返回结果
	response := map[string]interface{}{
		"orderId": rsp.Order.Id,
		"ref":     time.Now().UnixNano(),
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// 返回JSON结构
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
