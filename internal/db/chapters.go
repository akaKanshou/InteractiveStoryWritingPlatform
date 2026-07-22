package db

import (
	"context"
	"errors"
	"fmt"
	fverrors "forgeverse/internal/errors"
	"forgeverse/internal/models"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// chapterDir is the prefix for chapter storage
const chapterDir = "./chapters"

func GetChapterByID(chapterID string) (*models.Chapter, fverrors.Error) {
	chapter := &models.Chapter{
		ChapterID: chapterID,
	}

	err := dbConn.QueryRow(context.Background(),
		"SELECT chapter_name, story_id, file_id FROM chapters WHERE chapter_id=$1",
		chapterID).Scan(
		&chapter.ChapterName,
		&chapter.StoryID,
		&chapter.FileID,
	)

	if err == nil {
		return chapter, nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fverrors.ChapterNotFoundError
	}

	return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
}

func WriteToFile(fileID string, content string) fverrors.Error {
	file, err := os.OpenFile(filepath.Join(chapterDir, fileID), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fverrors.NewServerError(err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

	n, err := file.WriteString(content)
	if err != nil {
		return fverrors.NewServerError(err)
	}

	if n != len(content) {
		return fverrors.NewServerError(errors.New("could not write all bytes to file"))
	}

	return nil
}

func ReadFromFile(fileID string) (string, fverrors.Error) {
	file, err := os.ReadFile(filepath.Join(chapterDir, fileID))
	if err != nil {
		return "", fverrors.NewServerError(err)
	}

	return string(file), nil
}

func InsertNewChapter(chapter *models.Chapter) fverrors.Error {
	_, err := dbConn.Exec(context.Background(),
		"INSERT INTO chapters (chapter_id, chapter_name, story_id, file_id) VALUES ($1, $2, $3, $4)",
		chapter.ChapterID, chapter.ChapterName, chapter.StoryID, chapter.FileID,
	)

	if err == nil {
		return nil
	}

	if pgErr, okay := errors.AsType[*pgconn.PgError](err); okay && (pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) || pgerrcode.IsDataException(pgErr.Code)) {
		return fverrors.GenericInvalidRequestErr
	}

	return fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
}

func UpdateChapter(chapter *models.Chapter) fverrors.Error {
	_, err := dbConn.Exec(context.Background(),
		"UPDATE chapters SET chapter_name=$1 WHERE chapter_id=$2",
		chapter.ChapterName, chapter.ChapterID)

	if pgErr, okay := errors.AsType[*pgconn.PgError](err); okay && (pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) || pgerrcode.IsDataException(pgErr.Code)) {
		return fverrors.GenericInvalidRequestErr
	} else if err != nil {
		return fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred.", err)
	}

	return nil
}

func GetChaptersByStory(storyID string, page int) ([]*models.Chapter, fverrors.Error) {
	rows, err := dbConn.Query(context.Background(),
		"SELECT chapter_name, chapter_id FROM chapters WHERE story_id=$1 LIMIT 50 OFFSET $2", storyID, page)

	if err != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
	}

	defer rows.Close()

	chapters := make([]*models.Chapter, 0, 50)
	for rows.Next() {
		chapter := &models.Chapter{
			StoryID: storyID,
		}

		err := rows.Scan(&chapter.ChapterName, &chapter.ChapterID)
		if err != nil {
			return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
		}

		chapters = append(chapters, chapter)
	}

	if err := rows.Err(); err != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
	}

	return chapters, nil
}

func InsertNewEdge(from, to string) fverrors.Error {
	_, err := dbConn.Exec(context.Background(),
		"INSERT INTO edges (from_chap, to_chap) VALUES ($1, $2)",
		from, to)

	if err == nil {
		return nil
	}

	if pgErr, okay := errors.AsType[*pgconn.PgError](err); okay && (pgErr.Code == pgerrcode.UniqueViolation) {
		return fverrors.NewDBError(http.StatusInternalServerError, "Edge already exists", err)
	} else if okay && (pgerrcode.IsIntegrityConstraintViolation(pgErr.Code)) {
		return fverrors.GenericInvalidRequestErr
	}

	return fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
}

func GetEdgesByChapter(chapterID string) ([]*models.Edge, fverrors.Error) {
	rows, err := dbConn.Query(context.Background(),
		"SELECT from_chap, to_chap FROM edges WHERE from_chap=$1 OR to_chap=$1",
		chapterID)

	if err != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
	}

	defer rows.Close()

	edges := make([]*models.Edge, 0, 50)
	for rows.Next() {
		edge := &models.Edge{}

		if err := rows.Scan(&edge.FromChapter, &edge.ToChapter); err != nil {
			return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
		}

		edges = append(edges, edge)
	}

	if err := rows.Err(); err != nil {
		return nil, fverrors.NewDBError(http.StatusInternalServerError, "An unexpected error occurred", err)
	}

	return edges, nil
}
