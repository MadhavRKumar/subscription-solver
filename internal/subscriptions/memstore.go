package subscriptions

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type memStore struct {
	conn *pgxpool.Pool
}

type NotFoundError struct {
	message string
}

func (e *NotFoundError) Error() string {
	return e.message
}

func NewMemStore() *memStore {
	log.Println(os.Getenv("POSTGRES_URL"))
	conn, err := pgxpool.New(context.Background(), os.Getenv("POSTGRES_URL"))
	if err != nil {
		os.Exit(1)
	}

	return &memStore{
		conn: conn,
	}
}

func (m *memStore) Add(uuid string, subscription Subscription) (Subscription, error) {
	_, err := m.conn.Exec(context.Background(), "INSERT INTO subscriptions (uuid, name, profile_limit, cost) VALUES ($1, $2, $3, $4)", uuid, subscription.Name, subscription.ProfileLimit, subscription.Cost)

	subscription.UUID = uuid

	return subscription, err
}

func (m *memStore) Get(uuid string) (Subscription, error) {
	var sub Subscription
	err := m.conn.QueryRow(context.Background(), "SELECT uuid, name, profile_limit, cost FROM subscriptions where uuid=$1 and deleted_at is NULL", uuid).Scan(&sub.UUID, &sub.Name, &sub.ProfileLimit, &sub.Cost)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Subscription{}, &NotFoundError{"Subscription not found"}
		}
		return Subscription{}, err
	}

	return sub, nil
}

func (m *memStore) List() (map[string]Subscription, error) {
	subscriptions := make(map[string]Subscription)

	rows, err := m.conn.Query(context.Background(), "SELECT uuid, name, profile_limit, cost FROM subscriptions WHERE deleted_at is NULL")
	if err != nil {
		return subscriptions, err
	}

	defer rows.Close()

	for rows.Next() {
		var sub Subscription
		err := rows.Scan(&sub.UUID, &sub.Name, &sub.ProfileLimit, &sub.Cost)
		if err != nil {
			return subscriptions, err
		}

		subscriptions[sub.UUID] = sub
	}

	return subscriptions, err
}

func (m *memStore) Update(uuid string, subscription Subscription) error {
	_, err := m.conn.Exec(context.Background(), "UPDATE subscriptions SET name=$1, profile_limit=$2, cost=$3 WHERE uuid=$4", subscription.Name, subscription.ProfileLimit, subscription.Cost, uuid)

	if err == pgx.ErrNoRows {
		return &NotFoundError{"Subscription not found"}
	}

	return err
}

func (m *memStore) Remove(uuid string) error {
	if _, err := m.conn.Exec(context.Background(), "UPDATE subscriptions SET deleted_at=now() WHERE uuid=$1", uuid); err != nil {
		if err == pgx.ErrNoRows {
			return &NotFoundError{"Subscription not found"}
		}
		return err
	}

	return nil
}
