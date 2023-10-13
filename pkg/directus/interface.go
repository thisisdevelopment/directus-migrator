package directus

import "context"

type IDirectus interface {
	GetSchema(context.Context) (string, error)
	Diff(ctx context.Context, diff string) (string, error)
	Apply(ctx context.Context, diff string) (string, error)
}
