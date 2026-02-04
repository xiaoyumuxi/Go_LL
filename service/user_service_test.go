package service

import (
	"gin-crud/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// mockDB 创建一个 Mock 的 GORM DB 实例
func mockDB() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})

	return gormDB, mock, err
}

func TestUserService_GetUser(t *testing.T) {
	db, mock, err := mockDB()
	if err != nil {
		t.Fatalf("failed to mock db: %v", err)
	}

	userService := &UserService{DB: db}

	t.Run("UserExists", func(t *testing.T) {
		userID := "123"
		expectedUser := models.User{
			Username: "testuser",
		}

		rows := sqlmock.NewRows([]string{"id", "username"}).
			AddRow(userID, expectedUser.Username)

		// 修正点 1: 匹配 GORM 默认生成的 SQL
		// GORM 的 First() 会生成: SELECT * FROM `users` WHERE `users`.`id` = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT ?
		// 我们使用更宽泛的正则来匹配，或者精确匹配
		// 注意：LIMIT ? 意味着有两个参数
		mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
			WithArgs(userID, 1). // 修正点 2: 增加 LIMIT 的参数 1
			WillReturnRows(rows)

		user, err := userService.GetUser(userID)

		assert.NoError(t, err)
		// 修正点 3: 先判断 user 是否为 nil，防止 panic
		if assert.NotNil(t, user) {
			assert.Equal(t, expectedUser.Username, user.Username)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userID := "999"

		// 同样需要匹配完整的 SQL 和参数
		mock.ExpectQuery("^SELECT \\* FROM `users` WHERE `users`.`id` = \\? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT \\?$").
			WithArgs(userID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := userService.GetUser(userID)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, "用户不存在", err.Error())
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
