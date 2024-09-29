SELECT id,
       seller_username,
       buyer_username,
       product_id,
       product_count,
       order_comment,
       order_address,
       order_status,
       is_completed
FROM public.orders
WHERE seller_username = $1;