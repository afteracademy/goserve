package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DbConfig struct {
	User        string
	Pwd         string
	Host        string
	Port        uint16
	Name        string
	MinPoolSize uint16
	MaxPoolSize uint16
	Timeout     time.Duration
	SSLMode     string
}

type Database interface {
	GetInstance() *database
	Connect()
	Disconnect()
	Pool() *pgxpool.Pool
}

type database struct {
	pool    *pgxpool.Pool
	context context.Context
	config  DbConfig
}

func NewDatabase(ctx context.Context, config DbConfig) Database {
	return &database{
		context: ctx,
		config:  config,
	}
}

func (db *database) GetInstance() *database {
	return db
}

func (db *database) Pool() *pgxpool.Pool {
	return db.pool
}

func (db *database) Connect() {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.config.User,
		db.config.Pwd,
		db.config.Host,
		db.config.Port,
		db.config.Name,
		coalesce(db.config.SSLMode, "disable"),
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatal("failed to parse postgres config:", err)
	}

	cfg.MinConns = int32(db.config.MinPoolSize)
	cfg.MaxConns = int32(db.config.MaxPoolSize)
	cfg.MaxConnLifetime = time.Hour
	cfg.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(db.context, db.config.Timeout)
	defer cancel()

	fmt.Println("connecting postgres...")
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal("connection to postgres failed:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("pinging postgres failed:", err)
	}

	fmt.Println("connected to postgres!")
	db.pool = pool
}

func (db *database) Disconnect() {
	fmt.Println("disconnecting postgres...")
	if db.pool != nil {
		db.pool.Close()
	}
	fmt.Println("disconnected postgres")
}

func ParseUUID(id string) (uuid.UUID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.New(id + " is not a valid uuid")
	}
	return u, nil
}

func coalesce(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
