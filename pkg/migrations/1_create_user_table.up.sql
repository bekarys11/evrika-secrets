	CREATE TABLE IF NOT EXISTS users (
    id SERIAL UNIQUE,
    name VARCHAR (50) UNIQUE NOT NULL,
	email VARCHAR (100) UNIQUE NOT NULL,
    password VARCHAR (250) NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
    )