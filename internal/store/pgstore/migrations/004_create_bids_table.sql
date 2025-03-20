-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS bids(
    id SERIAL PRIMARY KEY,
    uuid UUID UNIQUE NOT NULL DEFAULT gen_random_uuid(),

    product_id SERIAL NOT NULL REFERENCES products(id),
    bidder_id SERIAL NOT NULL REFERENCES users(id),
    bid_amount INTEGER NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

---- create above / drop below ----

DROP TABLE IF EXISTS bids;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
