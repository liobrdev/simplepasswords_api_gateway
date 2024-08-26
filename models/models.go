package models

import "time"

type User struct {
	Slug            string    `json:"user_slug" gorm:"primaryKey;not null"`
	PasswordHash    []byte    `json:"-" gorm:"not null"`
	Name            string    `json:"name" gorm:"not null"`
	EmailAddress    string    `json:"email_address,omitempty" gorm:"unique;not null"`
	PhoneNumber     string    `json:"phone_number,omitempty" gorm:"unique;not null"`
	IsActive        bool      `json:"-" gorm:"default:true;not null"`
	CreatedAt       time.Time `json:"-" gorm:"autoCreateTime:nano;not null"`
	UpdatedAt       time.Time `json:"-" gorm:"autoUpdateTime:nano;not null"`
}

type ClientSession struct {
	UserSlug  string    `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserSlug;constraint:OnDelete:CASCADE"`
	ClientIP  string    `gorm:"not null"`
	Digest    []byte    `gorm:"unique;not null"`
	TokenKey  string    `gorm:"primaryKey;size:16;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:false;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

type MFAToken struct {
	UserSlug  string    `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserSlug;constraint:OnDelete:CASCADE"`
	KeyDigest	[]byte    `gorm:"primaryKey;not null"`
	OTPDigest	[]byte		`gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime:false;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

func (ClientSession) TableName() string {
	return "client_sessions"
}

func (MFAToken) TableName() string {
	return "mfa_tokens"
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
