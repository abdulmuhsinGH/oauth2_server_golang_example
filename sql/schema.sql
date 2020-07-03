-- Schema is work in progress

CREATE TABLE IF NOT EXISTS users (
  id uuid DEFAULT uuid_generate_v4(),
  username varchar(25) UNIQUE NOT NULL, 
	password text NOT NULL, 
	firstname varchar(100) NOT NULL,
	middlename varchar(100), 
	lastname varchar(100) NOT NULL, 
	gender varchar(25) NOT NULL,
	email text NOT NULL,
	PRIMARY KEY (id, username)
	
);

CREATE TABLE IF NOT EXISTS oauth_clients (
  id     TEXT  NOT NULL,
  secret TEXT  NOT NULL,
  domain TEXT  NOT NULL,
  data   JSONB NULL,
  CONSTRAINT oauth_clients_pkey PRIMARY KEY (id)
);