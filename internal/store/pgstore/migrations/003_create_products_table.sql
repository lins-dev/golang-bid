-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS products(
    id serial primary key,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),
    seller_id serial NOT NULL REFERENCES users(id),
    
    product_name TEXT NOT NULL,
    description TEXT NOT NULL,
    price INTEGER NOT NULL,
    auction_end TIMESTAMPTZ NOT NULL,
    is_sold boolean NOT NULL DEFAULT false,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

---- create above / drop below ----

DROP TABLE IF EXISTS products;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
