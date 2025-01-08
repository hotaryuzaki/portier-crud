CREATE TABLE keys (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT,
    is_active BOOLEAN,
    FOREIGN KEY (created_by) REFERENCES users (id)
);