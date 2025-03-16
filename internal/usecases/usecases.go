package usecases

import "net/http"

type Auth interface {
	SignUp(username, password string) (string, error)
	SignIn(username, password string) (string, error)
}

type Chat interface {
	HandleMessages()
	HandleConnections(w http.ResponseWriter, r *http.Request)
}
