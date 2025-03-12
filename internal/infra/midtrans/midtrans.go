package midtrans

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/env"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"time"
)

type MidtransIf interface {
	GeneratePaymentLink(req dto.RequestSnap) (string, error)
}

type Midtrans struct {
	Snap snap.Client
	env  *env.Env
}

func New(env *env.Env) *Midtrans {
	var s = snap.Client{}
	if env.MidtransEnvironment == "production" {
		s.New(env.MidtransServerKey, midtrans.Production)
	} else {
		s.New(env.MidtransServerKey, midtrans.Sandbox)
	}

	return &Midtrans{
		Snap: s,
		env:  env,
	}
}

func (m *Midtrans) GeneratePaymentLink(req dto.RequestSnap) (string, error) {
	custAddress := &midtrans.CustomerAddress{
		FName:       req.User.Name,
		Phone:       req.User.Phone,
		Address:     req.User.Address,
		CountryCode: "IDN",
	}

	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.Amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    req.User.Name,
			Email:    req.User.Email,
			Phone:    req.User.Phone,
			BillAddr: custAddress,
		},
		Expiry: &snap.ExpiryDetails{
			StartTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
			Unit:      "minutes",
			Duration:  m.env.MidtransMaxPaymentDuration,
		},
		CustomField1: req.TransactionID,
	}

	var items []midtrans.ItemDetails
	for _, p := range req.TransactionDetails {
		items = append(items, midtrans.ItemDetails{
			ID:           p.ProductID,
			Price:        int64(p.Product.FinalPrice),
			Qty:          int32(p.Qty),
			Name:         p.Product.Name,
			MerchantName: p.MerchantName,
		})
	}

	snapReq.Items = &items

	url, err := m.Snap.CreateTransactionUrl(snapReq)
	if err != nil {
		return "", err
	}

	return url, nil
}
