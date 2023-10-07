package domain

import "context"

type MidtransService interface {
	GenerateSnapURL(ctx context.Context, topUp *TopUp) error
	VerifyPayment(ctx context.Context, orderId string) (bool, error)
}
