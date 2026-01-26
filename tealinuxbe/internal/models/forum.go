package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `gorm:"size:100;not null"`
	Slug        string    `gorm:"size:100;uniqueIndex;not null"`
	Description string    `gorm:"type:text"`
	Order       int       `gorm:"default:0"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Topic struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title      string    `gorm:"size:255;not null"`
	Slug       string    `gorm:"size:255;uniqueIndex;not null"`
	UserID     uint      `gorm:"not null"` // References User.ID (uint)
	User       User      `gorm:"foreignKey:UserID"`
	CategoryID uuid.UUID `gorm:"type:uuid;not null"`
	Category   Category  `gorm:"foreignKey:CategoryID"`
	Views      int       `gorm:"default:0"`
	IsPinned   bool      `gorm:"default:false"`
	IsLocked   bool      `gorm:"default:false"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Posts      []Post `gorm:"foreignKey:TopicID"`
	Tags       []Tag  `gorm:"many2many:topic_tags;"`
}

type Post struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TopicID   uuid.UUID  `gorm:"type:uuid;not null"`
	UserID    uint       `gorm:"not null"` // References User.ID (uint)
	User      User       `gorm:"foreignKey:UserID"`
	Content   string     `gorm:"type:text;not null"`
	ReplyToID *uuid.UUID `gorm:"type:uuid"` // For threaded replies/quotes
	CreatedAt time.Time
	UpdatedAt time.Time
	Likes     []Like `gorm:"foreignKey:PostID"`
}

type Like struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PostID    uuid.UUID `gorm:"type:uuid;not null"`
	UserID    uint      `gorm:"not null"` // References User.ID (uint)
	CreatedAt time.Time
}

type Tag struct {
	ID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name string    `gorm:"size:50;uniqueIndex;not null"`
}

type TopicTag struct {
	TopicID uuid.UUID `gorm:"type:uuid;primaryKey"`
	TagID   uuid.UUID `gorm:"type:uuid;primaryKey"`
}
