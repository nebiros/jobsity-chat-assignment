package response

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type HTTPError struct {
	Type       string `json:"type"`
	Title      string `json:"title"`
	StatusCode int    `json:"status"`
	Detail     string `json:"detail"`
	Instance   string `json:"instance"`

	RequestID string `json:"requestId"`
}

func NewHTTPError(ctx context.Context, statusCode int, detail string) *HTTPError {
	reqID := middleware.GetReqID(ctx)

	return &HTTPError{
		Type:       "",
		Title:      http.StatusText(statusCode),
		StatusCode: statusCode,
		Detail:     detail,
		Instance:   "",
		RequestID:  reqID,
	}
}

func (e HTTPError) Error() string {
	return fmt.Sprintf(
		"HTTP error '%d' %s with request ID: %s: %s",
		e.StatusCode,
		e.Title,
		e.RequestID,
		e.Detail,
	)
}

func WriteHTTPError(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
	writeError(ctx, w, NewHTTPError(ctx, statusCode, err.Error()))
}

func writeError(ctx context.Context, w http.ResponseWriter, httpError *HTTPError) {
	b, err := json.Marshal(httpError)
	if err != nil {
		slog.ErrorContext(ctx, "unable to marshal error", slog.Any("error", err))

		return
	}

	slog.ErrorContext(ctx, "HTTP error", slog.String("error", string(b)))

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(httpError.StatusCode)

	if _, err := w.Write(b); err != nil {
		slog.ErrorContext(ctx, "http.ResponseWriter.Write errored", slog.Any("error", err))
	}
}
