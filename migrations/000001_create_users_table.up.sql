CREATE TABLE IF NOT EXISTS products (
    product_id VARCHAR(50) NOT NULL PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    price FLOAT NOT NULL,
    description MEDIUMTEXT NOT NULL,
    product_image VARCHAR(50),
    date_created DATE,
    date_modified DATE
)
