--create databases
CREATE DATABASE pc;
CREATE DATABASE pc_test;

---connect to dev and add uuid extension
\connect pc;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS plpgsql;

---connect to test and add uuid extension
\connect pc_test;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS plpgsql;
