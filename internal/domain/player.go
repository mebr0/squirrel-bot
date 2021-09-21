package domain

type Player struct {
	ID int64 `db:"id"`
	Username string `db:"username"`
	FirstName string `db:"first_name"`
	LastName string `db:"last_name"`
}

func NewPlayer(id int64, username, firstName, lastName string) Player {
	return Player{
		ID: id,
		Username: username,
		FirstName: firstName,
		LastName: lastName,
	}
}
