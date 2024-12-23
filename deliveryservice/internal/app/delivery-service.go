package app

import (
	"context"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/logger"
	"github.com/IvanChumakov/hotel-booking-project/hotel-lib/tracing"
	pb "github.com/IvanChumakov/hotel-booking-project/protos"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotificationDeliveryServer struct {
	pb.UnimplementedNotificationDeliveryServer
}

func NewServer() *NotificationDeliveryServer {
	return &NotificationDeliveryServer{}
}

func (n *NotificationDeliveryServer) SendNotification(ctx context.Context, notification *pb.Notification) (*emptypb.Empty, error) {
	log := logger.New()
	ctx, err := tracing.GetParentContext(ctx)
	if err != nil {
		log.Logger.Error("Error getting parent context")
	}
	_, span := tracing.StartTracerSpan(ctx, "delivery-system")
	defer span.End()

	log.Logger.Info("delivery system got notification")
	return &emptypb.Empty{}, nil
}
