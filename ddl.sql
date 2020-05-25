-- clean
DROP TABLE IF EXISTS enedis_value;
DROP INDEX IF EXISTS enedis_value_ts;

DROP SCHEMA IF EXISTS spurf;

-- schema
CREATE SCHEMA spurf;

-- enedis_value
CREATE TABLE spurf.enedis_value (
  ts TIMESTAMP WITH TIME ZONE NOT NULL,
  value REAL NOT NULL
);

CREATE UNIQUE INDEX enedis_value_ts ON spurf.enedis_value(ts);
