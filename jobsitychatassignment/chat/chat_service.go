package chat

import (
	"encoding/csv"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	domain "github.com/nebiros/jobsity-chat-assignment/jobsitychatassignment"
	"github.com/nebiros/jobsity-chat-assignment/pkg/websocket"
)

var (
	_hubs = make(map[string]*websocket.Hub[domain.ChatMessage])
)

var (
	_stocksServerBaseURL = "https://stooq.com"
)

type Service struct {
	httpClient *http.Client
}

func NewService(httpClient *http.Client) *Service {
	return &Service{httpClient: httpClient}
}

func (s *Service) InitWebSocket(hubID string, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	if _, ok := _hubs[hubID]; !ok {
		slog.Debug("new websocket hub created", slog.String("hubID", hubID))

		_hubs[hubID] = websocket.NewHub[domain.ChatMessage](ctx)
		go _hubs[hubID].Run()
	}

	conn, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	defer conn.Close()

	wsc := websocket.NewClient[domain.ChatMessage](ctx, _hubs[hubID], conn)

	slog.Info("registering websocket client", slog.String("hubID", hubID))
	_hubs[hubID].Register() <- wsc

	go wsc.WritePump()
	wsc.ReadPump(s.spyMessage)

	return nil
}

func (s *Service) spyMessage(broadcast chan domain.ChatMessage, chatMessage domain.ChatMessage) {
	message := strings.TrimSpace(chatMessage.Message)
	if message == "" {
		return
	}

	if strings.HasPrefix(message, "/") {
		// in: /stock=aapl.us
		comm := message[len("/"):]
		// out: stock=aapl.us

		// in: stock=aapl.us
		ss := strings.Split(comm, "=")
		if len(ss) != 2 {
			return
		}
		// out: stock aapl.us

		commandType, commandValue := ss[0], ss[1]

		slog.Debug("comm", slog.String("commandType", commandValue), slog.String("commandValue", commandValue))

		switch commandType {
		case "stock":
			u := fmt.Sprintf("%s/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", _stocksServerBaseURL, commandValue)

			slog.Debug("stocks server get endpoint url", slog.String("url", u))

			resp, err := s.httpClient.Get(u)
			if err != nil {
				broadcast <- domain.ChatMessage{
					UserID:   "",
					Username: "StockServ",
					Message:  err.Error(),
				}

				return
			}

			defer resp.Body.Close()

			csvReader := csv.NewReader(resp.Body)

			header, err := csvReader.Read()
			if err != nil {
				broadcast <- domain.ChatMessage{
					UserID:   "",
					Username: "StockServ",
					Message:  fmt.Errorf("unable to read csv header: %w", err).Error(),
				}

				return
			}

			data := make(map[string]string, 0)

			for {
				row, err := csvReader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					broadcast <- domain.ChatMessage{
						UserID:   "",
						Username: "StockServ",
						Message:  err.Error(),
					}

					return
				}

				for i, col := range header {
					data[col] = row[i]
				}
			}

			// APPL.US quote is $93.42 per share
			result := data["Symbol"] + " quote is $" + data["Close"] + " per share"

			broadcast <- domain.ChatMessage{
				UserID:   "",
				Username: "StockServ",
				Message:  result,
			}
		}
	}
}
