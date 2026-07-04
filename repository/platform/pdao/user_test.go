package pdao

import (
	"context"
	"log"
	"os"
	"testing"

	"server-api/global"
	"server-api/repository/platform/pattr"
	"server-api/repository/platform/pmodel"
	"server-api/testhelper"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain 统一初始化集成测试环境。无数据库时跳过整个包，而非视为失败。
func TestMain(m *testing.M) {
	if err := testhelper.SetupTestWithDB(); err != nil {
		log.Printf("skip pdao tests: db unavailable: %v", err)
		os.Exit(0)
	}
	code := m.Run()
	testhelper.Cleanup()
	os.Exit(code)
}

// TestCreateAndBindThirdUser_ReviveSoftDeletedUser 验证核心修复：
// 当 account 已存在但被软删除时，CreateAndBindThirdUser 应复活原记录（更新字段并清空 deleted_at），
// 而非走 Create 分支导致 account 唯一约束冲突或产生重复记录。
func TestCreateAndBindThirdUser_ReviveSoftDeletedUser(t *testing.T) {
	ctx := context.Background()
	account := "revive_test_user"

	db, _ := global.Database().Get("platform")

	// 清理历史数据（含软删除），保证测试隔离
	cleanupUser := func() {
		var ids []uint
		db.Unscoped().Model(&pmodel.User{}).Where("account = ?", account).Pluck("id", &ids)
		if len(ids) > 0 {
			db.Unscoped().Where("user_id IN ?", ids).Delete(&pmodel.UserStat{})
			db.Unscoped().Where("user_id IN ?", ids).Delete(&pmodel.UserRole{})
			db.Unscoped().Where("id IN ?", ids).Delete(&pmodel.User{})
		}
	}
	cleanupUser()
	defer cleanupUser()

	// 1. 先创建一个用户
	original := &pmodel.User{
		Account:  account,
		Nickname: "old_nick",
		Password: "old_pwd",
	}
	originalStat := &pmodel.UserStat{
		FromPlatformID: pattr.UserFromPlatformAccount,
	}
	genres := []int8{int8(global.RoleUser)}
	require.NoError(t, NewUser().Create(ctx, genres, original, originalStat))
	originalID := original.ID
	require.NotZero(t, originalID)

	// 2. 软删除该用户
	require.NoError(t, db.Delete(&pmodel.User{}, originalID).Error)

	// 3. 用同 account 再次调用（复活场景），Account 平台跳过第三方绑定
	revive := &pmodel.User{
		Account:  account,
		Nickname: "new_nick",
		Password: "new_pwd",
	}
	reviveStat := &pmodel.UserStat{
		FromPlatformID: pattr.UserFromPlatformAccount,
	}
	err := NewUser().CreateAndBindThirdUser(ctx, genres, revive, reviveStat, 0)
	assert.NoError(t, err)

	// 4. 应复活原记录，而非新建
	assert.Equal(t, originalID, revive.ID, "should revive original record, not create new")

	// 5. DB 中 deleted_at 应为 NULL（已复活），nickname/password 已更新
	var got pmodel.User
	require.NoError(t, db.Unscoped().First(&got, originalID).Error)
	assert.False(t, got.DeletedAt.Valid, "deleted_at should be NULL after revive")
	assert.Equal(t, "new_nick", got.Nickname)
	assert.Equal(t, "new_pwd", got.Password)
}

// TestCreateAndBindThirdUser_CreateNewWhenAbsent 验证 account 不存在时正常新建。
func TestCreateAndBindThirdUser_CreateNewWhenAbsent(t *testing.T) {
	ctx := context.Background()
	account := "create_new_test_user"

	db, _ := global.Database().Get("platform")

	cleanupUser := func() {
		var ids []uint
		db.Unscoped().Model(&pmodel.User{}).Where("account = ?", account).Pluck("id", &ids)
		if len(ids) > 0 {
			db.Unscoped().Where("user_id IN ?", ids).Delete(&pmodel.UserStat{})
			db.Unscoped().Where("user_id IN ?", ids).Delete(&pmodel.UserRole{})
			db.Unscoped().Where("id IN ?", ids).Delete(&pmodel.User{})
		}
	}
	cleanupUser()
	defer cleanupUser()

	record := &pmodel.User{
		Account:  account,
		Nickname: "fresh_nick",
		Password: "fresh_pwd",
	}
	stat := &pmodel.UserStat{
		FromPlatformID: pattr.UserFromPlatformAccount,
	}
	genres := []int8{int8(global.RoleUser)}

	err := NewUser().CreateAndBindThirdUser(ctx, genres, record, stat, 0)
	assert.NoError(t, err)
	assert.NotZero(t, record.ID, "should create new record")

	// DB 中确实存在且未被软删除
	var got pmodel.User
	require.NoError(t, db.Unscoped().First(&got, record.ID).Error)
	assert.Equal(t, account, got.Account)
	assert.False(t, got.DeletedAt.Valid)
}
