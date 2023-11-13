	CREATE TABLE IF NOT EXISTS roles (
    id SERIAL UNIQUE,
    name VARCHAR (50) UNIQUE NOT NULL,
    alias VARCHAR (50) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    )
