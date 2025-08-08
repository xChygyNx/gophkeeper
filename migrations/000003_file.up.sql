create table file
(
    file_id    serial PRIMARY KEY,
    user_id    int          NOT NULL references users (user_id) on delete cascade,
    name       varchar(100) NOT NULL,
    created_at timestamp    NOT NULL,
    deleted_at timestamp NULL
);

CREATE INDEX idx_file_user_id ON file (user_id);
CREATE INDEX idx_file_deleted_at ON file (deleted_at);