package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	"github.com/jlu-cow-studio/common/discovery"
	"github.com/sanity-io/litter"
)

func TestUpdateItem(t *testing.T) {
	// 初始化测试环境
	discovery.Init()
	redis.Init()
	mysql.Init()
	// 构造请求对象
	req := &product_core.UpdateItemReq{
		Base: &base.BaseReq{
			Token: "84dd7f33-f145-4160-bbd7-8a84226f1783",
		},
		Item: &product_core.ItemInfo{
			ItemId:             int32(27),
			Name:               "牛肉干",
			Description:        "美味牛肉干，口感鲜美，量大123管饱。",
			Category:           "cattle_product",
			Price:              28.00,
			Stock:              200,
			Province:           "湖南",
			City:               "长沙",
			District:           "岳麓",
			UserId:             3,
			UserType:           "service_provider",
			SpecificAttributes: `{"产地": "湖南", "口味": "五香味", "保质期": "180天"}`,
		},
	}

	// 创建 handler 对象
	h := &Handler{}

	// 调用 UpdateItem 函数
	res, err := h.UpdateItem(context.Background(), req)
	litter.Dump(res)
	if err != nil {
		t.Errorf("TestUpdateItem failed, err: %v", err)
	}
	if res.Base.Code != "200" {
		t.Errorf("TestUpdateItem failed, code: %v, message: %v", res.Base.Code, res.Base.Message)
	}
	fmt.Printf("TestUpdateItem succeed, res: %v\n", res)
}
