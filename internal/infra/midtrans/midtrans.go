package midtrans

import (
	"ambic/internal/domain/env"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransIf interface {
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
