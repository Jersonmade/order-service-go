CREATE USER wb_user WITH PASSWORD 'wb_password';

GRANT ALL PRIVILEGES ON DATABASE wb_test_db TO wb_user;

CREATE TABLE orders (
    order_uid VARCHAR PRIMARY KEY,
    track_number VARCHAR,
    entry VARCHAR,
    locale VARCHAR,
    internal_signature TEXT,
    customer_id VARCHAR,
    delivery_service VARCHAR,
    shardkey VARCHAR,
    sm_id INTEGER,
    date_created TIMESTAMP,
    oof_shard VARCHAR
);

CREATE TABLE deliveries (
    order_uid VARCHAR REFERENCES orders(order_uid) ON DELETE CASCADE,
    name VARCHAR,
    phone VARCHAR,
    zip VARCHAR,
    city VARCHAR,
    address VARCHAR,
    region VARCHAR,
    email VARCHAR
);

CREATE TABLE payments (
    order_uid VARCHAR REFERENCES orders(order_uid) ON DELETE CASCADE,
    transaction VARCHAR,
    request_id VARCHAR,
    currency VARCHAR,
    provider VARCHAR,
    amount INTEGER,
    payment_dt BIGINT,
    bank VARCHAR,
    delivery_cost INTEGER,
    goods_total INTEGER,
    custom_fee INTEGER
);

CREATE TABLE items (
    order_uid VARCHAR REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT,
    track_number VARCHAR,
    price INTEGER,
    rid VARCHAR,
    name VARCHAR,
    sale INTEGER,
    size VARCHAR,
    total_price INTEGER,
    nm_id INTEGER,
    brand VARCHAR,
    status INTEGER
);

