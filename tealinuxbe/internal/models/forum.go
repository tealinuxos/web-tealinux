package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Slug        string    `gorm:"size:100;uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Order       int       `gorm:"default:0" json:"order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Topic struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title      string    `gorm:"size:255;not null" json:"title"`
	Slug       string    `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Type       string    `gorm:"size:20;default:'discussion'" json:"type"` // bug, question, discussion
	UserID     uint      `gorm:"not null" json:"user_id"`                  // References User.ID (uint)
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null" json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category"`
	Views      int       `gorm:"default:0" json:"views"`
	IsPinned   bool      `gorm:"default:false" json:"is_pinned"`
	IsLocked   bool      `gorm:"default:false" json:"is_locked"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Posts      []Post    `gorm:"foreignKey:TopicID" json:"posts"`
	Tags       []Tag     `gorm:"many2many:topic_tags;" json:"tags"`
}

type Post struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	TopicID   uuid.UUID  `gorm:"type:uuid;not null" json:"topic_id"`
	UserID    uint       `gorm:"not null" json:"user_id"` // References User.ID (uint)
	User      User       `gorm:"foreignKey:UserID" json:"user"`
	Content   string     `gorm:"type:text;not null" json:"content"`
	ReplyToID *uuid.UUID `gorm:"type:uuid" json:"reply_to_id"` // For threaded replies/quotes
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Likes     []Like     `gorm:"foreignKey:PostID" json:"likes"`
}

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	PostID    uuid.UUID `gorm:"type:uuid;not null" json:"post_id"`
	UserID    uint      `gorm:"not null" json:"user_id"` // References User.ID (uint)
	CreatedAt time.Time `json:"created_at"`
}

type Tag struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name string    `gorm:"size:50;uniqueIndex;not null" json:"name"`
}

type TopicTag struct {
	TopicID uuid.UUID `gorm:"type:uuid;primaryKey" json:"topic_id"`
	TagID   uuid.UUID `gorm:"type:uuid;primaryKey" json:"tag_id"`
}
