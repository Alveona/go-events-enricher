-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events.logs
(
    client_time DateTime('Europe/Moscow'),
    server_time DateTime('Europe/Moscow'),
    device_id   String,
    device_os   String,
    session     String,
    sequence    Int64,
    event       String,
    param_int   Int64,
    ip          String,
    param_str   String,
) ENGINE = MergeTree
        PRIMARY KEY server_time
        PARTITION BY toYYYYMMDD(server_time)
        ORDER BY server_time
;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events.logs
-- +goose StatementEnd
