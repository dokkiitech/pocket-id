package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/pocket-id/pocket-id/backend/internal/common"
	"github.com/pocket-id/pocket-id/backend/internal/dto"
	"github.com/pocket-id/pocket-id/backend/internal/model"
	datatype "github.com/pocket-id/pocket-id/backend/internal/model/types"
	"github.com/pocket-id/pocket-id/backend/internal/storage"
	testutils "github.com/pocket-id/pocket-id/backend/internal/utils/testing"
)

func TestUserServiceVerifyEmail(t *testing.T) {
	t.Run("verifies the email the token was sent to", func(t *testing.T) {
		service, db := newTestUserService(t)
		user := createTestUserWithEmail(t, db, "alice@example.com")
		createTestEmailVerificationToken(t, db, user.ID, "token-alice", "alice@example.com")

		require.NoError(t, service.VerifyEmail(t.Context(), user.ID, "token-alice"))

		var updated model.User
		require.NoError(t, db.First(&updated, "id = ?", user.ID).Error)
		assert.True(t, updated.EmailVerified)
	})

	t.Run("rejects a token that was sent to a different email", func(t *testing.T) {
		service, db := newTestUserService(t)
		user := createTestUserWithEmail(t, db, "victim@example.com")
		createTestEmailVerificationToken(t, db, user.ID, "token-attacker", "attacker@example.com")

		err := service.VerifyEmail(t.Context(), user.ID, "token-attacker")

		var invalidTokenErr *common.InvalidEmailVerificationTokenError
		require.ErrorAs(t, err, &invalidTokenErr)

		var updated model.User
		require.NoError(t, db.First(&updated, "id = ?", user.ID).Error)
		assert.False(t, updated.EmailVerified)
	})
}

func TestUserServiceUpdateUserInvalidatesEmailVerificationTokens(t *testing.T) {
	service, db := newTestUserService(t)
	user := createTestUserWithEmail(t, db, "attacker@example.com")
	createTestEmailVerificationToken(t, db, user.ID, "token-attacker", "attacker@example.com")

	_, err := service.UpdateUser(t.Context(), user.ID, dto.UserCreateDto{
		Username: user.Username,
		Email:    new("victim@example.com"),
	}, true, false)
	require.NoError(t, err)

	var count int64
	require.NoError(t, db.Model(&model.EmailVerificationToken{}).Where("user_id = ?", user.ID).Count(&count).Error)
	assert.Zero(t, count)

	err = service.VerifyEmail(t.Context(), user.ID, "token-attacker")
	var invalidTokenErr *common.InvalidEmailVerificationTokenError
	require.ErrorAs(t, err, &invalidTokenErr)
}

func newTestUserService(t *testing.T) (*UserService, *gorm.DB) {
	t.Helper()

	db := testutils.NewDatabaseForTest(t)

	fileStorage, err := storage.NewDatabaseStorage(db)
	require.NoError(t, err)

	appConfig := NewTestAppConfigService(&model.AppConfig{
		RequireUserEmail:    model.AppConfigVariable{Value: "false"},
		AllowOwnAccountEdit: model.AppConfigVariable{Value: "true"},
		EmailsVerified:      model.AppConfigVariable{Value: "false"},
	})

	userService := NewUserService(
		db,
		nil,
		nil,
		nil,
		appConfig,
		NewCustomClaimService(db),
		NewAppImagesService(map[string]string{}, fileStorage),
		nil,
		fileStorage,
	)

	return userService, db
}

func createTestUserWithEmail(t *testing.T, db *gorm.DB, email string) model.User {
	t.Helper()

	user := model.User{
		Username: "user-" + email,
		Email:    &email,
	}
	require.NoError(t, db.Create(&user).Error)

	return user
}

func createTestEmailVerificationToken(t *testing.T, db *gorm.DB, userID string, token string, email string) {
	t.Helper()

	require.NoError(t, db.Create(&model.EmailVerificationToken{
		UserID:    userID,
		Token:     token,
		Email:     email,
		ExpiresAt: datatype.DateTime(time.Now().Add(time.Hour)),
	}).Error)
}
