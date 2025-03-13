package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/errs"
	"github.com/Zidan-Kharisma-Sakana/book-library/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

func routeWrapper(handler func(w http.ResponseWriter, r *http.Request) (interface{}, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqId := r.Header.Get("X-Request-Id")
		if reqId == "" {
			reqId = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), "reqId", reqId)

		data, err := handler(w, r.WithContext(ctx))

		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			w.Header().Set("X-Request-ID", reqId)

			var errorBuilder *errs.ErrorBuilder
			if !errors.As(err, &errorBuilder) {
				errorBuilder = errs.NewInternalServerError(err)
			}
			errorBuilder.RequestId = reqId

			logger.Error("Error", zap.Error(err), zap.Any("error", errorBuilder.Build()))

			w.WriteHeader(errorBuilder.StatusCode)
			if err := json.NewEncoder(w).Encode(errorBuilder.Build()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func decodeBody[T any](r *http.Request) (T, error) {
	var input T
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return input, err
	}
	return input, nil
}
