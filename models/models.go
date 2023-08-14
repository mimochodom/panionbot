package models

import (
	"time"
)

type Users struct {
	UserID            int64  `gorm:"primaryKey"`
	UserName          string `gorm:"size:255"`
	BunnyCountGlobal  int
	TomatoCountGlobal int

	Groups []Groups `gorm:"many2many:users_groups;"`
}

type Groups struct {
	GroupID        int64  `gorm:"primaryKey"`
	GroupName      string `gorm:"size:255"`
	LastGamePlayed time.Time

	Users []Users `gorm:"many2many:users_groups;"`
}

type UsersGroups struct {
	UserID      int64
	GroupID     int64
	Leaving     bool
	BunnyCount  int
	TomatoCount int
	User        Users  `gorm:"foreignKey:UserID"`
	Group       Groups `gorm:"foreignKey:GroupID"`
}

type GroupsBTGameResult struct {
	id           int64 `gorm:"primaryKey"`
	GroupID      int64
	GamePlayed   time.Time
	UserIDBunny  int64
	UserIDTomato int64
	User         Users  `gorm:"foreignKey:UserIDBunny"`
	User_        Users  `gorm:"foreignKey:UserIDTomato"`
	Group        Groups `gorm:"foreignKey:GroupID"`
}

var SmileyList = []string{"🎲", "🎯", "🏀", "⚽", "🎳", "🎰"}
