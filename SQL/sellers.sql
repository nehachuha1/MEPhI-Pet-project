DROP TABLE IF EXISTS sellers CASCADE;

CREATE TABLE IF NOT EXISTS public.sellers
(
    id                    SERIAL PRIMARY KEY,
    seller_username       text                  NOT NULL UNIQUE,
    accepted_by_moderator boolean DEFAULT true NOT NULL,
    moderator_username    text                  NOT NULL,
    is_active             boolean DEFAULT true  NOT NULL,
    balance               BIGINT  DEFAULT 0     NOT NULL,
    transactions          text[] DEFAULT ARRAY[]::text[]
)