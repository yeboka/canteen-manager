CREATE TABLE users
(
    id                 bigserial   not null primary key,
    email              varchar     not null unique,
    username           varchar(32) not null,
    role               varchar     not null,
    encrypted_password varchar     not null
);

CREATE TABLE orders
(
    id          bigserial not null primary key,
    user_id     int       not null,
    createdAt   date      not null,
    totalAmount int       not null,
    foreign key (user_id) references users (id)
);

CREATE TABLE categories
(
    id        SERIAL PRIMARY KEY,
    parent_id INTEGER REFERENCES categories (id),
    name      VARCHAR NOT NULL UNIQUE
);

CREATE TABLE menuitem
(
    id          SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES categories (id),
    name        VARCHAR NOT NULL UNIQUE,
    price       INTEGER NOT NULL,
    description VARCHAR NOT NULL
);


CREATE TABLE orderItem
(
    id           serial not null primary key,
    order_id     int    not null,
    menu_item_id int    not null,
    quantity     int    not null,
    foreign key (menu_item_id) references menuitem (id),
    foreign key (order_id) references orders(id)
);
