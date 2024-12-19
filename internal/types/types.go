package types

type Student struct {
	Id    int64
	Name  string `json: "name" validate:"required"`
	Email string `json: "email" validate:"required"`
	Age   int    `json: "age"validate:"required"`
}
