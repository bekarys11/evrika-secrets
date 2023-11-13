	CREATE TABLE IF NOT EXISTS users_roles (
	    user_id INT,
	    role_id INT,
	    PRIMARY KEY (user_id, role_id),
        CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id),
        CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles (id)
    )
