package handler

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/jlu-cow-studio/product-core/biz"
)

func (h *Handler) DeleteItem(ctx context.Context, req *product_core.DeleteItemReq) (res *product_core.DeleteItemRes, err error) {
	res = &product_core.DeleteItemRes{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
	}

	// 获取 token
	token := req.Base.Token
	cmd := redis.DB.Get(redis.GetUserTokenKey(token))
	if cmd.Err() != nil {
		res.Base.Message = cmd.Err().Error()
		res.Base.Code = "400"
		log.Printf("[DeleteItem] Redis get token error: %v", cmd.Err())
		return
	}

	// 解析 token 中的用户信息
	userInfo := new(redis_model.UserInfo)
	if err = json.Unmarshal([]byte(cmd.Val()), userInfo); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "401"
		log.Printf("[DeleteItem] Unmarshal token error: %v", err)
		return
	}

	// 检查商品的拥有者是否匹配当前用户
	if !biz.CheckItemExsit(strconv.Itoa(int(req.GetItemId()))) {
		res.Base.Message = "item does not exsit"
		res.Base.Code = "402"
		log.Printf("[DeleteItem] item does not exsit")
		return
	}

	// 检查商品的拥有者是否匹配当前用户
	if !biz.CheckItemOwner(strconv.Itoa(int(req.GetItemId())), userInfo.Uid) {
		res.Base.Message = "role and user not match"
		res.Base.Code = "403"
		log.Printf("[DeleteItem] Role and user not match")
		return
	}

	// 构建待删除的商品对象
	item := &mysql_model.Item{
		ID: req.ItemId, // 改正变量名错误
	}
	if err = biz.DeleteItem(item); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "404" // 修改状态码
		log.Printf("[DeleteItem] Delete item error: %v", err)
		return
	}

	if err := biz.SendItemDeleteMsg(item.ToRedis()); err != nil {
		log.Fatalln("send delete message failed! ", err)
	}

	res.Base.Message = ""
	res.Base.Code = "200"
	return
}
