package telegram

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type UpdateResponse struct {
	Result []Update
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}
