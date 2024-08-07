package helpers

import (
	"testing"

	"gorm.io/gorm"
)

func CountUsers(t *testing.T, db *gorm.DB, userCount *int64) {
	if result := db.Table("users").Count(userCount); result.Error != nil {
		t.Fatalf("User count failed: %s", result.Error.Error())
	}
}

func CountLogs(t *testing.T, db *gorm.DB, logCount *int64) {
	if result := db.Table("logs").Count(logCount); result.Error != nil {
		t.Fatalf("Log count failed: %s", result.Error.Error())
	}
}

func CountClientSessions(t *testing.T, db *gorm.DB, sessionCount *int64) {
	if result := db.Table("client_sessions").Count(sessionCount); result.Error != nil {
		t.Fatalf("Client session count failed: %s", result.Error.Error())
	}
}

func CountMFATokens(t *testing.T, db *gorm.DB, mfaTokenCount *int64) {
	if result := db.Table("mfa_tokens").Count(mfaTokenCount); result.Error != nil {
		t.Fatalf("MFA token count failed: %s", result.Error.Error())
	}
}
