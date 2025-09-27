package handlers

import (
	"errors"
	"fmt"
	"github.com/K1la/delayed-notifier/internal/api/response"
	"github.com/K1la/delayed-notifier/internal/models"
	"github.com/K1la/delayed-notifier/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"time"
)

func (h *Handler) CreateNotification(c *ginext.Context) {
	zlog.Logger.Info().Msgf("req: %+v", c.Request.Body)
	var notifReq CreateRequest
	if err := c.BindJSON(&notifReq); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to decode request body")
		response.Fail(c.Writer, fmt.Errorf("decode error: %s", err.Error()))
		return
	}

	zlog.Logger.Info().Msgf("after json decoder notifReq: %v", notifReq)

	if err := h.valid.Struct(&notifReq); err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to validate request body")
		response.Fail(c.Writer, fmt.Errorf("validation error: %s", err.Error()))
		return
	}

	zlog.Logger.Info().Msgf("after valid struct notifReq: %+v", notifReq)

	//if time.Until(notifReq.SendAt) <= 0 {
	//	zlog.Logger.Error().Msg("invalid payload: time is in the past")
	//	response.Fail(c.Writer, fmt.Errorf("invalid payload: time shuold be in the future"))
	//	return
	//}

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
		response.NotFound(c.Writer, err)
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
			response.Fail(c.Writer, err)
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
			response.Fail(c.Writer, err)
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

	err := h.service.UpdateNotification(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrNotificationNotFound) {
			zlog.Logger.Error().Err(err).Msg("failed to cancel notification")
			response.Fail(c.Writer, err)
			return
		}
		response.NotFound(c.Writer, err)
		return
	}
	zlog.Logger.Info().Msgf("successfuly handled cancel notification with id:%v", id)
	response.OK(c.Writer, "canceled")
}
