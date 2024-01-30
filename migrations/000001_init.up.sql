CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    email varchar(256) UNIQUE,
    password_hash bytea NOT NULL,
    name character varying(256),
    date_of_birth date,
    create_date date DEFAULT CURRENT_DATE
);

CREATE TABLE refresh_tokens
(
    id SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    token VARCHAR(256) NOT NULL,
    create_date DATE DEFAULT CURRENT_DATE,
    expired_date DATE,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
CREATE INDEX idx_user_id ON refresh_tokens (user_id);