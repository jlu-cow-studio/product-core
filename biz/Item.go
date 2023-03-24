package biz

import (
	"github.com/jlu-cow-studio/common/dal/mysql"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
)

func CheckCategoryAndRole(catagory, role string) bool {
	switch role {
	case mysql_model.RoleBreeder:
		return catagory == mysql_model.CategoryBreeding ||
			catagory == mysql_model.CategoryCattleProduct ||
			catagory == mysql_model.CategoryWholeCattle
	case mysql_model.RoleServiceProvider:
		return catagory == mysql_model.CategoryServiceProduct ||
			catagory == mysql_model.CategoryService
	case mysql_model.RoleNormal:
		return false
	default:
		return false
	}
}

func InsertItem(item *mysql_model.Item) error {

	conn := mysql.GetDBConn()
	return conn.Table("items").Create(item).Error
}
