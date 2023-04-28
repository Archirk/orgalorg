package main

import (
	// "bytes"
	// "encoding/json"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/google/go-querystring/query"
)
const (
	testChat = 144581052
)
type (
	Msg struct {
		ChatID int `url:"chat_id"`
		Text string `url:"text"`
	}
	GitlabEvent struct {
		EventName string `json:"event_name"`
	}
)
func sendMessage(text string) {
	client := http.DefaultClient
	q, _ := query.Values(Msg{testChat, text})
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TToken"))
	r, _ := http.NewRequest(http.MethodGet, url, nil)

	r.URL.RawQuery = q.Encode()
	_, err := client.Do(r)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("sent")
}

func handle(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var glEvent GitlabEvent
	if err = json.Unmarshal(body, &glEvent); err != nil {
		panic(err)
	}
	sendMessage(glEvent.EventName)
}

func main() {
	router := chi.NewRouter()
	router.Post("/", handle)
	http.ListenAndServe(":9000", router)
}