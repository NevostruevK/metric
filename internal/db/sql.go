package db

// CREATE TABLE IF NOT EXISTS metrics (
// CREATE TYPE metric_type AS ENUM ('gauge', 'counter');

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
	$1, $2, $3
)
`

	insertCounterSQL = `
INSERT INTO metrics (
	id, mType, delta
) VALUES (
	$1, $2, $3
)
`

	getMetricSQL = `
SELECT * FROM metrics WHERE id = $1 
`
	updateGaugeSQL = `
UPDATE metrics SET id = $1, mType = 'gauge', value = $2 WHERE id = $1
	`

	updateCounterSQL = `
UPDATE metrics SET id = $1, mType = 'counter', delta = $2 WHERE id = $1
	`
)
