DROP TABLE IF EXISTS orders;

CREATE TABLE IF NOT EXISTS public.orders
(
    id              SERIAL PRIMARY KEY,
    seller_username text                                      NOT NULL,
    buyer_username  text                                      NOT NULL,
    product_id      BIGINT                                    NOT NULL,
    product_count   INT     DEFAULT 1                         NOT NULL,
    order_comment   text    DEFAULT 'Комментарий отсутствует' NOT NULL,
    order_address   text                                      NOT NULL,
    order_status    text    DEFAULT 'ACTIVE'                  NOT NULL,
    is_completed    boolean DEFAULT false                     NOT NULL,
    CONSTRAINT order_st CHECK (order_status = 'ACTIVE' OR order_status = 'IN_PROGRESS' OR
                               order_status = 'COMPLETED'),
    CONSTRAINT prod_id FOREIGN KEY (product_id) REFERENCES public.products (id),
    CONSTRAINT seller_name FOREIGN KEY (seller_username) REFERENCES public.sellers(seller_username)
)