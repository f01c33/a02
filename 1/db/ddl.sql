CREATE DATABASE tb01
    WITH
    OWNER = postgres
    ENCODING = 'UTF8'
    LC_COLLATE = 'en_US.UTF-8'
    LC_CTYPE = 'en_US.UTF-8'
    LOCALE_PROVIDER = 'libc'
    TABLESPACE = pg_default
    CONNECTION LIMIT = -1
    IS_TEMPLATE = False;

CREATE TABLE IF NOT EXISTS tb01 (
    id SERIAL PRIMARY KEY NOT NULL,
    col_texto text NOT NULL,
    col_dt timestamp NOT NULL
);