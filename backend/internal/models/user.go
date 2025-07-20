package models

import (
	"time"

	"codewithdell/backend/internal/utils"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"`
	FirstName string         `json:"first_name" gorm:"not null"`
	LastName  string         `json:"last_name" gorm:"not null"`
	Avatar    string         `json:"avatar"`
	Bio       string         `json:"bio"`
	Role      UserRole       `json:"role" gorm:"default:'user'"`
	Status    UserStatus     `json:"status" gorm:"default:'active'"`
	Verified  bool           `json:"verified" gorm:"default:false"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Posts     []Post     `json:"posts,omitempty" gorm:"foreignKey:AuthorID"`
	Comments  []Comment  `json:"comments,omitempty" gorm:"foreignKey:UserID"`
	Projects  []Project  `json:"projects,omitempty" gorm:"foreignKey:AuthorID"`
	Likes     []Like     `json:"likes,omitempty" gorm:"foreignKey:UserID"`
	Bookmarks []Bookmark `json:"bookmarks,omitempty" gorm:"foreignKey:UserID"`
}

// UserRole represents user roles
type UserRole string

const (
	RoleAdmin  UserRole = "admin"
	RoleEditor UserRole = "editor"
	RoleUser   UserRole = "user"
)

// UserStatus represents user status
type UserStatus string

const (
	StatusActive   UserStatus = "active"
	StatusInactive UserStatus = "inactive"
	StatusBanned   UserStatus = "banned"
)

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UUID == "" {
		u.UUID = utils.GenerateUUID()
	}
	return nil
}

// GetFullName returns the full name of the user
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsEditor checks if the user is an editor
func (u *User) IsEditor() bool {
	return u.Role == RoleEditor || u.Role == RoleAdmin
}

// IsActive checks if the user is active
func (u *User) IsActive() bool {
	return u.Status == StatusActive
} 