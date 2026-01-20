package main

import (
	"github.com/jinzhu/gorm"
)

type Token struct {
	gorm.Model
	Account   Account
	AccountID uint
	Token     string
}

type Account struct {
	gorm.Model
	User      string `gorm:"unique"`
	Email     string `gorm:"unique"`
	Pass      string `json:"-"`
	Favorite  int
	Validated bool
	Admin     bool
}

type List struct {
	gorm.Model
	Name      string
	Items     []Item
	Archived  bool
	AccountID uint
	Account   Account
	Share     []Account `gorm:"many2many:share_list;"`
}

type Item struct {
	gorm.Model
	Order           int
	ListId          int
	Title           string
	Note            string
	Complete        bool
	Completed       int
	Time            string
	ReminderSent    bool
	RecurDays       int  // 0 = no recurrence, >0 = repeat every N days
	PreAlarmMinutes int  // 0 = no pre-alarm, >0 = minutes before
	PreAlarmSent    bool // track if pre-alarm was sent
	AccountID       uint // could be a co-creator
	Account         Account
}
