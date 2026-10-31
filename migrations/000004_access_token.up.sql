create table access_token
(
    access_token     varchar(100) NOT NULL,
    user_id          int         NOT NULL references users (user_id) on delete cascade,
    created_at       timestamp   NOT NULL,
    end_date_at      timestamp   NOT NULL
);

CREATE UNIQUE INDEX idx_access_token_unique ON access_token (access_token);
CREATE INDEX idx_access_token_end_date_at ON access_token (end_date_at);