package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	RecipientID string `gorm:"unique;not null"`
	LastWorkout *time.Time
	StartsAt    *time.Duration
	EveryHours  *int

	Schedules      []Schedule
	LastSchedule   *Schedule `gorm:"foreignkey:id;association_foreignkey:last_schedule_id"`
	LastScheduleID *uint
}

func GetOrInitUserById(id string) *User {
	var user User

	db.FirstOrInit(&user, User{RecipientID: id})

	return &user
}

func (u *User) Recipient() string {
	return u.RecipientID
}
