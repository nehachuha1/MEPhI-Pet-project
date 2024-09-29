SELECT id,
       intruder_username,
       moderator_username,
       ban_reason,
       ban_date,
       expires_at
FROM public.seller_bans
WHERE intruder_username=$1;