package messenger

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	LogStandard = iota
	LogDebug
)

type MessengerServer struct {
	VerificationToken   string
	MessageEventHandler func(*SubscriptionMessagingEvent)
	LogLevel int
}

// HandleRequestFromFacebook is the top-level http Handler for all requests coming from FB
// if the request is a GET, it's the verification check
// otherwise it is a subscription to a messaging event
func (m *MessengerServer) HandleRequestFromFacebook(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		m.HandleVerificationRequest(w, req)
	} else if req.Method == http.MethodPost {
		m.HandleIncomingEvent(w, req)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func NewMessengerServer(verificationSecret string) *MessengerServer {
	m := &MessengerServer{
		VerificationToken: verificationSecret,
	}

	// default to standard log level
	m.LogLevel = LogStandard

	return m
}

// HandleVerificationRequest responds to Facebook's verification for app subscriptions
func (m *MessengerServer) HandleVerificationRequest(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	mode := params.Get("hub.mode")
	token := params.Get("hub.verify_token")
	challenge := params.Get("hub.challenge")

	if len(mode) > 0 && len(token) > 0 {
		if mode == "subscribe" && token == m.VerificationToken && len(challenge) > 0 {
			io.WriteString(w, challenge)
		} else {
			w.WriteHeader(http.StatusForbidden)
		}
	}
}

// HandleIncomingEvent provides a default handler for doing something with events from Facebook
// unmarshalls the event body into a IncomingFacebookEvent struct then calls the MessageEventHandler
// with that parameter
func (m *MessengerServer) HandleIncomingEvent(w http.ResponseWriter, req *http.Request) {
	fbEvent := &IncomingFacebookEvent{}
	jsonBody, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Println("error reading incoming request body", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if m.LogLevel <= LogDebug {
		log.Println("REQUEST:")
		log.Println(string(jsonBody))
	}

	err = json.Unmarshal(jsonBody, fbEvent)
	if err != nil {
		log.Println("error unmarshalling json request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	// For all entries, handle the messaging event
	for _, entry := range fbEvent.Entry {
		// There's only one messaging event even though this is an array
		messagingEvent := entry.Messaging[0]
		messagingEvent.SetEventCategory()
		m.MessageEventHandler(messagingEvent)
	}
}
