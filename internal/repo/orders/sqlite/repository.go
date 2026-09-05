package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"layouts-orders-bot/internal/entity"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

func New(path string) (*Repository, error) {

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to database: %w", err)
	}

	return &Repository{db: db}, nil
}

func (s *Repository) NewOrder(order *entity.Order) error {
	properties, err := marshalProperties(order.Properties)
	if err != nil {
		return fmt.Errorf("can't marshal order properties: %w", err)
	}

	q := `
		INSERT INTO orders (
			user_id,
		    username,
			properties,
			min_cost,
			max_cost
		)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	row := s.db.QueryRow(q, order.UserID, order.Username, properties, order.MinCost, order.MaxCost)

	err = row.Scan(&order.ID)
	if err != nil {
		return fmt.Errorf("can't create order: %w", err)
	}

	return nil
}

func marshalProperties(properties map[*entity.State]string) ([]byte, error) {
	result := make(map[string]string, len(properties))

	for state, value := range properties {
		result[state.ID] = value
	}

	return json.Marshal(result)
}

func (s *Repository) GetOrders(userID int) ([]*entity.Order, error) {
	return s.getOrders(userID)
}

func (s *Repository) GetAllOrders() ([]*entity.Order, error) {
	return s.getOrders(0)
}

func (s *Repository) getOrders(userID int) ([]*entity.Order, error) {
	q := `
		SELECT
		    id,
			user_id,
			username,
			properties,
			min_cost,
			max_cost
		FROM orders
	`
	if userID != 0 {
		q += fmt.Sprintf(" WHERE user_id = %d", userID)
	}

	rows, err := s.db.Query(q)
	if err != nil {
		return nil, fmt.Errorf("can't get orders: %w", err)
	}
	defer rows.Close()

	orders := make([]*entity.Order, 0)

	for rows.Next() {
		var properties []byte
		order := &entity.Order{}
		if err := rows.Scan(&order.ID, &order.UserID, &order.Username, &properties, &order.MinCost, &order.MaxCost); err != nil {
			return nil, fmt.Errorf("can't scan order: %w", err)
		}
		if err := json.Unmarshal(properties, &order.Properties); err != nil {
			return nil, fmt.Errorf("can't unmarshal order properties: %w", err)
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("can't iterate orders: %w", err)
	}

	return orders, nil
}

func (s *Repository) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY,
		user_id TEXT NOT NULL,
		username TEXT NOT NULL,
		properties JSONB NOT NULL,
		min_cost INTEGER NOT NULL,
		max_cost INTEGER NOT NULL
	)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return fmt.Errorf("can't create orders table: %w", err)
	}

	return nil
}

func (s *Repository) Close() error {
	return s.db.Close()
}
