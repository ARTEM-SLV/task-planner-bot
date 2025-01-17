package database

type Repository interface {
	AddUser(userID int64, username string, lastMsgID int) error
}
