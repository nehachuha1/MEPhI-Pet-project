INSERT INTO public.orders (seller_username, buyer_username, buyer_name, product_id, product_count, order_comment, order_address,
                           order_status)
VALUES ($1, $2, $3, $4, $5, $6, $7,'ACTIVE');