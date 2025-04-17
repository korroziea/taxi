BEGIN;

CREATE TYPE car_type AS ENUM (
	'economy',
	'comfort',
	'business'
);

CREATE TABLE IF NOT EXISTS cars(
	id 		   VARCHAR(50) PRIMARY KEY,
	number 	   TEXT 	   NOT NULL,
	color 	   TEXT        NOT NULL,
	type       car_type    NOT NULL,
	created_at TIMESTAMP   DEFAULT NOW(),
   	updated_at TIMESTAMP   DEFAULT NOW()
);

CREATE TYPE work_status AS ENUM (
	'free',
	'busy',
	'off-shift'
);

CREATE TABLE IF NOT EXISTS drivers(
	id 		   VARCHAR(50) PRIMARY KEY,
	first_name VARCHAR(20) NOT NULL,
	email 	   TEXT        NOT NULL,
	phone 	   VARCHAR(15) NOT NULL,
	password   TEXT 	   NOT NULL,
	rate       SMALLINT    DEFAULT 0,
	status     work_status DEFAULT 'off-shift',
	car_id     VARCHAR(50),
	created_at TIMESTAMP   DEFAULT NOW(),
   	updated_at TIMESTAMP   DEFAULT NOW(),
    
   	FOREIGN KEY (car_id) REFERENCES cars (id)
      	ON UPDATE CASCADE
);

COMMIT;
