package subscriptions

import ( 
    "context"
    "os"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

type memStore struct {
    conn *pgxpool.Pool
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

func (*memStore) Get(uuid string) (Subscription, error) {
    return Subscription{}, nil
}

func (m *memStore) List() (map[string]Subscription, error) {
    var subscriptions = make(map[string]Subscription)


    rows, err := m.conn.Query(context.Background(), "SELECT uuid, name, profile_limit, cost FROM subscriptions")

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

func (*memStore) Update(uuid string, subscription Subscription) error {
    return nil
}

func (*memStore) Remove(uuid string) error {
    return nil
}
