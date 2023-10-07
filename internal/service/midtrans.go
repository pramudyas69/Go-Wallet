package service

import (
	"context"
	"e-wallet/domain"
	"e-wallet/internal/config"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type midtransService struct {
	config config.Midtrans
	env    midtrans.EnvironmentType
}

func NewMidtransService(cnf *config.Config) domain.MidtransService {
	env := midtrans.Sandbox
	if cnf.Midtrans.IsProd {
		env = midtrans.Production
	}

	return &midtransService{
		config: cnf.Midtrans,
		env:    env,
	}
}

func (m midtransService) GenerateSnapURL(ctx context.Context, topUp *domain.TopUp) error {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  topUp.ID,
			GrossAmt: int64(topUp.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "Tian Presti",
			LName: "Herlina",
			Email: "pandupramudya44@gmail.com",
			Phone: "085339040823",
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "Hutang Ke Pacar",
				Price: int64(topUp.Amount),
				Name:  "Hutang Ke Pacar",
				Qty:   1,
			},
		},
	}

	var client snap.Client
	client.New(m.config.Key, m.env)

	snapResp, err := client.CreateTransaction(req)
	if err != nil {
		return err
	}
	topUp.SnapURL = snapResp.RedirectURL
	return nil
}

func (m midtransService) VerifyPayment(ctx context.Context, orderId string) (bool, error) {
	var client coreapi.Client
	client.New(m.config.Key, m.env)

	transactionStatusResp, err := client.CheckTransaction(orderId)
	if err != nil {
		return false, err
	} else {
		if transactionStatusResp != nil {
			// 5. Do set transaction status based on response from check transaction status
			if transactionStatusResp.TransactionStatus == "capture" {
				if transactionStatusResp.FraudStatus == "challenge" {
					// TODO set transaction status on your database to 'challenge'
					// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
				} else if transactionStatusResp.FraudStatus == "accept" {
					return true, nil
				}
			} else if transactionStatusResp.TransactionStatus == "settlement" {
				return true, nil
			} else if transactionStatusResp.TransactionStatus == "deny" {
				// TODO you can ignore 'deny', because most of the time it allows payment retries
				// and later can become success
			} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
				// TODO set transaction status on your databaase to 'failure'
			} else if transactionStatusResp.TransactionStatus == "pending" {
				// TODO set transaction status on your databaase to 'pending' / waiting payment
			}
		}
	}
	return false, nil
}
