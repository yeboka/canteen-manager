CREATE TABLE users(
    id bigserial not null primary key ,
    email varchar not null unique ,
    username varchar(32),
    role varchar,
    encrypted_password varchar not null
);
