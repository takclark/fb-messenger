package messenger

import "encoding/json"

const (
	EventCategoryMessageReceived   = "messages"
	EventCategoryAccountLinking    = "messaging_account_linking"
	EventCategoryDeliveries        = "message_deliveries"
	EventCategoryEchoes            = "message_echoes"
	EventCategoryHandovers         = "messaging_handovers"
	EventCategoryOptins            = "messaging_optins"
	EventCategoryPayments          = "messaging_payments"
	EventCategoryPolicyEnforcement = "messaging_policy_enforcement"
	EventCategoryPostbacks         = "messaging_postbacks"
	EventCategoryPreCheckouts      = "messaging_pre_checkouts"
	EventCategoryReads             = "message_reads"
	EventCategoryReferrals         = "messaging_referrals"
	EventCategoryStandby           = "standby"
	EventCategoryUnknownEvent      = "unknown"
)

// IncomingEvent can be used as the destination for unmarshalling messaging event webhooks
type IncomingFacebookEvent struct {
	Object string `json:"object"`
	// Facebook returns messaging events as an array, there may be multiple
	Entry []*Entry `json:"entry"`
}

type Entry struct {
	ID        string                        `json:"id"`
	Time      int64                         `json:"time"`
	// Though this is an array, Facebook always sends only one
	Messaging []*SubscriptionMessagingEvent `json:"messaging"`
}

type Recipient struct {
	ID string `json:"id"`
}

type SubscriptionMessagingEvent struct {
	EventCategory string
	Sender        struct {
		ID string `json:"id"`
	} `json:"sender"`
	Recipient   *Recipient                  `json:"recipient"`
	Timestamp   int64                       `json:"timestamp"`
	Message     *SubscriptionMessageContent `json:"message"`
	Delivery    *Delivery                   `json:"delivery"`
	Read        *MessageRead                `json:"read"`
	AccountLink *AccountLink                `json:"account_linking"`
	Handover    *Handover                   `json:"pass_thread_control"`
	Optin       *Optin                      `json:"optin"`
	Enforcement *Enforcement                `json:"enforcement"`
	Postback    *Postback                   `json:"postback"`
	Referral    *Referral                   `json:"referral"`
}

type SubscriptionMessageContent struct {
	IsEcho      bool                      `json:"is_echo"`
	MID         string                    `json:"mid"`
	Text        string                    `json:"text"`
	Attachments []*SubscriptionAttachment `json:"attachments"`
	QuickReply  struct {
		Payload string `json:"payload"`
	} `json:"quick_reply"`
}

type SubscriptionAttachment struct {
	Type    string                         `json:"type"`
	Payload *SubscriptionAttachmentPayload `json:"payload"`
	Title   string                         `json:"title"`
	URL     string                         `json:"URL"`
}

type SubscriptionAttachmentPayload struct {
	URL         string                                    `json:"url"`
	Coordinates *SubscriptionCoordinatesAttachmentPayload `json:"coordinates"`
}

type SubscriptionCoordinatesAttachmentPayload struct {
	Lat  json.Number `json:"lat"`
	Long json.Number `json:"long"`
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

// GetSenderID returns the sender ID of a messaging event, useful to determine where to send reply
func (event *SubscriptionMessagingEvent) GetSenderID() string {
	return event.Sender.ID
}

// SetEventCategory determines what type of API subscription a messaging event is
// Facebook uses the same base for messaging events then within the messaging event will only include the relevant field
// e.g. A text message event will have a message field but not postback, and vice versa
func (event *SubscriptionMessagingEvent) SetEventCategory() {

	if event.Message != nil && event.Message.IsEcho {
		event.EventCategory = EventCategoryEchoes
		return
	}

	if event.Message != nil && event.Message.MID != "" {
		event.EventCategory = EventCategoryMessageReceived
		return
	}

	if event.Delivery != nil && event.Delivery.Watermark != 0 {
		event.EventCategory = EventCategoryDeliveries
		return
	}

	if event.Read != nil && event.Read.Watermark != 0 {
		event.EventCategory = EventCategoryReads
		return
	}

	if event.AccountLink != nil && event.AccountLink.Status != "" {
		event.EventCategory = EventCategoryAccountLinking
		return
	}

	if event.Handover != nil && event.Handover.NewOwnerAppID != "" {
		event.EventCategory = EventCategoryHandovers
		return
	}

	if event.Optin != nil && event.Optin.Ref != "" {
		event.EventCategory = EventCategoryOptins
	}

	if event.Enforcement != nil && event.Enforcement.Reason != "" {
		event.EventCategory = EventCategoryPolicyEnforcement
		return
	}

	if event.Postback != nil && event.Postback.Title != "" {
		event.EventCategory = EventCategoryPostbacks
		return
	}

	if event.Referral != nil && event.Referral.Source != "" {
		event.EventCategory = EventCategoryReferrals
		return
	}

	event.EventCategory = EventCategoryUnknownEvent
}
