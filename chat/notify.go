package chat

import (
	"encoding/json"

	"net/http"
)

// Notify sends notifications to the client based on the contents of the POST request
// support both JSON and HTML
func Notify(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		if r.Header.Get("Content-Type") == "application/json" {
			jq := make(map[string]string)
			// decode JSON into a string variable message
			json.NewDecoder(r.Body).Decode(&jq)
			msg := jq["message"]

			if msg == "" {
				http.Error(w, `Messages must have the key "messages"`, http.StatusBadRequest)
				return
			}
			// send the message to the hub
			DefaultHub.Echo <- msg
			w.Write([]byte("Sent message successfully"))
			return
		}

		// if not JSON, we assumed HTML, receive from POST form
		msg := r.FormValue("message")
		r.ParseForm()

		if msg == "" {
			http.Error(w, "No message found in request", http.StatusBadRequest)
			return
		}

		// also send the message to the hub
		DefaultHub.Echo <- msg
		w.Write([]byte("Sent message successfully"))
		return
	default:
		// only handle POST request
		http.Error(w, "Only POST method supported", http.StatusMethodNotAllowed)
		return
	}
}
