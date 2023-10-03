package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/dto"
)

type notificationService struct {
	notificationRepository domain.NotificationRepository
}

func NewNotification(notificationRepository domain.NotificationRepository) domain.NotificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
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
