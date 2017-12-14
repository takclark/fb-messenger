package messenger

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type MessengerServer struct {
	VerificationToken   string
	MessageEventHandler func(*IncomingFacebookEvent)
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

func NewMessengerClient(verificationSecret string) *MessengerServer {
	m := &MessengerServer{
		VerificationToken: verificationSecret,
	}

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

	log.Println("REQUEST:")
	log.Println(string(jsonBody))
	err = json.Unmarshal(jsonBody, fbEvent)
	if err != nil {
		log.Println("error unmarshalling json request", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	// Set the type of this event to be helpful to the handler
	fbEvent.SetEventCategory()

	m.MessageEventHandler(fbEvent)
}
