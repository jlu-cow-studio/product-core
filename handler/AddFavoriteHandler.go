package handler

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/jlu-cow-studio/common/model/http_struct/item"
	"github.com/jlu-cow-studio/product-core/biz"
)

func (h *Handler) AddFavorite(ctx context.Context, req *product_core.AddFavoriteReq) (res *product_core.AddFavoriteRes, err error) {
	res = &product_core.AddFavoriteRes{
		Base: &base.BaseRes{
			Message: "",
			Code:    "498",
		},
	}

	cmd := redis.DB.Get(redis.GetUserTokenKey(req.Base.Token))
	if cmd.Err() != nil {
		res.Base.Message = cmd.Err().Error()
		res.Base.Code = "401"
		return res, nil
	}

	info := &redis_model.UserInfo{}

	if err := json.Unmarshal([]byte(cmd.Val()), info); err != nil {
		res.Base.Message = err.Error()
		res.Base.Code = "402"
		return res, nil
	}

	if req.Action == item.AddFavoriteAction_Add {
		if ok, err := biz.CheckFavoriteAdded(info.Uid, strconv.Itoa(int(req.ItemId))); err != nil {
			res.Base.Message = err.Error()
			res.Base.Code = "403"
			return res, nil
		} else if ok {
			res.Base.Message = "favorite already added"
			res.Base.Code = "404"
			return res, nil
		} else if err := biz.AddFavorite(info.Uid, strconv.Itoa(int(req.ItemId))); err != nil {
			res.Base.Message = err.Error()
			res.Base.Code = "405"
			return res, nil
		}
	} else if req.Action == item.AddFovoriteAction_Del {
		if ok, err := biz.CheckFavoriteAdded(info.Uid, strconv.Itoa(int(req.ItemId))); err != nil {
			res.Base.Message = err.Error()
			res.Base.Code = "406"
			return res, nil
		} else if !ok {
			res.Base.Message = "favorite already added"
			res.Base.Code = "407"
			return res, nil
		} else if err := biz.DelFavorite(info.Uid, strconv.Itoa(int(req.ItemId))); err != nil {
			res.Base.Message = err.Error()
			res.Base.Code = "408"
			return res, nil
		}
	} else {
		res.Base.Message = "unknown action"
		res.Base.Code = "400"
		return res, nil
	}

	return
}
