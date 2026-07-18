package server

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5"
)

//TODO: Add page functionality in getStories
//TODO: Add "last modified" statistic to stories
//TODO: Add visibility (pub/priv/all) field to stories

func Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf(
		"postgres://%v:%v@%v:%v/%v",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DB"),
	))

	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to DB")
	return conn, nil
}

func validateUserID(id string) bool {
	for _, c := range id {
		if (c < '0') || (c > '9') {
			return false
		}
	}
	return true
}

func validateUsername(name string) error {
	name = strings.ToLower(name)

	_, err := regexp.Match("^[a-z0-9_]{3,18}", []byte(name))

	return err
}

func (s *Server) registerUser(user *User) error {
	_, err := s.PgConn.Exec(context.Background(),
		`INSERT INTO users (username, email, user_id) values ($1, $2, $3)`,
		user.Username, user.Email, user.GoogleUserID)

	if err != nil {
		return err
	}

	user.AuthState |= AuthDB
	return nil
}

func (s *Server) getUser(u *User) error {
	ok := validateUserID(u.GoogleUserID)

	if !ok {
		fmt.Println("Invalid user ID")
		return fmt.Errorf("invalid user id")
	}

	err := s.PgConn.QueryRow(context.Background(), `SELECT username FROM users WHERE user_id=$1`, u.GoogleUserID).Scan(&u.Username)
	if (err != nil) && (err == pgx.ErrNoRows) {
		return err
	} else if err != nil {
		fmt.Println("Error getting user:", u.GoogleUserID, err)
		return err
	}

	u.AuthState |= AuthDB
	return nil
}

func (s *Server) insertStory(story *Story, u *User) error {
	_, err := s.PgConn.Exec(context.Background(),
		`INSERT INTO stories (story_id, username, story_name, description, visibility) VALUES ($1, $2, $3, $4, $5)`,
		story.StoryID, u.Username, story.StoryName, story.Description, story.Visibility)

	return err
}

func (s *Server) updateStory(story *Story, u *User) error {
	_, err := s.PgConn.Exec(context.Background(),
		`UPDATE stories SET story_name=$1, description=$2, visibility=$3 WHERE story_id=$4`,
		story.StoryName, story.Description, story.Visibility, story.StoryID)

	return err
}

func (s *Server) getStory(storyId string) (*Story, error) {
	story := Story{
		StoryID: storyId,
	}

	err := s.PgConn.QueryRow(context.Background(),
		`SELECT story_name, description, visibility, username FROM stories WHERE story_id=$1`,
		storyId).Scan(&story.StoryName, &story.Description, &story.Visibility, &story.Username)

	if err != nil {
		return nil, err
	}

	return &story, err
}

func putStoryRowsToSlice(rows pgx.Rows) ([]Story, error) {
	stories := make([]Story, 0, 50)
	for rows.Next() {
		var story Story
		err := rows.Scan(&story.StoryName, &story.Description, &story.Visibility, &story.StoryID)
		if err != nil {
			return nil, err
		}
		stories = append(stories, story)
	}

	return stories, nil
}

func (s *Server) getAllUserStoriesFromDB(username string) ([]Story, error) {
	rows, err := s.PgConn.Query(context.Background(),
		`SELECT story_name, description, visibility, story_id FROM stories WHERE username=$1 LIMIT 50 OFFSET 0`,
		username)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stories, err := putStoryRowsToSlice(rows)
	if err != nil {
		return nil, err
	}

	for _, story := range stories {
		story.Username = username
	}

	return stories, nil
}

func (s *Server) getAllUserStoriesFromDBWithVisibility(username string, visibility int8) ([]Story, error) {
	rows, err := s.PgConn.Query(context.Background(),
		`SELECT story_name, description, visibility, story_id FROM stories WHERE username=$1 AND visibility=$2 LIMIT 50 OFFSET 0`,
		username, visibility)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	stories, err := putStoryRowsToSlice(rows)
	if err != nil {
		return nil, err
	}

	for _, story := range stories {
		story.Username = username
	}

	return stories, nil
}

func (s *Server) insertChapter(chapter *Chapter, user *User) error {
	_, err := s.PgConn.Exec(context.Background(),
		`INSERT INTO chapters (chapter_id, chapter_name, story_id, file_id) VALUES ($1, $2, $3, $4)`,
		chapter.ChapterID, chapter.ChapterName, chapter.StoryID, chapter.FileID)

	return err
}

func saveChapterContentToFile(chapter *Chapter, content string) error {
	filePath := filepath.Join("./chapters", chapter.FileID)
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	if err != nil {
		return err
	}

	defer file.Close()

	n, err := file.WriteString(content)

	if err != nil {
		return err
	}

	if n != len(content) {
		return errors.New("failed to write to file")
	}

	return nil
}

func (s *Server) getChapter(chapterId string) (*Chapter, error) {
	chapter := Chapter{
		ChapterID: chapterId,
	}

	err := s.PgConn.QueryRow(context.Background(),
		`SELECT chapter_name, story_id, file_id FROM chapters WHERE chapter_id=$1`,
		chapterId).Scan(&chapter.ChapterName, &chapter.StoryID, &chapter.FileID)

	if err != nil {
		return nil, err
	}

	return &chapter, err
}

func (s *Server) updateChapter(chapter *Chapter, user *User) error {
	_, err := s.PgConn.Exec(context.Background(),
		`UPDATE chapters SET chapter_name=$1 WHERE chapter_id=$2`,
		chapter.ChapterName, chapter.ChapterID)

	return err
}

func (s *Server) getChaptersByStory(storyID string) ([]Chapter, error) {
	rows, err := s.PgConn.Query(context.Background(),
		`SELECT chapter_id, chapter_name, file_id  FROM chapters WHERE story_id=$1 LIMIT 50 OFFSET 0`,
		storyID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	chapters := make([]Chapter, 0, 50)
	for rows.Next() {
		chapter := Chapter{
			StoryID: storyID,
		}

		if err := rows.Scan(&chapter.ChapterID, &chapter.ChapterName, &chapter.FileID); err != nil {
			return nil, err
		}

		chapters = append(chapters, chapter)
	}

	return chapters, nil
}

func getContent(fileID string) (string, error) {
	file, err := os.ReadFile(filepath.Join("./chapters", fileID))
	if err != nil {
		return "", err
	}

	return string(file), nil
}
