package definitions

type ChatType string

var PrivateChatType ChatType = "private"
var GroupChatType ChatType = "group"
var SuperGroupChatType ChatType = "supergroup"
var ChannelChatType ChatType = "channel"

type Chat struct {
	Id       int      `json:"id"`
	Type     ChatType `json:"type"`
	Title    string   `json:"title"`
	Username string   `json:"username"`
}

type User struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Message struct {
	Id       int    `json:"message_id"`
	ThreadId int    `json:"message_thread_id"`
	From     *User  `json:"from"`
	Date     int    `json:"date"`
	Text     string `json:"text"`
	Chat     Chat   `json:"chat"`
}
