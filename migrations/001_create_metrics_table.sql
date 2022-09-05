-- Write your migrate up statements here
CREATE TYPE metric_type AS ENUM('counter', 'gauge');

CREATE TABLE IF NOT EXISTS metrics(
    id serial PRIMARY KEY,
    name VARCHAR (255) UNIQUE NOT NULL,
    type metric_type NOT NULL,
    int_value BIGINT,
    float_value DOUBLE PRECISION
);

---- create above / drop below ----

DROP TABLE metrics;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
