package handlers

import (
	"time"

	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service ServiceI
	valid   *validator.Validate
}

func New(s ServiceI, v *validator.Validate) *Handler {
	return &Handler{service: s, valid: v}
}

type CreateRequest struct {
	Message string         `json:"message" validate:"required"`
	SendAt  time.Time      `json:"send_at" validate:"required"`
	Retries int            `json:"retries" validate:"required"`
	To      string         `json:"to"      validate:"required"`
	Channel models.Channel `json:"channel" validate:"required"`
}
