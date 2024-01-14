DO $$ BEGIN
CREATE TYPE types AS ENUM ('auth', 'env', 'ssh');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

    CREATE TABLE IF NOT EXISTS secrets (
    id SERIAL UNIQUE,
    title VARCHAR (100),
	key VARCHAR (100) UNIQUE NOT NULL,
    data TEXT,
    stype types,
	author_id INT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_author
	FOREIGN KEY (author_id)
    REFERENCES users(id) ON DELETE NO ACTION ON UPDATE CASCADE
);

COMMIT;

