package helpers

import (
	"testing"

	"gorm.io/gorm"

	"github.com/liobrdev/simplepasswords_api_gateway/models"
)

func QueryTestUserBySlug(t *testing.T, db *gorm.DB, user *models.User, slug string) {
	if result := db.First(&user, "slug = ?", slug); result.Error != nil {
		t.Fatalf("User query by slug failed: %s", result.Error.Error())
	}
}

func QueryTestUserByEmail(t *testing.T, db *gorm.DB, user *models.User, email string) {
	if result := db.First(&user, "email_address = ?", email); result.Error != nil {
		t.Fatalf("User query by email failed: %s", result.Error.Error())
	}
}

func QueryTestLogLatest(t *testing.T, db *gorm.DB, log *models.Log) {
	if result := db.Last(&log); result.Error != nil {
		t.Fatalf("Latest log query failed: %s", result.Error.Error())
	}
}

func QueryTestClientSessionLatest(t *testing.T, db *gorm.DB, session *models.ClientSession) {
	if result := db.Order("created_at DESC").Limit(1).Find(&session); result.Error != nil {
		t.Fatalf("Latest client session query failed: %s", result.Error.Error())
	}
}

func QueryTestMFATokenLatest(t *testing.T, db *gorm.DB, mfaToken *models.MFAToken) {
	if result := db.Order("created_at DESC").Limit(1).Find(&mfaToken); result.Error != nil {
		t.Fatalf("Latest mfa token query failed: %s", result.Error.Error())
	}
}

func QueryTestEmailTokenLatest(
	t *testing.T, db *gorm.DB, emailToken *models.EmailVerificationToken,
) {
	if result := db.Order("created_at DESC").Limit(1).Find(&emailToken); result.Error != nil {
		t.Fatalf("Latest email token query failed: %s", result.Error.Error())
	}
}
