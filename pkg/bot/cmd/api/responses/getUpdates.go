package responses

type Message struct {
	Id       int    `json:"message_id"`
	ThreadId int    `json:"message_thread_id"`
	From     *User  `json:"from"`
	Date     int    `json:"date"`
	Text     string `json:"text"`
}

type User struct {
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type InlineQuery struct {
	Id       string `json:"id"`
	From     *User  `json:"user"`
	Query    string `json:"query"`
	ChatType string `json:"chat_type"`
}

type Update struct {
	UpdateId    int          `json:"update_id"`
	Message     *Message     `json:"message"`
	InlineQuery *InlineQuery `json:"inline_query"`
}

type GetUpdatesResponse struct {
	Ok     bool      `json:"ok"`
	Result []*Update `json:"result"`
}
