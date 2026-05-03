package authx

type Claims struct {
	UserID          string
	Email           string
	Roles           []string
	PasswordVersion int
}
