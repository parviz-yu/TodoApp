CREATE TABLE users (
  id BIGSERIAL not NULL PRIMARY KEY,
  name VARCHAR not NULL,
  email VARCHAR not NULL UNIQUE,
  encrypted_password VARCHAR not NULL
);

CREATE TABLE tasks (
  id BIGSERIAL not NULL PRIMARY KEY,
  user_id BIGSERIAL NOT NULL,
  title VARCHAR NOT NULL,
  description VARCHAR not NULL,
  done BOOLEAN NOT NULL,
  creation_date TIMESTAMP NOT NULL
);