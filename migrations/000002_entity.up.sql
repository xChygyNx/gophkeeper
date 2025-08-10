create table entity
(
    entity_id     serial PRIMARY KEY,
    user_id     int       NOT NULL references users (user_id) on delete cascade,
    data        bytea     NOT NULL,
    metadata    jsonb not null,
    created_at  timestamp NOT NULL,
    updated_at  timestamp NOT NULL,
    deleted_at  timestamp NULL
);

CREATE INDEX idx_entity_user_id ON entity (user_id);
CREATE INDEX idx_entity_metadata ON entity (metadata);
CREATE INDEX idx_entity_deleted_at ON entity (deleted_at);