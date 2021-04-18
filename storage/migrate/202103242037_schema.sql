-- +goose Up

create table users(
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    id serial primary key,
    first_name varchar(30) not null,
    last_name varchar(30) not null,
    email varchar(30) not null,
    password varchar(60) not null,
    username varchar(10) not null
);

-- +goose Down

drop table users;