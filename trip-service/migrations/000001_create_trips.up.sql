BEGIN;

CREATE TYPE status_type AS ENUM (
	'processing',
	'waiting',
	'executing',
	'finished',
	'canceled'
);

CREATE TABLE IF NOT EXISTS trips (
	id 		     VARCHAR(50)  PRIMARY KEY,
	status       status_type  DEFAULT 'processing',
	user_id 	 VARCHAR(50)  NOT NULL,
	cost 		 BIGSERIAl    NOT NULL,
	start_point  JSONB        NOT NULL,
	end_point 	 JSONB        NOT NULL,
	distance 	 INTEGER 	  DEFAULT 0,
	duration 	 INTEGER 	  DEFAULT 0,
	driver_id 	 VARCHAR(50),
	driver_name  VARCHAR(20),
	driver_rate  SMALLINT,
	car_id       VARCHAR(50),
	car_number   TEXT,
	car_color    TEXT,
	waiting_time INTEGER      DEFAULT 0,
 	created_at   TIMESTAMP    DEFAULT NOW(),
    updated_at   TIMESTAMP    DEFAULT NOW()
);

COMMIT;
