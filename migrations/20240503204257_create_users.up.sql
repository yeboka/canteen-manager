CREATE TABLE users
(
    id                 bigserial   not null primary key,
    email              varchar     not null unique,
    username           varchar(32) not null,
    role               varchar     not null,
    encrypted_password varchar     not null
);

CREATE TABLE orders(
    id bigserial not null primary key,
    user_id int not null,
    createdAt date not null,
    totalAmount int not null ,
    foreign key (user_id) references users(id)
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
    name        VARCHAR NOT NULL UNIQUE ,
    price       INTEGER NOT NULL,
    description VARCHAR NOT NULL
);
