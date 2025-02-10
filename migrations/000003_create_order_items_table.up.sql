CREATE TABLE IF NOT EXISTS order_items (
    order_id VARCHAR(50) NOT NULL,
    product_id VARCHAR(50),
    quantity INT DEFAULT 1,
    cost FLOAT
);
