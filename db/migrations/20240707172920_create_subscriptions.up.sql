CREATE TABLE IF NOT EXISTS subscriptions(
    id serial PRIMARY KEY,
    uuid UUID NOT NULL,
    name VARCHAR (50) NOT NULL,
    profile_limit int NOT NULL,
    cost int NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_your_table_modtime
    BEFORE UPDATE ON subscriptions
    FOR EACH ROW EXECUTE FUNCTION update_modified_column();
