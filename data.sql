SELECT * FROM public.orders
ORDER BY id ASC LIMIT 100;

SELECT 
	m.id, 
	p.price, 
	m.discount,
	m.stock
FROM menus m
	JOIN products p ON p.id = m.product_id
	WHERE m.id = 1;

CREATE TABLE "vouchers" (
  "id" serial PRIMARY KEY,
  "code" varchar(20) DEFAULT '',
  "name" varchar(255) DEFAULT '',
  "description" text DEFAULT '',
  "discount" float DEFAULT 0,
  "start_date" date,
  "end_date" date,
  "usage_limit" int DEFAULT 0,
  "usage_count" int DEFAULT 0,
  created_at timestamp default now()
  deleted_at timestamp
);
INSERT INTO vouchers(code, name, description, discount, start_date, end_date, usage_limit, usage_count)
VALUES ('ADD10PERCENT','additional discount 10%', 'get additional discount 10% for every order above Rp.100.000',
	   0.1, '2026-02-04', '2026-02-10', 50, 0);

ALTER TABLE vouchers ADD COLUMN deleted_at timestamp;

ALTER TABLE "orders" ADD FOREIGN KEY ("voucher_id") REFERENCES "vouchers" ("id");

-- 
WITH avg_rating AS (
  	SELECT
		AVG(r.rating) AS "rating_product",
  		d.menu_id AS "idmenu"
  		FROM reviews r
  		JOIN dt_order d ON d.id = r.id
  		JOIN menus m ON m.id = d.menu_id
  		GROUP BY d.menu_id
	)
	
	SELECT
		p.id,
    	p.name,
    	string_agg(pi.image, ',') AS "image product",
    	p.price,
		p.description,
    	CAST(m.discount AS FLOAT4),
    	ar."rating_product",
		COUNT(ar."idmenu") AS "count reviews"
  	FROM menus m
  	LEFT JOIN avg_rating ar ON ar."idmenu"= m.id
  	LEFT JOIN products p ON p.id = m.product_id
  	LEFT JOIN product_images pi ON pi.product_id = m.product_id
	WHERE m.id = 4
	GROUP BY p.id, m.id, ar."rating_product";

-- 
SELECT * FROM public.orders
ORDER BY id ASC LIMIT 100;

-- SELECT * FROM public.dt_order
-- ORDER BY id ASC 

SELECT
	o.id,
	TO_CHAR(o.created_at, 'dd/mm/yyyy') AS "date"
	STRING_AGG(p.name || dt.qty, ',')
FROM orders o
JOIN dt_order dt ON dt.order_id = o.id
JOIN menus m ON dt.menu_id = m.id
JOIN products p ON p.id = m.product_id;
