package mocks

import (
	"github.com/midtrans/midtrans-go"
	"github.com/stretchr/testify/mock"
)

type MockSnapClient struct {
	mock.Mock
}

func (m *MockSnapClient) CreateTransaction(req *midtrans.Request) (*midtrans.Response, *midtrans.Error) {
	args := m.Called(req)

	var response *midtrans.Response
	var err *midtrans.Error

	if tmp := args.Get(0); tmp != nil {
		response = tmp.(*midtrans.Response)
	}
	if tmp := args.Get(1); tmp != nil {
		err = tmp.(*midtrans.Error)
	}

	return response, err
}
