package dao

import (
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/types/platformtypes"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/save95/xerror"
	"github.com/save95/xerror/xcode"
)

type user struct {
	db *gorm.DB
}

func NewUser(options ...interface{}) *user {
	impl := user{}
	for _, option := range options {
		if db, ok := option.(*gorm.DB); ok {
			impl.db = db
		}
	}

	if impl.db == nil {
		impl.db, _ = global.Database().Get("platform")
	}

	return &impl
}

func (u *user) Create(genres []int8, record *platform.User, stat *platform.UserStat) error {
	if record.ID > 0 || len(genres) == 0 {
		return xerror.WithXCode(xcode.DBRequestParamError)
	}

	return u.db.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := tx.Create(record).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写统计
		stat.UserID = record.ID
		if err := tx.Create(stat).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写角色
		roles := make([]*platform.UserRole, 0)
		for _, genre := range genres {
			roles = append(roles, &platform.UserRole{
				Genre:  uint8(genre),
				UserID: record.ID,
			})
		}
		if err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(roles, 100).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		return nil
	})
}

func (u *user) Update(record *platform.User, genres []int8) error {
	if record.ID == 0 {
		return xerror.WithXCode(xcode.DBRecordNotFound)
	}

	return u.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(record).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 清空角色
		if err := tx.Model(platform.UserRole{}).Unscoped().
			Where("user_id = ?").Delete(&platform.UserRole{}).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写角色
		roles := make([]*platform.UserRole, 0)
		for _, genre := range genres {
			roles = append(roles, &platform.UserRole{
				Genre:  uint8(genre),
				UserID: record.ID,
			})
		}
		if err := tx.CreateInBatches(roles, 100).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		return nil
	})
}

func (u *user) CreateAndBindThirdUser(genres []int8, record *platform.User, stat *platform.UserStat, thirdUserID uint) error {
	if record.ID > 0 {
		return xerror.WithXCode(xcode.DBRequestParamError)
	}

	return u.db.Transaction(func(tx *gorm.DB) error {
		// 如果存在则更新
		if err := tx.Where("account = ?", record.Account).
			Assign(map[string]interface{}{
				"nickname":   record.Nickname,
				"password":   record.Password,
				"deleted_at": nil, // 清空删除状态
			}).
			FirstOrCreate(record).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写统计
		stat.UserID = record.ID
		if err := tx.FirstOrCreate(stat).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 写角色
		roles := make([]*platform.UserRole, 0)
		for _, genre := range genres {
			roles = append(roles, &platform.UserRole{
				Genre:  uint8(genre),
				UserID: record.ID,
			})
		}
		if err := tx.Clauses(clause.Insert{Modifier: "IGNORE"}).
			CreateInBatches(roles, 100).Error; nil != err {
			return xerror.WrapWithXCode(err, xcode.DBFailed)
		}

		// 绑定第三方用户
		switch stat.FromPlatformID {
		case platformtypes.UserFromPlatformAccount:
			// skip
		case platformtypes.UserFromPlatformWechat:
			if err := tx.Model(platform.UserWechat{}).
				Where("id = ?", thirdUserID).
				Updates(map[string]interface{}{
					"user_id": record.ID,
					"bind_at": time.Now(),
				}).Error; nil != err {
				return xerror.WrapWithXCode(err, xcode.DBFailed)
			}
		case platformtypes.UserFromPlatformAlipay:
			if err := tx.Model(platform.UserAlipay{}).
				Where("id = ?", thirdUserID).
				Updates(map[string]interface{}{
					"user_id": record.ID,
					"bind_at": time.Now(),
				}).Error; nil != err {
				return xerror.WrapWithXCode(err, xcode.DBFailed)
			}
		default:
			return xerror.New("not support user from platform")
		}

		return nil
	})
}

func (u *user) Save(record *platform.User) error {
	if record.ID == 0 {
		return xerror.WithXCode(xcode.DBRecordNotFound)
	}

	if err := u.db.Save(record).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}
