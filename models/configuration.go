package models

type IConfiguration interface {
	GetSecret(name string) (*string, error)
}
