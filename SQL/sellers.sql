DROP TABLE IF EXISTS sellers CASCADE;

CREATE TABLE IF NOT EXISTS public.sellers
(
    id                    SERIAL PRIMARY KEY,
    seller_username       text                  NOT NULL UNIQUE,
    accepted_by_moderator boolean DEFAULT false NOT NULL,
    moderator_username    text                  NOT NULL,
    is_active             boolean DEFAULT true  NOT NULL,
    is_banned             boolean DEFAULT false NOT NULL,
    ban_id                SERIAL                NOT NULL UNIQUE,
    balance               BIGINT  DEFAULT 0     NOT NULL,
    transactions          text[],
    CONSTRAINT user_check FOREIGN KEY (seller_username) REFERENCES public.users (login)
)