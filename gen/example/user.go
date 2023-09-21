package example

import "database/sql"

const (
	UserId       = "Id"
	UserAge      = "Age"
	UserName     = "UserName"
	UserNickName = "NickName"
)

type User struct {
	Id       uint64
	Age      *int
	Name     string
	NickName *sql.NullString
}

const (
	UserDetailId      = "Id"
	UserDetailUserId  = "UserId"
	UserDetailAddress = "Address"
	UserDetailPicture = "Picture"
)

type UserDetail struct {
	Id      uint64
	UserId  uint64
	Address string
	Picture []byte
}
