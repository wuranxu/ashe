package models

type TUserLog struct {
	ID       int `gorm:"primary_key;unique" json:"id"`
	SchoolId int `gorm:"type:varchar(16)" json:"school_id"` // 用户名
}

func (*TUserLog) TableName() string { return "t_user_log" }
