package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/K1la/delayed-notifier/internal/api/response"
	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/K1la/delayed-notifier/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
	"time"
)

type Handler struct {
	service *service.NotificationService
	valid   *validator.Validate
}

func New(s *service.NotificationService, v *validator.Validate) *Handler {
	return &Handler{service: s, valid: v}
}

type CreateRequest struct {
	Message string         `json:"message" validate:"required"`
	SendAt  string         `json:"send_at" validate:"required"`
	Retries int            `json:"retries" validate:"required"`
	To      string         `json:"to"      validate:"required"`
	Channel models.Channel `json:"channel" validate:"required"`
}

func (h *Handler) CreateNotification(c *ginext.Context) {
	zlog.Logger.Info().Msgf("req: %+v", c.Request.Body)
	var req CreateRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to decode request body")
		response.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("decode error: %s", err.Error()))
		return
	}

	if err := h.valid.Struct(req); err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to validate request body")
		response.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("validation error: %s", err.Error()))
		return
	}

	// parse RFC3339 timestamp and convert to Moscow time
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		zlog.Logger.Fatal().Err(err).Msg("failed to load Moscow time zone")
	}

	parsedUTC, err := time.Parse(time.RFC3339, req.SendAt)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to parse send_at time")
		response.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("parse send_at time (use RFC3339): %s", err.Error()))
		return
	}
	parsedTime := parsedUTC.In(loc)

	notif := &models.Notification{
		ID:        uuid.New().String(),
		Message:   req.Message,
		Channel:   req.Channel,
		To:        req.To,
		SendAt:    parsedTime,
		Status:    models.StatusPending,
		Retries:   req.Retries,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = h.service.CreateNotification(c.Request.Context(), notif)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to create notification")
	}

	zlog.Logger.Info().Msgf("notif created: %+v", notif)
	response.Created(c.Writer, notif)

}

//	func (h *Handler) GetNotificationStatusByID(c *gin.Context) {
//		id := c.Param("id")
//
//		// TODO: заменить вместо in-memo на service
//		n, err := h.storage.Get(id)
//
//		if err != nil {
//			zlog.Logger.Error().Err(err).Msg("failed to find notifications")
//			response.NotFound(c.Writer, err)
//			return
//		}
//		response.OK(c.Writer, n)
//	}

func (h *Handler) GetAllNotifications(c *gin.Context) {
	list, err := h.service.GetNotifications(c.Request.Context())
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to find notifications")
		response.NotFound(c.Writer, err)
		return
	}
	response.OK(c.Writer, list)
}

//func (h *Handler) Delete(c *gin.Context) {
//	id := c.Param("id")
//
//	// TODO: заменить вместо in-memo на service
//	err := h.storage.Delete(id)
//	if err != nil {
//		response.NotFound(c.Writer, err)
//		return
//	}
//	response.OK(c.Writer, "canceled")
//}
