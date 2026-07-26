package models

import "strconv"

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
