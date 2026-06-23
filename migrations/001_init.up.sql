-- Create products and orders tables
CREATE TABLE IF NOT EXISTS products (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    price       DECIMAL(10,2) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id           UUID PRIMARY KEY,
    product_id   INTEGER NOT NULL REFERENCES products(id),
    user_id      TEXT NOT NULL,
    ordered_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Seed a product with initial inventory (managed in Redis)
INSERT INTO products (id, name, price)
VALUES (1, 'Flash Sale Item', 99.99)
ON CONFLICT (id) DO NOTHING;