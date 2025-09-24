package interfaces

import (
	"context"
)

type Module interface {
	Start(ctx context.Context) error
	Close(ctx context.Context)
}
