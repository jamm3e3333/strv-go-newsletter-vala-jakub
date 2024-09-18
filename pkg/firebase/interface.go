package firebase

import "context"

type Connector interface {
	Create(ctx context.Context, opName, path string, data any) error
	Delete(ctx context.Context, opName, path string) error
	GetForData(ctx context.Context, path, opName string, data any) error
	Update(ctx context.Context, path, opName string, data map[string]any) error
}
