package pdao

import (
	"context"
	"errors"
	"server-api/global"
	"server-api/repository/platform/pattr"
	"server-api/repository/platform/pmodel"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/gomooth/xerror"
	"github.com/gomooth/xerror/xcode"
)

type user struct {
	db *gorm.DB
}

type IUser interface {
	Create(ctx context.Context, genres []int8, record *pmodel.User, stat *pmodel.UserStat) error
	Update(ctx context.Context, record *pmodel.User, genres []int8) error
	CreateAndBindThirdUser(ctx context.Context, genres []int8, record *pmodel.User, stat *pmodel.UserStat, thirdUserID uint) error
	Save(ctx context.Context, record *pmodel.User) error
}

func NewUser() IUser {
	db, _ := global.Database().Get("platform")
	return &user{
		db: db,
	}
}

func (u *user) Create(ctx context.Context, genres []int8, record *pmodel.User, stat *pmodel.UserStat) error {
	if record.ID > 0 || len(genres) == 0 {
		return xerror.NewXCode(xcode.DBRequestParamError)
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := tx.Create(record).Error; err != nil {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写统计
		stat.UserID = record.ID
		if err := tx.Create(stat).Error; err != nil {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写角色
		roles := make([]*pmodel.UserRole, 0)
		for _, genre := range genres {
			roles = append(roles, &pmodel.UserRole{
				Genre:  uint8(genre),
				UserID: record.ID,
			})
		}
		if err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(roles, 100).Error; err != nil {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		return nil
	})
}

func (u *user) Update(ctx context.Context, record *pmodel.User, genres []int8) error {
	if record.ID == 0 {
		return xerror.NewXCode(xcode.DBRecordNotFound)
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; err != nil {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 清空角色
		if err := tx.Model(pmodel.UserRole{}).Unscoped().
			Where("user_id = ?", record.ID).Delete(&pmodel.UserRole{}).Error; err != nil {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写角色
		roles := make([]*pmodel.UserRole, 0)
		for _, genre := range genres {
			roles = append(roles, &pmodel.UserRole{
				Genre:  uint8(genre),
				UserID: record.ID,
			})
		}
		if err := tx.CreateInBatches(roles, 100).Error; err != nil {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		return nil
	})
}

func (u *user) CreateAndBindThirdUser(ctx context.Context, genres []int8, record *pmodel.User, stat *pmodel.UserStat, thirdUserID uint) error {
	if record.ID > 0 {
		return xerror.NewXCode(xcode.DBRequestParamError)
	}

	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 查找已有记录（含软删除）：若 account 曾被软删除则恢复，否则新建
		var existing pmodel.User
		err := tx.Unscoped().Where("account = ?", record.Account).First(&existing).Error
		if err == nil {
			record.ID = existing.ID
			record.CreatedAt = existing.CreatedAt
			if err := tx.Unscoped().Model(&pmodel.User{}).
				Where("id = ?", existing.ID).
				Updates(map[string]interface{}{
					"nickname":   record.Nickname,
					"password":   record.Password,
					"deleted_at": nil, // 清空删除状态
				}).Error; err != nil {
				return xerror.WrapWithXCode(err, xcode.DBFailed)
			}
		} else if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tx.Create(record).Error; err != nil {
				return xerror.WrapWithXCode(err, xcode.DBFailed)
			}
		} else {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写统计
		stat.UserID = record.ID
		if err := tx.Where("user_id = ?", record.ID).FirstOrCreate(stat).Error; err != nil {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写角色
		roles := make([]*pmodel.UserRole, 0)
		for _, genre := range genres {
			roles = append(roles, &pmodel.UserRole{
				Genre:  uint8(genre),
				UserID: record.ID,
			})
		}
		if err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).
			CreateInBatches(roles, 100).Error; err != nil {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 绑定第三方用户
		switch stat.FromPlatformID {
		case pattr.UserFromPlatformAccount:
			// skip
		case pattr.UserFromPlatformWechat:
			if err := tx.Model(pmodel.UserWechat{}).
				Where("id = ?", thirdUserID).
				Updates(map[string]interface{}{
					"user_id": record.ID,
					"bind_at": time.Now(),
				}).Error; err != nil {
				return xerror.WrapWithXCode(err, xcode.DBFailed)
			}
		case pattr.UserFromPlatformAlipay:
			if err := tx.Model(pmodel.UserAlipay{}).
				Where("id = ?", thirdUserID).
				Updates(map[string]interface{}{
					"user_id": record.ID,
					"bind_at": time.Now(),
				}).Error; err != nil {
				return xerror.WrapWithXCode(err, xcode.DBFailed)
			}
		default:
			return xerror.New("not support user from platform")
		}

		return nil
	})
}

func (u *user) Save(ctx context.Context, record *pmodel.User) error {
	if record.ID == 0 {
		return xerror.NewXCode(xcode.DBRecordNotFound)
	}

	if err := u.db.WithContext(ctx).Save(record).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}
	return nil
}
