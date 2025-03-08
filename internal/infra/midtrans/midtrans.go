package midtrans

import (
	"ambic/internal/domain/dto"
	"ambic/internal/domain/env"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransIf interface {
	GeneratePaymentLink(req dto.RequestSnap) (string, error)
}

type Midtrans struct {
	Snap snap.Client
}

func New(env *env.Env) *Midtrans {
	var s = snap.Client{}
	s.New(env.MidtransServerKey, midtrans.Sandbox)

	return &Midtrans{
		Snap: s,
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
		CustomField1: req.TransactionID,
	}

	var items []midtrans.ItemDetails
	for _, p := range req.TransactionDetails {
		items = append(items, midtrans.ItemDetails{
			ID:    p.ProductID,
			Price: int64(p.Product.FinalPrice),
			Qty:   int32(p.Qty),
			Name:  p.Product.Name,
		})
	}

	url, err := m.Snap.CreateTransactionUrl(snapReq)
	if err != nil {
		return "", err
	}

	return url, nil
}
