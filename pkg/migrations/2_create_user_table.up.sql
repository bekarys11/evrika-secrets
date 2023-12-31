	CREATE TABLE IF NOT EXISTS users (
    id SERIAL UNIQUE,
    name VARCHAR (50) UNIQUE NOT NULL,
	email VARCHAR (100) UNIQUE NOT NULL,
    password VARCHAR (250) NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT TRUE,
	role_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_role
    FOREIGN KEY (role_id)
    REFERENCES roles(id) ON DELETE NO ACTION ON UPDATE CASCADE
    )
