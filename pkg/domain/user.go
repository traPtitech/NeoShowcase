package domain

func NewUser(name string) *User {
	return &User{
		ID:   NewID(),
		Name: name,
	}
}

type User struct {
	ID   string
	Name string
}
