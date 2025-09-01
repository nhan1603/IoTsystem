package cassandra

import (
	"context"

	"github.com/gocql/gocql"
)

type ctxKey struct{}

var batchKey ctxKey

// WithBatch puts *gocql.Batch into context so repos can append queries.
func WithBatch(ctx context.Context, b *gocql.Batch) context.Context {
	return context.WithValue(ctx, batchKey, b)
}

// BatchFrom extracts *gocql.Batch from context.
func BatchFrom(ctx context.Context) (*gocql.Batch, bool) {
	b, ok := ctx.Value(batchKey).(*gocql.Batch)
	return b, ok && b != nil
}

// BatchLike is the minimal surface we need from *gocql.Batch`.
type BatchLike interface {
	Query(stmt string, names []string) *gocql.Batch
	WithContext(ctx interface{}) *gocql.Batch // we won’t use this from repos
}
