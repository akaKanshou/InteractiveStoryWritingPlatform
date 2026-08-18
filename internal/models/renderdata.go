package models

import (
	"strconv"
	"time"
)

type MyStoriesData struct {
	User                                *User
	Stories                             []Story
	StatScenes, StatPublic, StatPrivate int
	PFP                                 int

	Err error
}

func NewMyStoriesData(user *User, stories []Story) *MyStoriesData {
	data := &MyStoriesData{
		User:        user,
		Stories:     stories,
		StatScenes:  0,
		StatPublic:  0,
		StatPrivate: 0,
		PFP:         0,
	}

	for _, story := range stories {
		if story.Visibility == VisibilityPrivate {
			data.StatPrivate++
		} else {
			data.StatPublic++
		}

		data.StatScenes += story.Chapters
	}

	data.PFP, _ = strconv.Atoi(user.UserID[len(user.UserID)-2:])

	return data
}

type StoryInfoDat struct {
	StoryInfo       Story
	LastUpdatedDate string
	Chapters        []Chapter
}

func NewStoryInfoDat(story *Story) *StoryInfoDat {
	return &StoryInfoDat{
		StoryInfo:       *story,
		LastUpdatedDate: time.Unix(story.LastUpdated, 0).Format("2006-01-02"),
	}
}

type ChapterInfoDat struct {
	ChapterInfo     Chapter
	LastUpdatedDate string
}

func NewChapterInfoDat(chapter *Chapter) *ChapterInfoDat {
	return &ChapterInfoDat{
		ChapterInfo:     *chapter,
		LastUpdatedDate: time.Unix(chapter.LastUpdated, 0).Format("2006-01-02"),
	}
}

type HomePageData struct {
	User *User
	PFP  int

	TrendingStories, EditorsPicks, RecentlyUpdated []Story

	Err error
}

func NewHomePageData(user *User, editorsPicks, trendingStories, recentlyUpdated []Story) *HomePageData {
	data := &HomePageData{
		User: user,

		TrendingStories: trendingStories,
		EditorsPicks:    editorsPicks,
		RecentlyUpdated: recentlyUpdated,
	}

	if user.UserID != "" {
		data.PFP, _ = strconv.Atoi(user.UserID[len(user.UserID)-2:])
	}

	return data
}
