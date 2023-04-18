package db

const (
	schemaSQL = `
	DO $$ BEGIN
	    CREATE TYPE metric_type AS ENUM ('gauge', 'counter');
	EXCEPTION
    	WHEN duplicate_object THEN null;
	END $$;
	CREATE TABLE IF NOT EXISTS metrics (
    id text NOT NULL UNIQUE,
	mType metric_type NOT NULL,
    delta bigint,
    value double precision
);
CREATE INDEX IF NOT EXISTS metrics_id ON metrics(id);
`

	insertGaugeSQL = `
INSERT INTO metrics (
	id, mType, value
) VALUES (
	$1, 'gauge', $2
)
`

	insertCounterSQL = `
INSERT INTO metrics (
	id, mType, delta
) VALUES (
	$1, 'counter', $2
)
`

	getMetricSQL = `
SELECT * FROM metrics WHERE id = $1 LIMIT(1) 
`

	getCounterMetricSQL = `
SELECT mType, delta FROM metrics WHERE id = $1 LIMIT(1)
`

	updateGaugeSQL = `
UPDATE metrics SET mType = 'gauge', value = $2 WHERE id = $1
	`

	updateCounterSQL = `
UPDATE metrics SET mType = 'counter', delta = $2 WHERE id = $1
	`
)
