package chat

import (
	"encoding/csv"
	"net/http"
	"net/http/httptest"
	"testing"

	domain "github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment"
)

func TestService_spyMessage(t *testing.T) {
	tests := map[string]struct {
		in struct {
			broadcast   chan domain.ChatMessage
			chatMessage domain.ChatMessage
		}
		out struct {
			chatMessage domain.ChatMessage
		}
	}{
		"should accept /stock command": {
			in: struct {
				broadcast   chan domain.ChatMessage
				chatMessage domain.ChatMessage
			}{
				broadcast: make(chan domain.ChatMessage),
				chatMessage: domain.ChatMessage{
					UserID:   "",
					Username: "",
					Message:  "/stock=aapl.us",
				},
			},
			out: struct{ chatMessage domain.ChatMessage }{chatMessage: domain.ChatMessage{
				UserID:   "",
				Username: "StockServ",
				Message:  "AAPL.US quote is $247.1725 per share",
			}},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				records := [][]string{
					{"Symbol", "Date", "Time", "Open", "High", "Low", "Close", "Volume"},
					{"AAPL.US", "2025-02-24", "22:00:19", "244.925", "248.86", "244.62", "247.1725", "51207835"},
				}

				csvWriter := csv.NewWriter(w)
				if err := csvWriter.WriteAll(records); err != nil {
					t.Fatal(err)
				}

				w.Header().Set("Content-Type", "text/csv")
			}))

			defer testServer.Close()

			_stocksServerBaseURL = testServer.URL

			s := &Service{
				httpClient: testServer.Client(),
			}

			go s.spyMessage(tt.in.broadcast, tt.in.chatMessage)

			select {
			case message, ok := <-tt.in.broadcast:
				if !ok {
					t.Fatal("broadcast channel closed")
				}

				if tt.out.chatMessage.Username != message.Username {
					t.Errorf("tt.out.chatMessage.Username != \"%s\"", message.Username)
				}

				if tt.out.chatMessage.Message != message.Message {
					t.Errorf("tt.out.chatMessage.Message != \"%s\"", message.Message)
				}
			}

			close(tt.in.broadcast)
		})
	}
}
