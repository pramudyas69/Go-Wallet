package service

import (
	"bytes"
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
	"errors"
	"text/template"
	"time"
)

type notificationService struct {
	notificationRepository domain.NotificationRepository
	templateRepository     domain.TemplateRepository
	hub                    *dto.Hub
}

func NewNotification(notificationRepository domain.NotificationRepository,
	templateRepository domain.TemplateRepository,
	hub *dto.Hub) domain.NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
		templateRepository:     templateRepository,
		hub:                    hub,
	}
}

func (n notificationService) FindByUser(ctx context.Context, userID int64) ([]dto.NotificationData, error) {
	notifications, err := n.notificationRepository.FindByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var notificationData []dto.NotificationData
	for _, notification := range notifications {
		notificationData = append(notificationData, dto.NotificationData{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Title:     notification.Title,
			Body:      notification.Body,
			Status:    notification.Status,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt,
		})
	}

	return notificationData, nil
}

func (n notificationService) Insert(ctx context.Context, userId int64, code string, data map[string]string) error {
	tmpl, err := n.templateRepository.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if tmpl == (domain.Template{}) {
		return errors.New("template not found")
	}

	body := new(bytes.Buffer)
	t := template.Must(template.New("notification").Parse(tmpl.Body))
	if err := t.Execute(body, data); err != nil {
		return err
	}

	notification := domain.Notification{
		UserID:    userId,
		Title:     tmpl.Title,
		Body:      body.String(),
		Status:    1,
		IsRead:    0,
		CreatedAt: time.Now(),
	}
	err = n.notificationRepository.Insert(ctx, &notification)
	if err != nil {
		return err
	}

	if channel, ok := n.hub.NotificationChannel[userId]; ok {
		channel <- dto.NotificationData{
			ID:        notification.ID,
			UserID:    notification.UserID,
			Title:     notification.Title,
			Body:      notification.Body,
			Status:    notification.Status,
			IsRead:    notification.IsRead,
			CreatedAt: notification.CreatedAt,
		}
	}

	return nil
}
