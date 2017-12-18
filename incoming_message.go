package messenger

const (
	MessageReceived   = "messages"
	AccountLinking    = "messaging_account_linking"
	Deliveries        = "message_deliveries"
	Echoes            = "message_echoes"
	Handovers         = "messaging_handovers"
	Optins            = "messaging_optins"
	Payments          = "messaging_payments"
	PolicyEnforcement = "messaging_policy_enforcement"
	Postbacks         = "messaging_postbacks"
	PreCheckouts      = "messaging_pre_checkouts"
	Reads             = "message_reads"
	Referrals         = "messaging_referrals"
	Standby           = "standby"
	UnknownEvent      = "unknown"
)

// IncomingEvent can be used as the destination for unmarshalling messaging event webhooks
type IncomingFacebookEvent struct {
	Object string `json:"object"`
	// Facebook returns messaging events as an array, but there's always only one element
	Entry         []*Entry `json:"entry"`
	EventCategory string
}

type Entry struct {
	ID        string            `json:"id"`
	Time      int64             `json:"time"`
	Messaging []*MessagingEvent `json:"messaging"`
}

type MessagingEvent struct {
	Sender struct {
		ID string `json:"id"`
	} `json:"sender"`
	Recipient struct {
		ID string `json:"id"`
	} `json:"recipient"`
	Timestamp   int64        `json:"timestamp"`
	Message     *Message     `json:"message"`
	Delivery    *Delivery    `json:"delivery"`
	Read        *MessageRead `json:"read"`
	AccountLink *AccountLink `json:"account_linking"`
	Handover    *Handover    `json:"pass_thread_control"`
	Optin       *Optin       `json:"optin"`
	Enforcement *Enforcement `json:"enforcement"`
	Postback    *Postback    `json:"postback"`
	Referral    *Referral    `json:"referral"`
}

type Message struct {
	IsEcho     bool   `json:"is_echo"`
	MID        string `json:"mid"`
	Text       string `json:"text"`
	QuickReply struct {
		Payload string `json:"payload"`
	} `json:"quick_reply"`
}

type Delivery struct {
	MIDS      []string `json:"mids"`
	Watermark int64    `json:"watermark"`
	Seq       int64    `json:"seq"`
}

type MessageRead struct {
	Watermark int64 `json:"watermark"`
	Seq       int64 `json:"seq"`
}

type AccountLink struct {
	Status            string `json:"status"`
	AuthorizationCode string `json:"authorization_code"`
}

type Handover struct {
	NewOwnerAppID string `json:"new_owner_app_id"`
	Metadata      string `json:"metadata"`
}

type Optin struct {
	Ref     string `json:"ref"`
	UserRef string `json:"user_ref"`
}

type Enforcement struct {
	Action string `json:"action"`
	Reason string `json:"reason"`
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

// SetEventCategory determines what type of API subscription this event is
// Facebook uses the same base for all messaging events with object, array of entries, etc.
// then within the messaging event will only include the relevant field
// e.g. A text message event will have a message field but not postback, and vice versa
func (event *IncomingFacebookEvent) SetEventCategory() {
	messagingEvent := event.Entry[0].Messaging[0]

	if messagingEvent.Message != nil && messagingEvent.Message.IsEcho {
		event.EventCategory = Echoes
		return
	}

	if messagingEvent.Message != nil && messagingEvent.Message.MID != "" {
		event.EventCategory = MessageReceived
		return
	}

	if messagingEvent.Delivery != nil && messagingEvent.Delivery.Watermark != 0 {
		event.EventCategory = Deliveries
		return
	}

	if messagingEvent.Read != nil && messagingEvent.Read.Watermark != 0 {
		event.EventCategory = Reads
		return
	}

	if messagingEvent.AccountLink != nil && messagingEvent.AccountLink.Status != "" {
		event.EventCategory = AccountLinking
		return
	}

	if messagingEvent.Handover != nil && messagingEvent.Handover.NewOwnerAppID != "" {
		event.EventCategory = Handovers
		return
	}

	if messagingEvent.Optin != nil && messagingEvent.Optin.Ref != "" {
		event.EventCategory = Optins
	}

	if messagingEvent.Enforcement != nil && messagingEvent.Enforcement.Reason != "" {
		event.EventCategory = PolicyEnforcement
		return
	}

	if messagingEvent.Postback != nil && messagingEvent.Postback.Title != "" {
		event.EventCategory = Postbacks
		return
	}

	if messagingEvent.Referral != nil && messagingEvent.Referral.Source != "" {
		event.EventCategory = Referrals
		return
	}

	event.EventCategory = UnknownEvent
}
