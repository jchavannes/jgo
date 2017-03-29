package chat

type Message struct {
	Id       uint
	Date     int64
	Message  string
	User     *User
}
