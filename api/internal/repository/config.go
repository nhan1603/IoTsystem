package repository

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/nhan1603/IoTsystem/api/internal/appconfig/db/pg"
)

type Backend string

const (
	BackendPostgres  Backend = "postgres"
	BackendCassandra Backend = "cassandra"
)

type Config struct {
	Backend Backend

	// Postgres
	PGURL string // e.g. postgres://postgres:postgres@localhost:5432/iotsystem-pg?sslmode=disable

	// Cassandra
	CassHosts       []string // ["127.0.0.1"]
	CassKeyspace    string   // "iot"
	CassConsistency string   // "QUORUM","LOCAL_QUORUM","ONE"
	CassTimeout     time.Duration
}

// FromEnv builds Config from environment variables.
func FromEnv() Config {
	return Config{
		Backend:         Backend(strings.ToLower(getenv("DB_BACKEND", "postgres"))),
		PGURL:           getenv("PG_URL", "postgres://postgres:postgres@localhost:5432/iotsystem-pg?sslmode=disable"),
		CassHosts:       split(getenv("CASSANDRA_HOSTS", "127.0.0.1")),
		CassKeyspace:    getenv("CASSANDRA_KEYSPACE", "iotsystem"),
		CassConsistency: getenv("CASSANDRA_CONSISTENCY", "QUORUM"),
		CassTimeout:     10 * time.Second,
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func split(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// NewFromConfig constructs a Registry impl for the chosen backend and
// returns a cleanup function you should defer.
func NewFromConfig(ctx context.Context, cfg Config) (Registry, func(), error) {
	switch cfg.Backend {
	case BackendPostgres:
		conn, err := pg.Connect(cfg.PGURL)
		if err != nil {
			return impl{}, nil, fmt.Errorf("postgres create session: %w", err)
		}
		impl := newPostGres(conn)
		cleanup := func() { conn.Close() }
		return impl, cleanup, nil

	case BackendCassandra:
		cluster := gocql.NewCluster(cfg.CassHosts...)
		cluster.Keyspace = cfg.CassKeyspace
		cluster.Timeout = cfg.CassTimeout
		cluster.Consistency = parseConsistency(cfg.CassConsistency)
		// Token-aware, DC-aware policy can be added here if you have LocalDC
		session, err := cluster.CreateSession()
		if err != nil {
			return impl{}, nil, fmt.Errorf("cassandra create session: %w", err)
		}
		impl := newCassandraImpl(session)
		cleanup := func() { session.Close() }
		return impl, cleanup, nil

	default:
		return impl{}, nil, fmt.Errorf("unknown backend: %s", cfg.Backend)
	}
}

func parseConsistency(s string) gocql.Consistency {
	switch strings.ToUpper(s) {
	case "ALL":
		return gocql.All
	case "ONE":
		return gocql.One
	case "LOCAL_ONE":
		return gocql.LocalOne
	case "LOCAL_QUORUM":
		return gocql.LocalQuorum
	default:
		return gocql.Quorum
	}
}
