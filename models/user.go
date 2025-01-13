package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"type:varchar(50);unique;not null"`
	Password     string    `gorm:"type:varchar(255);not null"`
	CurrentToken string    `gorm:"type:varchar(512)"`
	LastLogin    time.Time `gorm:"type:timestamp with time zone"`
	CreatedAt    time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP;autoUpdateTime:true"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "users"
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies if the provided password matches the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// ToDict converts the user to a map representation
func (u *User) ToDict() map[string]interface{} {
	var lastLogin *string
	if !u.LastLogin.IsZero() {
		formatted := u.LastLogin.Format(time.RFC3339)
		lastLogin = &formatted
	}

	return map[string]interface{}{
		"user_id":    u.ID,
		"username":   u.Username,
		"created_at": u.CreatedAt.Format(time.RFC3339),
		"last_login": lastLogin,
	}
}
