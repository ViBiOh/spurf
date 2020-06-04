ALTER TABLE spurf.enedis_value RENAME TO enedis_value_old;

CREATE TABLE spurf.enedis_value (
  name TEXT NOT NULL,
  ts TIMESTAMP WITH TIME ZONE NOT NULL,
  value REAL NOT NULL
);

INSERT INTO spurf.enedis_value (name, ts, value) SELECT 'home', ts, value FROM spurf.enedis_value_old;
DROP TABLE spurf.enedis_value_old;
