package producerinterface

import "context"

type ProdInterInterface interface {
	Send(ctx context.Context, key string, value []byte) error
}
