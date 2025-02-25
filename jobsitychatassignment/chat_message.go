package jobsitychatassignment

type ChatMessage struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Message  string `json:"message"`
}
