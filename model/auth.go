package model

type Users struct {
	IdUser   uint64 `gorm:"column:id;primaryKey;autoIncrement"`
	username string `gorm:"column:username"`
	password string `gorm:"column:password"`
}

type Logins struct {
	username string
	password string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type Todo struct {
	UserID uint64 `json:"user_id"`
	Title  string `json:"title"`
}
