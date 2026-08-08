package models

type User struct {
	Name   string `json:"name,omitempty" db:"-"`
	Email  string `json:"email,omitempty" db:"email"`
	UserID string `json:"id" db:"user_id"`

	Username string `json:"-" db:"username"`

	AuthState uint8 `json:"-"`
	Remember  bool  `json:"-"`
}

type StoryVisibility = int8

const (
	VisibilityPrivate StoryVisibility = 0x1
	VisibilityPublic  StoryVisibility = 0x2
)

type Story struct {
	StoryID     string          `json:"story_id" db:"story_id"`
	StoryName   string          `json:"story_name" db:"story_name"`
	Description string          `json:"description" db:"description"`
	Visibility  StoryVisibility `json:"visibility" db:"visibility"`
	Forkable    bool            `json:"forkable" db:"forkable"`
	Username    string          `json:"username" db:"username"`
	LastUpdated int64           `json:"last_updated,omitempty" db:"last_updated"`
	Chapters    int             `json:"chapters" db:"chapters"`
}

type Chapter struct {
	ChapterID   string `json:"chapter_id" db:"chapter_id"`
	ChapterName string `json:"chapter_name" db:"chapter_name"`
	StoryID     string `json:"story_id" db:"story_id"`
	Content     string `json:"content,omitempty" db:"content"`
	LastUpdated int64  `json:"last_updated,omitempty" db:"last_updated"`
	Username    string `json:"username" db:"username"`

	FileID string `json:"-" db:"file_id"`
}

type Edge struct {
	FromChapter string `json:"from_chapter" db:"from_chapter"`
	ToChapter   string `json:"to_chapter" db:"to_chapter"`
}
