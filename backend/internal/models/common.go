package models

import (
	"time"

	"codewithdell/backend/internal/utils"
	"gorm.io/gorm"
)

// Category represents a category for posts and projects
type Category struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UUID        string         `json:"uuid" gorm:"uniqueIndex;not null"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Color       string         `json:"color"`
	Icon        string         `json:"icon"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Posts    []Post    `json:"posts,omitempty" gorm:"many2many:post_categories;"`
	Projects []Project `json:"projects,omitempty" gorm:"many2many:project_categories;"`
}

// Tag represents a tag for posts and projects
type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Slug      string         `json:"slug" gorm:"uniqueIndex;not null"`
	Color     string         `json:"color"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Posts    []Post    `json:"posts,omitempty" gorm:"many2many:post_tags;"`
	Projects []Project `json:"projects,omitempty" gorm:"many2many:project_tags;"`
}

// Technology represents a technology used in projects
type Technology struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UUID        string         `json:"uuid" gorm:"uniqueIndex;not null"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null"`
	Description string         `json:"description"`
	Icon        string         `json:"icon"`
	Color       string         `json:"color"`
	Website     string         `json:"website"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Projects []Project `json:"projects,omitempty" gorm:"many2many:project_technologies;"`
}

// Comment represents a comment on posts or projects
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    *uint          `json:"post_id"`
	ProjectID *uint          `json:"project_id"`
	ParentID  *uint          `json:"parent_id"`
	Status    CommentStatus  `json:"status" gorm:"default:'approved'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     User      `json:"user" gorm:"foreignKey:UserID"`
	Post     *Post     `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Project  *Project  `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
	Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []Comment `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// CommentStatus represents comment status
type CommentStatus string

const (
	CommentStatusPending  CommentStatus = "pending"
	CommentStatusApproved CommentStatus = "approved"
	CommentStatusSpam     CommentStatus = "spam"
)

// Like represents a like on posts or projects
type Like struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    *uint          `json:"post_id"`
	ProjectID *uint          `json:"project_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User    User     `json:"user" gorm:"foreignKey:UserID"`
	Post    *Post    `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// Bookmark represents a bookmark on posts or projects
type Bookmark struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;not null"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	PostID    *uint          `json:"post_id"`
	ProjectID *uint          `json:"project_id"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User    User     `json:"user" gorm:"foreignKey:UserID"`
	Post    *Post    `json:"post,omitempty" gorm:"foreignKey:PostID"`
	Project *Project `json:"project,omitempty" gorm:"foreignKey:ProjectID"`
}

// Screenshot represents a screenshot for projects
type Screenshot struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UUID      string         `json:"uuid" gorm:"uniqueIndex;not null"`
	ProjectID uint           `json:"project_id" gorm:"not null"`
	Title     string         `json:"title"`
	ImageURL  string         `json:"image_url" gorm:"not null"`
	AltText   string         `json:"alt_text"`
	Order     int            `json:"order" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Project Project `json:"project" gorm:"foreignKey:ProjectID"`
}

// TableName functions
func (Category) TableName() string {
	return "categories"
}

func (Tag) TableName() string {
	return "tags"
}

func (Technology) TableName() string {
	return "technologies"
}

func (Comment) TableName() string {
	return "comments"
}

func (Like) TableName() string {
	return "likes"
}

func (Bookmark) TableName() string {
	return "bookmarks"
}

func (Screenshot) TableName() string {
	return "screenshots"
}

// BeforeCreate hooks
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.UUID == "" {
		c.UUID = utils.GenerateUUID()
	}
	return nil
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.UUID == "" {
		t.UUID = utils.GenerateUUID()
	}
	return nil
}

func (t *Technology) BeforeCreate(tx *gorm.DB) error {
	if t.UUID == "" {
		t.UUID = utils.GenerateUUID()
	}
	return nil
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.UUID == "" {
		c.UUID = utils.GenerateUUID()
	}
	return nil
}

func (l *Like) BeforeCreate(tx *gorm.DB) error {
	if l.UUID == "" {
		l.UUID = utils.GenerateUUID()
	}
	return nil
}

func (b *Bookmark) BeforeCreate(tx *gorm.DB) error {
	if b.UUID == "" {
		b.UUID = utils.GenerateUUID()
	}
	return nil
}

func (s *Screenshot) BeforeCreate(tx *gorm.DB) error {
	if s.UUID == "" {
		s.UUID = utils.GenerateUUID()
	}
	return nil
} 