CREATE TABLE users (
  id BIGSERIAL not NULL PRIMARY KEY,
  name VARCHAR not NULL,
  email VARCHAR not NULL UNIQUE,
  encrypted_password VARCHAR not NULL
);

CREATE TABLE tasks (
  id BIGSERIAL not NULL PRIMARY KEY,
  user_id BIGSERIAL NOT NULL,
  description VARCHAR not NULL,
  status BOOLEAN NoT NULL,
  creation_date TIMESTAMP NOT NULL,
  due_date TIMESTAMP NOT NULL
);