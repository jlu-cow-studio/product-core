package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	"github.com/jlu-cow-studio/common/dal/rpc/tag_core"
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

	log.Printf("[AddItem] request: %v", req)

	// 获取 token
	token := req.Base.Token
	cmd := redis.DB.Get(redis.GetUserTokenKey(token))
	if cmd.Err() != nil {
		res.Base.Message = cmd.Err().Error()
		res.Base.Code = "400"
		log.Printf("[AddItem] Redis get token error: %v", cmd.Err())
		return
	}
	log.Printf("[AddItem] Redis get token success, token: %s", token)

	// 解析 token 中的用户信息
	userInfo := new(redis_model.UserInfo)
	if err = json.Unmarshal([]byte(cmd.Val()), userInfo); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "401"
		log.Printf("[AddItem] Unmarshal token error: %v", err)
		return
	}
	log.Printf("[AddItem] Unmarshal token success, userinfo: %v", userInfo)

	// 校验类别和角色匹配
	if !biz.CheckCategoryAndRole(req.GetItemInfo().GetCategory(), userInfo.Role) {
		res.Base.Message = "role and category not match"
		res.Base.Code = "402"
		log.Printf("[AddItem] Category and role not match, category: %v, role: %v", req.GetItemInfo().GetCategory(), userInfo.Role)
		return
	}
	log.Printf("[AddItem] CheckCategoryAndRole success, category: %s, role: %s", req.GetItemInfo().GetCategory(), userInfo.Role)

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

	tx := biz.InsertItem(item)
	if err = tx.Error; err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "403"
		log.Printf("[AddItem] Insert item error: %v", err)
		tx.Rollback()
		return
	}
	log.Printf("[AddItem] Insert item success, item: %v", item)

	if err := biz.SendItemCreateMsg(ctx, item.ToRedis()); err != nil {
		log.Fatalf("[AddItem] Send create message failed! error: %v", err)
		tx.Rollback()
	}
	log.Printf("[AddItem] Send create message success, item: %v", item)

	cli, err := rpc.GetTagCoreCli()
	if err != nil {
		log.Printf("get rpc conn error: %s\n", err.Error())
		res.Base.Message = err.Error()
		res.Base.Code = "405"
		return
	}

	tagUpdateItemTagsReq := &tag_core.UpdateItemTagsRequest{
		Base:    req.Base,
		TagList: req.TagList,
		ItemId:  strconv.FormatInt(int64(item.ID), 10),
	}

	tagUpdateItemTagsRes, err := cli.UpdateItemTags(ctx, tagUpdateItemTagsReq)
	if err != nil {
		res.Base.Message = fmt.Sprintf("error when update tag list: %v", err.Error())
		res.Base.Code = "406"
		tx.Rollback()
		return res, nil
	}

	if tagUpdateItemTagsRes.Base.Code != "200" {
		res.Base.Message = fmt.Sprintf("error when update tag list %v, %v", tagUpdateItemTagsRes.Base.Code, tagUpdateItemTagsRes.Base.Message)
		res.Base.Code = "407"
		tx.Rollback()
		return res, nil
	}

	tx.Commit()

	res.ItemId = item.ID
	res.Base.Message = ""
	res.Base.Code = "200"
	log.Printf("[AddItem] Add item success, item ID: %v", item.ID)
	return
}
