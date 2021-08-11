package database

type User struct {
	Uid       uint64
	Uname     string
	Upassword string
	Uadmin    bool
}
