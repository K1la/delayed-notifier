package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/K1la/delayed-notifier/internal/api/response"
	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/K1la/delayed-notifier/internal/repository"
	"github.com/K1la/delayed-notifier/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
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
	SendAt  time.Time      `json:"send_at" validate:"required"`
	Retries int            `json:"retries" validate:"required"`
	To      string         `json:"to"      validate:"required"`
	Channel models.Channel `json:"channel" validate:"required"`
}

func (h *Handler) CreateNotification(c *ginext.Context) {
	zlog.Logger.Info().Msgf("req: %+v", c.Request.Body)
	var notifReq CreateRequest
	//if err := json.NewDecoder(c.Request.Body).Decode(&notifReq); err != nil {
	if err := c.BindJSON(&notifReq); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to decode request body")
		response.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("decode error: %s", err.Error()))
		return
	}

	zlog.Logger.Info().Msgf("after json decoder notifReq: %v", notifReq)

	if err := h.valid.Struct(&notifReq); err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to validate request body")
		response.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("validation error: %s", err.Error()))
		return
	}

	zlog.Logger.Info().Msgf("after valid struct notifReq: %+v", notifReq)

	//// parse RFC3339 timestamp and convert to Moscow time
	//loc, err := time.LoadLocation("Europe/Moscow")
	//if err != nil {
	//	zlog.Logger.Fatal().Err(err).Msg("failed to load Moscow time zone")
	//	response.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("failed to load Moscow time zone: %s", err.Error()))
	//	return
	//}
	//
	//parsedUTC, err := time.Parse(time.RFC3339, req.SendAt)
	//if err != nil {
	//	zlog.Logger.Error().Err(err).Msg("failed to parse send_at time")
	//	response.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("parse send_at time (use RFC3339): %s", err.Error()))
	//	return
	//}
	//parsedTime := parsedUTC.In(loc)

	if time.Until(notifReq.SendAt) <= 0 {
		zlog.Logger.Error().Msg("invalid payload: time is in the past")
		response.Fail(c.Writer, http.StatusBadRequest, fmt.Errorf("invalid payload: time shuold be in the future"))
		return
	}

	notif := &models.Notification{
		ID:        uuid.New().String(),
		Message:   notifReq.Message,
		Channel:   notifReq.Channel,
		To:        notifReq.To,
		SendAt:    notifReq.SendAt,
		Status:    models.StatusPending,
		Retries:   notifReq.Retries,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	notif, err := h.service.CreateNotification(c.Request.Context(), notif)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to create notification")
		response.Fail(c.Writer, http.StatusInternalServerError, err)
		return
	}

	zlog.Logger.Info().Msgf("successfuly created notification: %+v", notif)
	response.Created(c.Writer, notif)

}

func (h *Handler) GetNotificationStatusByID(c *gin.Context) {
	zlog.Logger.Info().Msgf("req: %+v", c.Request.Body)

	id := c.Param("id")

	status, err := h.service.GetNotificationStatusByID(c.Request.Context(), id)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to find notification status")
		if errors.Is(err, repository.ErrNotificationNotFound) {
			response.Fail(c.Writer, http.StatusBadRequest, err)
		}
		response.NotFound(c.Writer, err)
		return
	}

	zlog.Logger.Info().Msgf("successfuly get status notification (id=%v) : %+v", id, status)
	response.OK(c.Writer, status)
}

func (h *Handler) GetAllNotifications(c *gin.Context) {
	list, err := h.service.GetNotifications(c.Request.Context())
	if err != nil {
		if errors.Is(err, repository.ErrNoNotificationsFound) {
			zlog.Logger.Error().Err(err).Msg("failed to find all notifications")
			response.NotFound(c.Writer, err)
			return
		}
		zlog.Logger.Error().Err(err).Msg("failed to find notifications")
		response.NotFound(c.Writer, err)
		return
	}

	response.OK(c.Writer, list)
}

func (h *Handler) CancelNotification(c *gin.Context) {
	id := c.Param("id")

	err := h.service.CancelNotification(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationNotFound) {
			zlog.Logger.Error().Err(err).Msg("failed to cancel notification")
			response.Fail(c.Writer, http.StatusBadRequest, err)
			return
		}
		response.NotFound(c.Writer, err)
		return
	}
	zlog.Logger.Info().Msgf("successfuly handled cancel notification with id:%v", id)
	response.OK(c.Writer, "canceled")
}
