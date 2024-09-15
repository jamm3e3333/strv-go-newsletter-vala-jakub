package create_user

import "fmt"

type CreateUser struct {
}

func NewCreateUserHandler() *CreateUser {
	return &CreateUser{}
}

func (h *CreateUser) Handle(c *Command) error {
	fmt.Println(c.Email, c.Password)
	return nil
}
