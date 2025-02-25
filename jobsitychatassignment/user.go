package jobsitychatassignment

type User struct {
	ID             string `sql:"id"`
	Username       string `sql:"username"`
	HashedPassword string `sql:"hashedPassword"`
}
