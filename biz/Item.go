package biz

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jlu-cow-studio/common/dal/mq"
	"github.com/jlu-cow-studio/common/dal/mysql"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/segmentio/kafka-go"
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

func UpdateItem(item *mysql_model.Item) error {
	return mysql.GetDBConn().Table("items").Where("id = ?", item.ID).UpdateColumns(item).Error
}

func SendItemUpdateMsg(item *redis_model.Item) error {
	writer := mq.GetWritter(mq.Topic_ItemChange)

	val, err := json.Marshal(item)
	if err != nil {
		return nil
	}

	msg := kafka.Message{
		Value: val,
	}
	err = writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}
	return err
}
