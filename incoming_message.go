package messenger

// IncomingEvent can be used as the destination for unmarshalling messaging event webhooks
type IncomingFacebookEvent struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID        string           `json:"id"`
	Time      int64            `json:"time"`
	Messaging []MessagingEvent `json:"messaging"`
}

type MessagingEvent struct {
	Sender struct {
		ID int64 `json:"id"`
	} `json:"sender"`
	Recipient struct {
		ID int64 `json:"id"`
	} `json:"recipient"`
	Timestamp int64    `json:"timestamp"`
	Message   Message  `json:"message"`
	Postback  Postback `json:"postback"`
}

type Message struct {
	MID        string `json:"mid"`
	Text       string `json:"text"`
	QuickReply struct {
		Payload string `json:"payload"`
	} `json:"quick_reply"`
}

type Postback struct {
	Title    string   `json:"title"`
	Payload  string   `json:"payload"`
	Referral Referral `json:"referral"`
}

type Referral struct {
	Ref    string `json:"ref"`
	Source string `json:"source"`
	Type   string `json:"type"`
}
