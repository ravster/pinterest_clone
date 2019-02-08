--create databases
CREATE DATABASE pc;

---connect to dev and add uuid extension
\connect pc;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS plpgsql;

insert into users (created_at, updated_at, username, email, token) values(now(), now(), 'ravi', 'foo', 'foo');
