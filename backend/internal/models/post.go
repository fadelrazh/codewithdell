package models

import (
	"time"

	"codewithdell/backend/internal/utils"
	"gorm.io/gorm"
)

// Post represents a blog post
type Post struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UUID        string         `json:"uuid" gorm:"uniqueIndex;not null"`
	Title       string         `json:"title" gorm:"not null"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Excerpt     string         `json:"excerpt"`
	FeaturedImage string       `json:"featured_image"`
	Status      PostStatus     `json:"status" gorm:"default:'draft'"`
	PublishedAt *time.Time     `json:"published_at"`
	AuthorID    uint           `json:"author_id" gorm:"not null"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	LikeCount   int            `json:"like_count" gorm:"default:0"`
	CommentCount int           `json:"comment_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Author    User       `json:"author" gorm:"foreignKey:AuthorID"`
	Categories []Category `json:"categories,omitempty" gorm:"many2many:post_categories;"`
	Tags      []Tag      `json:"tags,omitempty" gorm:"many2many:post_tags;"`
	Comments  []Comment  `json:"comments,omitempty" gorm:"foreignKey:PostID"`
	Likes     []Like     `json:"likes,omitempty" gorm:"foreignKey:PostID"`
	Bookmarks []Bookmark `json:"bookmarks,omitempty" gorm:"foreignKey:PostID"`
}

// PostStatus represents post status
type PostStatus string

const (
	PostStatusDraft     PostStatus = "draft"
	PostStatusPublished PostStatus = "published"
	PostStatusArchived  PostStatus = "archived"
)

// TableName specifies the table name for Post
func (Post) TableName() string {
	return "posts"
}

// BeforeCreate is a GORM hook that runs before creating a post
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	if p.UUID == "" {
		p.UUID = utils.GenerateUUID()
	}
	return nil
}

// IsPublished checks if the post is published
func (p *Post) IsPublished() bool {
	return p.Status == PostStatusPublished && p.PublishedAt != nil
}

// IncrementViewCount increments the view count
func (p *Post) IncrementViewCount() {
	p.ViewCount++
}

// IncrementLikeCount increments the like count
func (p *Post) IncrementLikeCount() {
	p.LikeCount++
}

// DecrementLikeCount decrements the like count
func (p *Post) DecrementLikeCount() {
	if p.LikeCount > 0 {
		p.LikeCount--
	}
}

// IncrementCommentCount increments the comment count
func (p *Post) IncrementCommentCount() {
	p.CommentCount++
}

// DecrementCommentCount decrements the comment count
func (p *Post) DecrementCommentCount() {
	if p.CommentCount > 0 {
		p.CommentCount--
	}
} 