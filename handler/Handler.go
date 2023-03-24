package handler

import "github.com/jlu-cow-studio/common/dal/rpc/product_core"

type Handler struct {
	product_core.UnimplementedProductCoreServiceServer
}
