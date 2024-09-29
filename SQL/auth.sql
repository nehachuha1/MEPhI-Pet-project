DROP TABLE IF EXISTS public.auth CASCADE;

CREATE TABLE IF NOT EXISTS public.auth
(
    id BIGSERIAL PRIMARY KEY,
    login text NOT NULL UNIQUE,
    password text NOT NULL
)