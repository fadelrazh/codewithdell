package models

import (
	"time"

	"codewithdell/backend/internal/utils"
	"gorm.io/gorm"
)

// Project represents a showcase project
type Project struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UUID        string         `json:"uuid" gorm:"uniqueIndex;not null"`
	Title       string         `json:"title" gorm:"not null"`
	Slug        string         `json:"slug" gorm:"uniqueIndex;not null"`
	Description string         `json:"description" gorm:"type:text;not null"`
	Content     string         `json:"content" gorm:"type:text"`
	FeaturedImage string       `json:"featured_image"`
	Status      ProjectStatus  `json:"status" gorm:"default:'draft'"`
	PublishedAt *time.Time     `json:"published_at"`
	AuthorID    uint           `json:"author_id" gorm:"not null"`
	ViewCount   int            `json:"view_count" gorm:"default:0"`
	LikeCount   int            `json:"like_count" gorm:"default:0"`
	CommentCount int           `json:"comment_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Project specific fields
	LiveURL    string `json:"live_url"`
	SourceURL  string `json:"source_url"`
	DemoURL    string `json:"demo_url"`
	Difficulty ProjectDifficulty `json:"difficulty" gorm:"default:'intermediate'"`
	Duration   string `json:"duration"`
	TeamSize   int    `json:"team_size" gorm:"default:1"`

	// Relationships
	Author     User        `json:"author" gorm:"foreignKey:AuthorID"`
	Categories []Category  `json:"categories,omitempty" gorm:"many2many:project_categories;"`
	Tags       []Tag       `json:"tags,omitempty" gorm:"many2many:project_tags;"`
	Technologies []Technology `json:"technologies,omitempty" gorm:"many2many:project_technologies;"`
	Comments   []Comment   `json:"comments,omitempty" gorm:"foreignKey:ProjectID"`
	Likes      []Like      `json:"likes,omitempty" gorm:"foreignKey:ProjectID"`
	Bookmarks  []Bookmark  `json:"bookmarks,omitempty" gorm:"foreignKey:ProjectID"`
	Screenshots []Screenshot `json:"screenshots,omitempty" gorm:"foreignKey:ProjectID"`
}

// ProjectStatus represents project status
type ProjectStatus string

const (
	ProjectStatusDraft     ProjectStatus = "draft"
	ProjectStatusPublished ProjectStatus = "published"
	ProjectStatusArchived  ProjectStatus = "archived"
)

// ProjectDifficulty represents project difficulty level
type ProjectDifficulty string

const (
	DifficultyBeginner     ProjectDifficulty = "beginner"
	DifficultyIntermediate ProjectDifficulty = "intermediate"
	DifficultyAdvanced     ProjectDifficulty = "advanced"
	DifficultyExpert       ProjectDifficulty = "expert"
)

// TableName specifies the table name for Project
func (Project) TableName() string {
	return "projects"
}

// BeforeCreate is a GORM hook that runs before creating a project
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.UUID == "" {
		p.UUID = utils.GenerateUUID()
	}
	return nil
}

// IsPublished checks if the project is published
func (p *Project) IsPublished() bool {
	return p.Status == ProjectStatusPublished && p.PublishedAt != nil
}

// IncrementViewCount increments the view count
func (p *Project) IncrementViewCount() {
	p.ViewCount++
}

// IncrementLikeCount increments the like count
func (p *Project) IncrementLikeCount() {
	p.LikeCount++
}

// DecrementLikeCount decrements the like count
func (p *Project) DecrementLikeCount() {
	if p.LikeCount > 0 {
		p.LikeCount--
	}
}

// IncrementCommentCount increments the comment count
func (p *Project) IncrementCommentCount() {
	p.CommentCount++
}

// DecrementCommentCount decrements the comment count
func (p *Project) DecrementCommentCount() {
	if p.CommentCount > 0 {
		p.CommentCount--
	}
}

// HasLiveDemo checks if the project has a live demo
func (p *Project) HasLiveDemo() bool {
	return p.LiveURL != "" || p.DemoURL != ""
}

// HasSourceCode checks if the project has source code available
func (p *Project) HasSourceCode() bool {
	return p.SourceURL != ""
} 