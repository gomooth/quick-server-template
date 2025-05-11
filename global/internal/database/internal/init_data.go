package internal

import (
	"server-api/global/internal/database/internal/platform"

	"github.com/gomooth/pkg/framework/dbmanager"
	"github.com/gomooth/utils/userutil"

	"github.com/pkg/errors"
)

type initData struct {
	dbs dbmanager.IDatabaseManager
}

func NewInit(dbs dbmanager.IDatabaseManager) *initData {
	return &initData{
		dbs: dbs,
	}
}

func (m *initData) AdminUser(account, password string) error {
	dbPlatform, err := m.dbs.Get("platform")
	if nil != err {
		return err
	}

	if len(account) == 0 || len(password) < 6 {
		return nil
	}

	genre := uint8(1)

	pwd, err := userutil.NewHasher().Sum(password)
	if nil != err {
		return errors.Wrap(err, "make user password failed")
	}

	// 初始数据
	datas := []interface{}{
		&platform.User{
			ID:      1,
			Account: account,
			//CheckedGenre: uint8(global.RoleSuper),
			CheckedGenre: genre,
			Nickname:     "boss",
			Password:     pwd,
			State:        1,
		},
		&platform.UserRole{
			ID:     1,
			UserID: 1,
			//Genre:  uint8(global.RoleSuper),
			Genre: genre,
		},
	}
	for _, data := range datas {
		if err := dbPlatform.FirstOrCreate(data).Error; nil != err {
			return err
		}
	}

	return nil
}
