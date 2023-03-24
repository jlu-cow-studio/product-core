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
	"github.com/stretchr/testify/assert"
)

func TestAddItem(t *testing.T) {
	discovery.Init()
	redis.Init()
	mysql.Init()
	// 构造测试请求
	req := &product_core.AddItemReq{
		Base: &base.BaseReq{
			Token: "84dd7f33-f145-4160-bbd7-8a84226f1783",
		},
		ItemInfo: &product_core.ItemInfo{
			Name:               "牛肉干",
			Description:        "美味牛肉干，口感鲜美。",
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

	// 创建测试处理器
	handler := &Handler{}

	// 调用处理器方法
	res, err := handler.AddItem(context.Background(), req)
	fmt.Println(litter.Sdump(res), err)

	// 断言响应结果
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "200", res.Base.Code)
	assert.True(t, res.ItemId > 0)
}
