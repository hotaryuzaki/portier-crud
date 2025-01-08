CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username varchar NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    gender BOOLEAN,
    id_number VARCHAR(20),
    user_image TEXT,
    tenant_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT,
    is_active BOOLEAN
    FOREIGN KEY (tenant_id) REFERENCES tenants (id)
);