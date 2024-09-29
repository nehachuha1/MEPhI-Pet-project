DROP TABLE IF EXISTS public.moderators;

CREATE TABLE IF NOT EXISTS public.moderators
(
    id                 SERIAL PRIMARY KEY,
    moderator_username text               NOT NULL,
    access_group       SMALLINT DEFAULT 1 NOT NULL,
    granted_by         text               NOT NULL,
    CONSTRAINT grant_check FOREIGN KEY (granted_by) REFERENCES public.users (login),
    CONSTRAINT user_check FOREIGN KEY (moderator_username) REFERENCES public.users (login)
);