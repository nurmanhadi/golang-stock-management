CREATE TYPE move AS ENUM('IN', 'OUT');
CREATE TABLE movements(
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    type move NOT NULL,
    qty INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_movements_products FOREIGN KEY(product_id) REFERENCES products(id)
);