package domain

func NewUser(name string) *User {
	return &User{
		ID:    NewID(),
		Name:  name,
		Admin: false,
	}
}

type User struct {
	ID    string
	Name  string
	Admin bool
}
