DROP TABLE IF EXISTS public.users CASCADE;

CREATE TABLE IF NOT EXISTS public.users
(
    id            SERIAL PRIMARY KEY,
    login         TEXT UNIQUE,
    first_name    TEXT,
    second_name   TEXT,
    sex           TEXT,
    age           INTEGER,
    address       TEXT,
    register_date TEXT,
    edit_date     TEXT,
    CONSTRAINT auth_check FOREIGN KEY (login) REFERENCES public.auth(login)
);
