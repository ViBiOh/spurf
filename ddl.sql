-- clean
DROP TABLE IF EXISTS enedis_value;

DROP INDEX IF EXISTS enedis_value_ts;

-- enedis_value
CREATE TABLE enedis_value (
  ts TIMESTAMP WITH TIME ZONE NOT NULL,
  value REAL NOT NULL
);

CREATE UNIQUE INDEX enedis_value_ts ON enedis_value (ts);
