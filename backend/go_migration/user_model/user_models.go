package user_model

import "time"

// all data models for users db

type Session struct {
	IDSession int64     `gorm:"column:idSession;primaryKey;autoIncrement"`
	IDProfile int64     `gorm:"column:idProfile;index:fk_Session_Profile1_idx"`
	Start     time.Time `gorm:"column:start;default:CURRENT_TIMESTAMP"`
	End       time.Time `gorm:"column:end"`
}
func (Session) TableName() string { return "Session" }

type Profile struct {
	IDProfile   int64     `gorm:"column:idProfile;primaryKey;autoIncrement"`
	DateCreated time.Time `gorm:"column:date_created;default:CURRENT_TIMESTAMP"`
	DateUpdated time.Time `gorm:"column:date_updated;default:CURRENT_TIMESTAMP"`

	Sessions []Session `gorm:"foreignKey:IDProfile;references:IDProfile"`
}
func (Profile) TableName() string { return "Profile" }