package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"pawang-backend/config"
	"pawang-backend/entity"
	"pawang-backend/repository"

	"github.com/OneSignal/onesignal-go-api"
)

type NotificationService interface {
	AddUserNotification(userId int, playerId string) (entity.UserNotification, error)
	DeleteUserNotification(userId int, playerId string) error
	SendNotification(title, subtitle, playerId string)
}

type notificationService struct {
	userNotificationRepository repository.UserNotificationRepository
}

func NewNotificationService(userNotificationRepository repository.UserNotificationRepository) *notificationService {
	return &notificationService{userNotificationRepository}
}

func (service *notificationService) SendNotification(title, subtitle, playerId string) {
	notification := *onesignal.NewNotification(config.GetEnv("ONESIGNAL_APP_ID"))
	notification.Name = onesignal.PtrString("Testing Notif")
	notification.SetDeliveryTimeOfDay("11:05")
	notification.SetHeadings(onesignal.StringMap{En: onesignal.PtrString(title)})
	notification.SetContents(onesignal.StringMap{En: onesignal.PtrString(subtitle)})
	notification.SetIncludePlayerIds([]string{playerId})

	configuration := onesignal.NewConfiguration()
	apiClient := onesignal.NewAPIClient(configuration)

	appAuth := context.WithValue(context.Background(), onesignal.AppAuth, config.GetEnv("ONESIGNAL_APP_KEY"))

	resp, r, err := apiClient.DefaultApi.CreateNotification(appAuth).Notification(notification).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreateNotification``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreateNotification`: CreateNotificationSuccessResponse
	fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreateNotification`: %v\n", resp)
}

func (service *notificationService) AddUserNotification(userId int, playerId string) (entity.UserNotification, error) {
	notification := entity.UserNotification{}
	notification.UserID = userId
	notification.OnesignalID = playerId

	player := *onesignal.NewPlayer(playerId, int32(1))
	player.SetAppId(config.GetEnv("ONESIGNAL_APP_ID"))

	configuration := onesignal.NewConfiguration()
	apiClient := onesignal.NewAPIClient(configuration)

	appAuth := context.WithValue(context.Background(), onesignal.AppAuth, config.GetEnv("ONESIGNAL_APP_KEY"))

	resp, r, err := apiClient.DefaultApi.CreatePlayer(appAuth).Player(player).Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultApi.CreatePlayer``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `CreatePlayer`: CreatePlayerSuccessResponse
	fmt.Fprintf(os.Stdout, "Response from `DefaultApi.CreatePlayer`: %v\n", resp)

	newUserNotification, err := service.userNotificationRepository.Create(notification)
	if err != nil {
		return newUserNotification, err
	}

	return newUserNotification, nil
}

func (service *notificationService) DeleteUserNotification(userId int, playerId string) error {
	notification, err := service.userNotificationRepository.FindByOneSignalID(playerId)
	if err != nil {
		return err
	}

	if notification.UserID != userId {
		return errors.New("the notification does no exists")
	}

	if err = service.userNotificationRepository.Delete(notification); err != nil {
		return err
	}

	return nil
}
