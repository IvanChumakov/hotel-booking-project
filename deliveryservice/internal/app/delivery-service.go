package app

import (
	"context"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	pb "github.com/IvanChumakov/hotel-booking-project/protos"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotificationDeliveryServer struct {
	pb.UnimplementedNotificationDeliveryServer
}

func NewServer() *NotificationDeliveryServer {
	return &NotificationDeliveryServer{}
}

func (n *NotificationDeliveryServer) SendNotification(_ context.Context, notification *pb.Notification) (*emptypb.Empty, error) {
	log := logger.New()
	log.Logger.Info("delivery system got notification")
	return &emptypb.Empty{}, nil
}
