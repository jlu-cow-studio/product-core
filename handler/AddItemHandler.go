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

func (h *Handler) AddItem(ctx context.Context, req *product_core.AddItemReq) (res *product_core.AddItemRes, err error) {
	res = &product_core.AddItemRes{
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
		log.Printf("[AddItem] Redis get token error: %v", cmd.Err())
		return
	}

	// 解析 token 中的用户信息
	userInfo := new(redis_model.UserInfo)
	if err = json.Unmarshal([]byte(cmd.Val()), userInfo); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "401"
		log.Printf("[AddItem] Unmarshal token error: %v", err)
		return
	}

	// 校验类别和角色匹配
	if !biz.CheckCategoryAndRole(req.GetItemInfo().GetCategory(), userInfo.Role) {
		res.Base.Message = "role and category not match"
		res.Base.Code = "402"
		log.Printf("[AddItem] Category and role not match, category: %v, role: %v", req.GetItemInfo().GetCategory(), userInfo.Role)
		return
	}

	// 添加商品到数据库
	itemInfo := req.GetItemInfo()
	uid, _ := strconv.Atoi(userInfo.Uid)
	item := &mysql_model.Item{
		Name:         itemInfo.GetName(),
		Description:  itemInfo.GetDescription(),
		Category:     itemInfo.GetCategory(),
		Price:        itemInfo.GetPrice(),
		Stock:        itemInfo.GetStock(),
		Province:     itemInfo.GetProvince(),
		City:         itemInfo.GetCity(),
		District:     itemInfo.GetDistrict(),
		UserID:       int32(uid),
		UserType:     userInfo.Role,
		SpecificAttr: itemInfo.GetSpecificAttributes(),
	}

	if err = biz.InsertItem(item); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "403"
		log.Printf("[AddItem] Insert item error: %v", err)
		return
	}

	res.ItemId = item.ID
	res.Base.Message = ""
	res.Base.Code = "200"
	log.Printf("[AddItem] Add item success, item ID: %v", item.ID)
	return
}
