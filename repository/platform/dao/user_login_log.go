package dao

import (
	"errors"
	"server-api/global"
	"server-api/repository/platform"
	"server-api/repository/types/platformfilter"

	"github.com/gomooth/pkg/framework/dbquery"

	"github.com/save95/xerror"
	"github.com/save95/xerror/xcode"

	"gorm.io/gorm"
)

type userLoginLog struct {
	db *gorm.DB
}

func NewUserLoginLog(options ...interface{}) *userLoginLog {
	impl := userLoginLog{}
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

func (u *userLoginLog) Create(record *platform.UserLoginLog) error {
	if record.ID > 0 {
		return xerror.WithXCode(xcode.DBRequestParamError)
	}

	if err := u.db.Create(record).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *userLoginLog) Save(record *platform.UserLoginLog) error {
	if record.ID == 0 {
		return xerror.WithXCode(xcode.DBRecordNotFound)
	}

	if err := u.db.Save(record).Error; nil != err {
		return xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return nil
}

func (u *userLoginLog) First(id uint) (*platform.UserLoginLog, error) {
	if id == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	var record platform.UserLoginLog
	if err := u.db.Where("id = ?", id).First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}
		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

func (u *userLoginLog) FirstByUser(userID uint) (*platform.UserLoginLog, error) {
	if userID == 0 {
		return nil, xerror.WithXCode(xcode.DBRequestParamError)
	}

	db := u.db.Model(platform.UserLoginLog{}).
		Where("user_id = ?", userID)

	var record platform.UserLoginLog
	if err := db.First(&record).Error; nil != err {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, xerror.WithXCode(xcode.DBRecordNotFound)
		}

		return nil, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return &record, nil
}

func (u *userLoginLog) Paginate(start, limit int, option dbquery.IFilter[platformfilter.UserLoginLog]) ([]*platform.UserLoginLog, uint, error) {
	db := u.db.Model(platform.UserLoginLog{})

	filter := option.Filter()

	if v := filter.UserID; v > 0 {
		db = db.Where("user_id = ?", v)
	}

	var total int64
	_ = db.Count(&total).Error

	db = option.Build(db,
		dbquery.BuildWithPage(start, limit),
	)

	var records []*platform.UserLoginLog
	if err := db.Find(&records).Error; nil != err {
		return nil, 0, xerror.WrapWithXCode(err, xcode.DBFailed)
	}

	return records, uint(total), nil
}
