DROP TABLE IF EXISTS public.products CASCADE;

CREATE TABLE IF NOT EXISTS public.products
(
    id             BIGSERIAL PRIMARY KEY,
    name           text    NOT NULL,
    owner_username text    NOT NULL,
    price          integer NOT NULL,
    description    text    NOT NULL,
    create_date    text    NOT NULL,
    edit_date      text    NOT NULL,
    is_active      boolean NOT NULL,
    views          bigint  NOT NULL,
    photo_urls     text[],
    main_photo     text,
    CONSTRAINT user_check FOREIGN KEY (owner_username) REFERENCES public.users (login)
);