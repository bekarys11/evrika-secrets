	CREATE TABLE IF NOT EXISTS users_secrets (
        id SERIAL UNIQUE,
	    user_id INT,
	    secret_id INT,
	    PRIMARY KEY (user_id, secret_id),
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT fk_role FOREIGN KEY (secret_id) REFERENCES secrets (id)
    )
