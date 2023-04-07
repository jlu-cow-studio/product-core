package biz

import (
	"context"
	"log"
	"time"

	"github.com/jlu-cow-studio/common/dal/mq"
	"github.com/jlu-cow-studio/common/dal/mysql"
	mysql_model "github.com/jlu-cow-studio/common/model/dao_struct/mysql"
	redis_model "github.com/jlu-cow-studio/common/model/dao_struct/redis"
	"github.com/jlu-cow-studio/common/model/mq_struct"
	"gorm.io/gorm"
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

func InsertItem(item *mysql_model.Item) *gorm.DB {
	tx := mysql.GetDBConn().Begin()
	tx.Table("items").Create(item)
	return tx
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

func UpdateItem(item *mysql_model.Item) *gorm.DB {
	tx := mysql.GetDBConn().Begin()
	tx.Table("items").Where("id = ?", item.ID).UpdateColumns(item)
	return tx
}

func SendItemUpdateMsg(ctx context.Context, item *redis_model.Item) (err error) {
	if err = mq.SendMessage(ctx, mq.Topic_ItemChange, &mq_struct.ItemChangeMsg{
		Op:   mq_struct.ItemOp_Update,
		Info: item,
	}); err != nil {
		log.Fatal("failed to write messages:", err)
	}
	return err
}

func SendItemCreateMsg(ctx context.Context, item *redis_model.Item) (err error) {

	if err = mq.SendMessage(ctx, mq.Topic_ItemChange, &mq_struct.ItemChangeMsg{
		Op:   mq_struct.ItemOp_Create,
		Info: item,
	}); err != nil {
		log.Fatal("failed to write messages:", err)
	}
	return err
}

func SendItemDeleteMsg(ctx context.Context, item *redis_model.Item) (err error) {
	if err = mq.SendMessage(ctx, mq.Topic_ItemChange, &mq_struct.ItemChangeMsg{
		Op:   mq_struct.ItemOp_Delete,
		Info: item,
	}); err != nil {
		log.Fatal("failed to write messages:", err)
	}
	return err
}

func AddFavorite(userId, itemId string) error {

	return mysql.GetDBConn().Table("user_item_follow").Create(&struct {
		UserId   string    `gorm:"column:user_id"`
		ItemId   string    `gorm:"column:item_id"`
		CreateAt time.Time `gorm:"create_at"`
	}{
		UserId:   userId,
		ItemId:   itemId,
		CreateAt: time.Now(),
	}).Error
}

func DelFavorite(userId, itemId string) error {

	return mysql.GetDBConn().Table("user_item_follow").Delete(nil, "user_id = ? and item_id = ?").Error
}

func CheckFavoriteAdded(userId, itemId string) (bool, error) {

	count := new(int64)

	tx := mysql.GetDBConn().Table("user_item_follow").Where("user_id = ? and item_id = ?", userId, itemId).Count(count)

	return *count == 1, tx.Error
}

func SendAddFavoriteMsg(ctx context.Context, userId, itemId string) error {

	return mq.SendMessage(ctx, mq.Topic_UserAction, &mq_struct.UserActionMsg{
		Op: mq_struct.UserActionOp_AddFavorite,
		Extra: map[string]string{
			mq_struct.UserActionExtraKey_ItemId: itemId,
			mq_struct.UserActionExtraKey_UserId: userId,
		},
	})
}
func SendDelFavoriteMsg(ctx context.Context, userId, itemId string) error {

	return mq.SendMessage(ctx, mq.Topic_UserAction, &mq_struct.UserActionMsg{
		Op: mq_struct.UserActionOp_DelFavorite,
		Extra: map[string]string{
			mq_struct.UserActionExtraKey_ItemId: itemId,
			mq_struct.UserActionExtraKey_UserId: userId,
		},
	})
}
