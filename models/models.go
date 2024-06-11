package models

import "time"

type User struct {
	Slug            string    `json:"user_slug" gorm:"primaryKey;not null"`
	PasswordSalt    string    `json:"-" gorm:"not null"`
	PasswordHash    []byte    `json:"-" gorm:"not null"`
	Name            string    `json:"name" gorm:"not null"`
	EmailAddress    string    `json:"email_address" gorm:"uniqueIndex;not null"`
	EmailIsVerified bool      `json:"email_is_verified" gorm:"default:false;not null"`
	PhoneNumber     string    `json:"phone_number" gorm:"not null"`
	PhoneIsVerified bool      `json:"phone_is_verified" gorm:"default:false;not null"`
	MfaIsEnabled    bool      `json:"mfa_is_enabled" gorm:"default:false;not null"`
	IsActive        bool      `json:"is_active" gorm:"default:true;not null"`
	CreatedAt       time.Time `json:"-" gorm:"autoCreateTime:nano;not null"`
	UpdatedAt       time.Time `json:"-" gorm:"autoUpdateTime:nano;not null"`
}

type DiscreetUser struct {
	Slug 									string `json:"user_slug"`
	Name            			string `json:"name"`
	TruncatedEmailAddress string `json:"truncated_email_address"`
	EmailIsVerified 			bool   `json:"email_is_verified"`
	TruncatedPhoneNumber  string `json:"truncated_phone_number"`
	PhoneIsVerified 			bool   `json:"phone_is_verified"`
	MfaIsEnabled    			bool   `json:"mfa_is_enabled"`
}

type ClientSession struct {
	UserSlug  string    `gorm:"index;not null"`
	User      User      `gorm:"foreignKey:UserSlug;constraint:OnDelete:CASCADE"`
	ClientIP  string    `gorm:"not null"`
	Digest    []byte    `gorm:"unique;not null"`
	TokenKey  string    `gorm:"primaryKey;size:16;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:false;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

type EmailVerificationToken struct {
	UserSlug  string    `gorm:"index;not null"`
	User      User      `gorm:"foreignKey:UserSlug;constraint:OnDelete:CASCADE"`
	Digest    []byte    `gorm:"unique;not null"`
	TokenKey  string    `gorm:"primaryKey;size:16;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:false;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (ClientSession) TableName() string {
	return "client_sessions"
}

type Log struct {
	ID              uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Caller          string    `json:"caller"`
	ClientIP        string    `json:"client_ip"`
	ClientOperation string    `json:"client_operation"`
	ContextString   string    `json:"context_string"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime:nano;not null"`
	Detail          string    `json:"detail"`
	Extra           string    `json:"extra"`
	Level           string    `json:"level"`
	Message         string    `json:"message"`
	RequestBody     string    `json:"request_body"`
	UserSlug        string    `json:"user_slug"`
}
