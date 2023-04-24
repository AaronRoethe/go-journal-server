package api

import (
	"fmt"
	"io"
	"net/http"

	"github.com/AaronRoethe/go-journal-server/src/storage"
)

func MessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body into a byte buffer
	body := make([]byte, r.ContentLength)
	if _, err := io.ReadFull(r.Body, body); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Save the message to the blob
	err := storage.SaveMessageToBlob(body)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Message saved successfully")
}
