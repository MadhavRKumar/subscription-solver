CREATE TABLE IF NOT EXISTS subscriptions(
    id serial PRIMARY KEY,
    name VARCHAR (50) NOT NULL,
    profile_limit int NOT NULL,
    cost int NOT NULL
);
