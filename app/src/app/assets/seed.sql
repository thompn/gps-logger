CREATE TABLE IF NOT EXISTS appdb(timestamp DATETIME, longitude TEXT, latitude TEXT, surface TEXT, tracking TEXT);
INSERT OR IGNORE INTO appdb VALUES('2020-01-01', '52.385379', '5.3413432', 'paved', 'init-value');
