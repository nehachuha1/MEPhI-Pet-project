DROP TABLE IF EXISTS public.seller_bans;

CREATE TABLE IF NOT EXISTS public.seller_bans
(
    id                 SERIAL PRIMARY KEY,
    intruder_username  text               NOT NULL,
    moderator_username text               NOT NULL,
    ban_reason         text               NOT NULL,
    ban_date           date DEFAULT now() NOT NULL,
    expires_at         date               NOT NULL,
    CONSTRAINT intr_username FOREIGN KEY (moderator_username) REFERENCES public.sellers (seller_username),
    CONSTRAINT moder_username FOREIGN KEY (intruder_username) REFERENCES public.users (login)
)