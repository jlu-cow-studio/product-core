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
	"github.com/sanity-io/litter"
)

func (h *Handler) UpdateItem(ctx context.Context, req *product_core.UpdateItemReq) (res *product_core.UpdateItemRes, err error) {

	res = &product_core.UpdateItemRes{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
	}

	log.Println(litter.Sdump(req))

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

	itemId := strconv.Itoa(int(req.GetItem().ItemId))

	// 检查商品的拥有者是否匹配当前用户
	if !biz.CheckItemExsit(itemId) {
		res.Base.Message = "item does not exsit"
		res.Base.Code = "402"
		log.Printf("[DeleteItem] item does not exsit")
		return
	}

	// 检查商品的拥有者是否匹配当前用户
	if !biz.CheckItemOwner(itemId, userInfo.Uid) {
		res.Base.Message = "role and user not match"
		res.Base.Code = "403"
		log.Printf("[DeleteItem] Role and user not match")
		return
	}

	item := req.GetItem()
	updateItem := &mysql_model.Item{
		ID:           item.ItemId,
		Name:         item.GetName(),
		Description:  item.GetDescription(),
		Price:        item.GetPrice(),
		Stock:        item.GetStock(),
		Province:     item.GetProvince(),
		City:         item.GetCity(),
		District:     item.GetDistrict(),
		UserType:     item.GetUserType(),
		SpecificAttr: item.GetSpecificAttributes(),
	}

	tx := biz.UpdateItem(updateItem)
	if tx.Error != nil {
		res.Base.Message = tx.Error.Error()
		res.Base.Code = "404" // 修改状态码
		log.Printf("[DeleteItem] update item error: %v", err)
		return
	}

	if err := biz.SendItemUpdateMsg(ctx, updateItem.ToRedis()); err != nil {
		log.Fatalln("send update message failed! ", err)
		tx.Rollback()
	}

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
		ItemId:  strconv.FormatInt(int64(req.Item.ItemId), 10),
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

	res.Base.Message = ""
	res.Base.Code = "200"
	return
}
