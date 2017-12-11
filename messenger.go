package server

import (
	"io"
	"net/http"
)

type MessengerServer struct {
	VerificationToken string
	PostHandler       func(http.ResponseWriter, *http.Request)
}

// HandleRequestFromFacebook is the top-level http Handler for all requests coming from FB
// if the request is a GET, it's the verification check
// otherwise it is a subscription to a messaging event
func (m *MessengerServer) HandleRequestFromFacebook(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		m.HandleVerificationRequest(w, req)
	} else if req.Method == http.MethodPost {
		m.PostHandler(w, req)
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
