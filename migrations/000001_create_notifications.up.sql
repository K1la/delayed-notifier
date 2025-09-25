CREATE TABLE IF NOT EXISTS notifications
(
    id          SERIAL PRIMARY KEY,
    message     TEXT                NOT NULL,
    send_at     TIMESTAMP           NOT NULL,
    status      TEXT                NOT NULL CHECK (status IN ('pending', 'sent', 'canceled')) DEAFULT 'pending',
    retries     INT                 NOT NULL DEAFULT 0,
    "to"        TEXT                NOT NULL,
    channel     TEXT                NOT NULL,
    created_at  TIMESTAMP                    DEAFULT NOW(),
    updated_at  TIMESTAMP                    DEAFULT NOW(),
);