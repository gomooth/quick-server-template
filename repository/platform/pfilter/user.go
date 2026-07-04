package pfilter

type User struct {
	Account     string
	AccountLike string
}

type UserLoginLog struct {
	UserID  uint
	UserIDs []uint
}
