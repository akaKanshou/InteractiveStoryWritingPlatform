package models

const (
	visibilityPrivate int8 = 0x1
	visibilityPublic  int8 = 0x2
)

type User struct {
	Name   string `json:"name" db:"-"`
	Email  string `json:"email,omitempty" db:"email"`
	UserID string `json:"id" db:"user_id"`

	Username string `json:"-" db:"username"`

	AuthState uint8 `json:"-"`
	Remember  bool  `json:"-"`
}

type Story struct {
	StoryID     string `json:"story_id" db:"story_id"`
	StoryName   string `json:"story_name" db:"story_name"`
	Description string `json:"description" db:"description"`
	Visibility  int    `json:"visibility" db:"visibility"`
	Username    string `json:"username" db:"username"`
}

type Chapter struct {
	ChapterID   string `json:"chapter_id" db:"chapter_id"`
	ChapterName string `json:"chapter_name" db:"chapter_name"`
	StoryID     string `json:"story_id" db:"story_id"`
	Content     string `json:"content,omitempty" db:"content"`

	FileID string `json:"-" db:"file_id"`
}

type Edge struct {
	FromChapter string `json:"from_chapter" db:"from_chapter"`
	ToChapter   string `json:"to_chapter" db:"to_chapter"`
}
