CREATE TABLE IF NOT EXISTS todos (
    id SERIAL primary key,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN 
NEW.updated_at=CURRENT_TIMESTAMPS;
RETURN NEW;

END;
$$ language "plpgsql" ;


CREATE TRIGGER update_updated_at_column
    BEFORE Update ON todos
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();