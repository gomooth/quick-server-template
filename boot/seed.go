package boot

import (
	"server-api/global"
	"server-api/repository/platform/pmodel"

	"github.com/gomooth/utils/userutil"
	"github.com/gomooth/xerror"
)

// seedAdminUser 初始化超级管理员账号。
// 仅在账号非空且密码不少于 6 位时执行；通过 FirstOrCreate 保证幂等，重复启动不重复创建、不覆盖已有密码。
// 角色使用 global.RoleSuper，避免在 global/internal 层因循环依赖无法引用而硬编码魔数。
func seedAdminUser(account, password string) error {
	if len(account) == 0 || len(password) < 6 {
		return nil
	}

	db, err := global.Database().Get("platform")
	if err != nil {
		return err
	}

	pwd, err := userutil.Sum(password)
	if err != nil {
		return xerror.Wrap(err, "make admin password failed")
	}

	genre := uint8(global.RoleSuper)

	// 管理员用户
	admin := &pmodel.User{
		Account:      account,
		CheckedGenre: genre,
		Nickname:     "boss",
		Password:     pwd,
		State:        1,
	}
	admin.ID = 1
	if err := db.FirstOrCreate(admin).Error; err != nil {
		return xerror.Wrap(err, "seed admin user failed")
	}

	// 管理员角色
	role := &pmodel.UserRole{
		UserID: 1,
		Genre:  genre,
	}
	role.ID = 1
	if err := db.FirstOrCreate(role).Error; err != nil {
		return xerror.Wrap(err, "seed admin role failed")
	}

	return nil
}
