BEGIN;

CREATE TYPE wallet_type AS ENUM (
	'personal',
	'family'
);

CREATE TABLE IF NOT EXISTS wallets(
	id 		   VARCHAR(50)   PRIMARY KEY,
	owner_id   VARCHAR(50),
	type	   wallet_type   DEFAULT 'personal',
	balance    INTEGER		 DEFAULT 0,
	created_at TIMESTAMP     DEFAULT NOW(),
	updated_at TIMESTAMP     DEFAULT NOW(),
	
	FOREIGN KEY (owner_id) REFERENCES users (id)
    	ON UPDATE CASCADE
);

CREATE TYPE role_type AS ENUM (
	'owner',
	'member'
);

CREATE TABLE IF NOT EXISTS user_wallets(
	user_id   VARCHAR(50),
	wallet_id VARCHAR(50),
	role      role_type    DEFAULT 'member',
	
	PRIMARY KEY (user_id, wallet_id),
	
 	FOREIGN KEY (user_id) REFERENCES users (id)
    	ON UPDATE CASCADE,
   	FOREIGN KEY (wallet_id) REFERENCES wallets (id)
   		ON UPDATE CASCADE
);

COMMIT;
