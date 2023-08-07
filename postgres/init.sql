DROP TABLE IF EXISTS tb_users;
CREATE TABLE tb_users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255),
    UNIQUE (email)
);

DROP TABLE IF EXISTS tb_categories;
CREATE TABLE tb_categories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id INTEGER,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    UNIQUE (name),
    UNIQUE (slug),
    FOREIGN KEY (user_id) REFERENCES tb_users (id)
);

DROP TABLE IF EXISTS tb_products;
CREATE TABLE tb_products (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price FLOAT,
    quantity INTEGER,
    slug VARCHAR(255) NOT NULL,
    image_url VARCHAR(255),
    category_id INTEGER,
    user_id INTEGER,
    UNIQUE (name),
    UNIQUE (slug),
    FOREIGN KEY (category_id) REFERENCES tb_categories (id),
    FOREIGN KEY (user_id) REFERENCES tb_users (id)
);

DROP TABLE IF EXISTS tb_payment_categories;
CREATE TABLE tb_payment_categories (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id INTEGER,
    category_name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    UNIQUE (category_name),
    UNIQUE (slug),
    FOREIGN KEY (user_id) REFERENCES tb_users (id)
);

DROP TABLE IF EXISTS tb_payment_methods;
CREATE TABLE tb_payment_methods (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id INTEGER NOT NULL,
    category_payment_id INTEGER NOT NULL,
    method_name VARCHAR(255) NOT NULL,
    number VARCHAR(255),
    owner_name VARCHAR(255),
    category_name VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES tb_users (id),
    FOREIGN KEY (category_payment_id) REFERENCES tb_payment_categories (id)
);

DROP TABLE IF EXISTS tb_carts;
CREATE TABLE tb_carts (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES tb_users (id)
);

DROP TABLE IF EXISTS tb_cart_items;
CREATE TABLE tb_cart_items (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    cart_id INTEGER,
    product_id INTEGER,
    quantity INTEGER,
    price FLOAT,
    total_price FLOAT,
    FOREIGN KEY (cart_id) REFERENCES tb_carts (id),
    FOREIGN KEY (product_id) REFERENCES tb_products (id)
);

DROP TABLE IF EXISTS tb_orders;
CREATE TABLE tb_orders (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    user_id INTEGER,
    payment_category_id INTEGER,
    payment_method_id INTEGER,
    receiver_name VARCHAR(255),
    proof_of_payment VARCHAR(255),
    total_price FLOAT,
    status VARCHAR(255),
    FOREIGN KEY (user_id) REFERENCES tb_users (id),
    FOREIGN KEY (payment_category_id) REFERENCES tb_payment_categories (id),
    FOREIGN KEY (payment_method_id) REFERENCES tb_payment_methods (id)
);

DROP TABLE IF EXISTS tb_order_items;
CREATE TABLE tb_order_items (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    order_id INTEGER,
    product_id INTEGER,
    product_name VARCHAR(255),
    quantity INTEGER,
    price FLOAT,
    total_price FLOAT,
    FOREIGN KEY (order_id) REFERENCES tb_orders (id),
    FOREIGN KEY (product_id) REFERENCES tb_products (id)
);
