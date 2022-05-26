package user

type User struct {
	ID           string
	Login        string
	PasswordHash []byte
	Token        string
}
