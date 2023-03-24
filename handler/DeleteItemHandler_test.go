package handler

import (
	"context"
	"testing"

	"github.com/jlu-cow-studio/common/dal/mysql"
	"github.com/jlu-cow-studio/common/dal/redis"
	"github.com/jlu-cow-studio/common/dal/rpc/base"
	"github.com/jlu-cow-studio/common/dal/rpc/product_core"
	"github.com/jlu-cow-studio/common/discovery"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
)

const ItemId = 25

func TestDeleteItem(t *testing.T) {
	// 初始化测试环境
	discovery.Init()
	redis.Init()
	mysql.Init()

	// 构造请求参数
	req := &product_core.DeleteItemReq{
		Base: &base.BaseReq{
			Token: "84dd7f33-f145-4160-bbd7-8a84226f1783",
		},
		ItemId: ItemId,
	}

	h := &Handler{}
	// 调用接口
	res, err := h.DeleteItem(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Base.Code != "200" {
		t.Fatalf("unexpected response code: %v", res.Base.Code)
	}

	// 验证商品是否已经被删除
	var item mysql_model.Item
	if err := mysql.GetDBConn().Table("items").Where("id = ?", ItemId).First(&item).Error; err == nil {
		t.Fatalf("item not deleted: %v", item)
	}
}
