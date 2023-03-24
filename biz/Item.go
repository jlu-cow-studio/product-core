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

func CheckItemExsit(ItemId string) bool {
	count := new(int64)
	if mysql.GetDBConn().Table("items").Where("id = ?", ItemId).Count(count).Error != nil {
		return false
	}
	return *count >= 1
}

func CheckItemOwner(ItemId, UId string) bool {
	count := new(int64)
	if mysql.GetDBConn().Table("items").Where("id = ?", ItemId).Where("user_id = ?", UId).Count(count).Error != nil {
		return false
	}
	return *count == 1
}

func DeleteItem(item *mysql_model.Item) error {
	return mysql.GetDBConn().Table("items").Delete(item).Error
}
