CREATE TABLE IF NOT EXISTS notifications
(
    id          SERIAL PRIMARY KEY,
    message     TEXT                NOT NULL,
    channel     TEXT                NOT NULL,
    "to"        TEXT                NOT NULL,
    send_at     TIMESTAMP           NOT NULL,
    status      TEXT                NOT NULL CHECK (status IN ('pending', 'sent', 'canceled')) DEFAULT 'pending',
    retries     INT                 NOT NULL DEFAULT 0,
    created_at  TIMESTAMP           DEFAULT NOW(),
    updated_at  TIMESTAMP           DEFAULT NOW()
);