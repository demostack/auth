package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Render JSON.
func Render(w http.ResponseWriter, status int, data string) {
	type Response struct {
		Status int    `json:"status"`
		Data   string `json:"data"`
	}

	res := Response{}
	res.Status = status
	res.Data = data

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)

	b, _ := json.Marshal(res)
	fmt.Fprintf(w, string(b))
}

// RenderAccessKeys JSON.
func RenderAccessKeys(w http.ResponseWriter, status int, accessKeyID string,
	secretAccessKey string, sessionToken string) {
	type Response struct {
		Status          int    `json:"status"`
		AccessKeyID     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		SessionToken    string `json:"session_token"`
	}

	res := Response{}
	res.Status = status
	res.AccessKeyID = accessKeyID
	res.SecretAccessKey = secretAccessKey
	res.SessionToken = sessionToken

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Status)

	b, _ := json.Marshal(res)
	fmt.Fprintf(w, string(b))
}
