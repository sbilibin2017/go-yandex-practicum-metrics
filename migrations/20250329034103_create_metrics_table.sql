
-- +migrate Up
CREATE TABLE IF NOT EXISTS metrics (
    id VARCHAR NOT NULL,
    type VARCHAR NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION,
    PRIMARY KEY (id, type)
);

-- +migrate Down
DROP TABLE IF EXISTS metrics;
