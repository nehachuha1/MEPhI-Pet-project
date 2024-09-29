INSERT INTO public.orders (seller_username, buyer_username, product_id, product_count, order_comment, order_address,
                           order_status, is_completed)
VALUES ($1, $2, $3, $4, $5, $6, 'ACTIVE', $7);