package data

type UsersRepo interface {
	GetAll() Users
	GetUsers() Users
	AddUser(u *User)
	PutUser(u *User, id int) error
	DeleteUser(id int) error
}
