CREATE TABLE IF NOT EXISTS orders (
    order_id VARCHAR(50) PRIMARY KEY NOT NULL,
    user_id VARCHAR(50),
    order_status VARCHAR(15),
    order_date DATE
);
