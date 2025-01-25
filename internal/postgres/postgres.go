package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(user, password, host string) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s", user, password, host)
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &DB{Pool: pool}, nil
}

func (db *DB) ExecuteFunction(ctx context.Context, query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, field := range rows.FieldDescriptions() {
			row[string(field.Name)] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}
