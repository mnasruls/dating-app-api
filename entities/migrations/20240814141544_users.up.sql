begin;

create extension if not exists "uuid-ossp";

CREATE TABLE IF NOT EXISTS users
(
    id            uuid            NOT NULL default uuid_generate_v4() primary key,
    username      varchar(50)     NOT NULL,
    phone_number  varchar(50)     NOT NULL,
    password      text            NOT NULL,
    verified      boolean         NOT NULL DEFAULT 'false',
    created_at    timestamp       NOT NULL,
    updated_at    timestamp       NULL,
    deleted_at    timestamptz     NULL
);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);


commit;