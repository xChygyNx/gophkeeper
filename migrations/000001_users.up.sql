create table users
(
    user_id    serial PRIMARY KEY,
    username   varchar(255) NOT NULL,
    password   text NOT NULL,
    created_at timestamp   NOT NULL,
    updated_at  timestamp  NULL,
    deleted_at timestamp NULL
);

CREATE UNIQUE INDEX idx_username_unique ON users (username);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);