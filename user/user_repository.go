package user

type UserRepository interface {
	FindUser(username string) (*User, error)
	ListUsers() ([]User, error)
	DeleteUser(username string) (*User, error)
}
