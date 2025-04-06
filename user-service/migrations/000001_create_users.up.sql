CREATE TABLE IF NOT EXISTS users (
	id 		   VARCHAR(50) PRIMARY KEY,
	first_name VARCHAR(20) NOT NULL,
	email 	   TEXT        NOT NULL,
	phone 	   VARCHAR(15) NOT NULL,
	password   TEXT 	   NOT NULL,
 	created_at TIMESTAMP   DEFAULT NOW(),
    updated_at TIMESTAMP   DEFAULT NOW()
);
