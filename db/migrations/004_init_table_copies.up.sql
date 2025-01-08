CREATE TABLE IF NOT EXISTS copies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    key_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT NULL,  -- Optional: Allows NULLs if the creator is not specified
    is_active BOOLEAN DEFAULT TRUE,
    FOREIGN KEY (created_by) REFERENCES users (id),
    FOREIGN KEY (key_id) REFERENCES keys (id)
);