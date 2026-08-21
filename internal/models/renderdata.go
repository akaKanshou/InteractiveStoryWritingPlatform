package models

import (
	"fmt"
	"time"
)

type MyStoriesData struct {
	User    *User
	Stories []StoryInfoDat

	Err error
}

func NewMyStoriesData(user *User, stories []Story) *MyStoriesData {
	storiesInfo := make([]StoryInfoDat, len(stories))
	for i, story := range stories {
		storiesInfo[i] = StoryInfoDat{
			StoryInfo:       story,
			LastUpdatedDate: getTimeString(story.LastUpdated),
		}
	}
	data := &MyStoriesData{
		User:    user,
		Stories: storiesInfo,
	}

	return data
}

type StoryInfoDat struct {
	StoryInfo       Story
	LastUpdatedDate string
	Chapters        []Chapter
}

func getTimeString(lastTime int64) string {
	goTime := time.Unix(lastTime, 0)
	timeSince := time.Since(goTime)

	if timeSince < time.Minute {
		return "Updated less than a minute ago"
	}

	if timeSince < time.Hour {
		return fmt.Sprintf("Updated %d minutes ago", timeSince/time.Minute)
	}

	if timeSince < time.Hour*24 {
		return fmt.Sprintf("Updated %d hours ago", timeSince/time.Hour)
	}

	if timeSince < time.Hour*24*2 {
		return "Updated yesterday"
	}

	year, month, date := goTime.Date()
	return fmt.Sprintf("Updated %d %s %d", date, month.String(), year)
}

func NewStoryInfoDat(story *Story) *StoryInfoDat {
	return &StoryInfoDat{
		StoryInfo:       *story,
		LastUpdatedDate: getTimeString(story.LastUpdated),
	}
}

type ChapterInfoDat struct {
	User            *User
	ChapterInfo     Chapter
	EdgeDetailsInc  []ChapterList
	EdgeDetailsOut  []ChapterList
	LastUpdatedDate string
}

func NewChapterInfoDat(user *User, chapter *Chapter, edgeDetailsInc []EdgeDetails,
	edgeDetailsOut []EdgeDetails) *ChapterInfoDat {
	chapList := make([]ChapterList, 0, len(edgeDetailsInc))
	for _, edge := range edgeDetailsInc {
		toChapter := edge.ToChapter
		lenList := len(chapList)
		if (lenList == 0) || (chapList[lenList-1].Index != toChapter.Index) {
			chapList = append(chapList, ChapterList{
				Index:    toChapter.Index,
				Fork:     chapter.StoryID != toChapter.StoryID,
				Chapters: []*ChapterInfoDat{NewChapterInfoDat(nil, toChapter, nil, nil)},
			})
		} else {
			chapList[lenList-1].Chapters = append(chapList[lenList-1].Chapters, NewChapterInfoDat(nil, toChapter,
				nil, nil))
		}
	}

	chapList2 := make([]ChapterList, 0, len(edgeDetailsOut))
	for _, edge := range edgeDetailsOut {
		fromChapter := edge.FromChapter
		lenList := len(chapList2)
		if (lenList == 0) || (chapList2[lenList-1].Index != fromChapter.Index) {
			chapList2 = append(chapList2, ChapterList{
				Index:    fromChapter.Index,
				Fork:     chapter.StoryID != fromChapter.StoryID,
				Chapters: []*ChapterInfoDat{NewChapterInfoDat(nil, fromChapter, nil, nil)},
			})
		} else {
			chapList2[lenList-1].Chapters = append(chapList2[lenList-1].Chapters, NewChapterInfoDat(nil, fromChapter,
				nil, nil))
		}
	}

	return &ChapterInfoDat{
		User:            user,
		ChapterInfo:     *chapter,
		EdgeDetailsInc:  chapList,
		EdgeDetailsOut:  chapList2,
		LastUpdatedDate: getTimeString(chapter.LastUpdated),
	}
}

type StoryPageDat struct {
	User     *User
	Story    *StoryInfoDat
	ChapList []ChapterList
}

func NewStoryPageDat(user *User, story *Story, chapterList []*Chapter) *StoryPageDat {
	chapList := make([]ChapterList, 0, len(chapterList))
	for _, chapter := range chapterList {
		lenList := len(chapList)
		if (lenList == 0) || (chapList[lenList-1].Index != chapter.Index) {
			chapList = append(chapList, ChapterList{
				Index:    chapter.Index,
				Chapters: []*ChapterInfoDat{NewChapterInfoDat(nil, chapter, nil, nil)},
			})
		} else {
			chapList[lenList-1].Chapters = append(chapList[lenList-1].Chapters, NewChapterInfoDat(nil, chapter, nil,
				nil))
		}
	}

	return &StoryPageDat{
		User:     user,
		Story:    NewStoryInfoDat(story),
		ChapList: chapList,
	}
}

type ChapterList struct {
	Index    int
	Fork     bool
	Chapters []*ChapterInfoDat
}
