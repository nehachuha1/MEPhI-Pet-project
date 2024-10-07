package database

// Order scripts

const (
	acceptOrderScript = `UPDATE public.orders SET order_status='IN_PROGRESS' WHERE id=$1;`
	createOrderScript = `INSERT INTO public.orders (seller_username, buyer_username, buyer_name, product_id, product_count,
                           order_comment, order_address, order_status)
						VALUES ($1, $2, $3, $4, $5, $6, $7, 'ACTIVE');`
	completeOrderScript = `UPDATE public.orders SET order_status='COMPLETED' WHERE id=$1;`
	getUserOrdersScript = `SELECT id,
							   seller_username,
							   buyer_username,
							   buyer_name,
							   product_id,
							   product_count,
							   order_comment,
							   order_address,
							   order_status,
							   is_completed
						FROM public.orders
						WHERE buyer_username = $1;`
	getOrderScript = `SELECT id,
							   seller_username,
							   buyer_username,
							   buyer_name,
							   product_id,
							   product_count,
							   order_comment,
							   order_address,
							   order_status,
							   is_completed
						FROM public.orders
						WHERE id = $1;`
	getSellerOrdersScript = `SELECT id,
							   seller_username,
							   buyer_username,
							   buyer_name,
							   product_id,
							   product_count,
							   order_comment,
							   order_address,
							   order_status,
							   is_completed
						FROM public.orders
						WHERE seller_username = $1;`
	checkUserBlockScript = `SELECT id,
								intruder_username,
								moderator_username,
								ban_reason,
								ban_date,
								expires_at
							FROM public.seller_bans
							WHERE intruder_username=$1;`
	blockUserScript = `INSERT INTO public.seller_bans 
    					(intruder_username, moderator_username, ban_reason, ban_date, expires_at)
						VALUES ($1, $2, $3, $5);`
	unblockUserScript = `DELETE FROM public.seller_bans WHERE intruder_username = $1;`
)
