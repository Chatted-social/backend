-- +goose Up
-- +goose StatementBegin
create table channels
(
    id              integer   default nextval('channels_id_seq'::regclass) not null,
    title           varchar(35)                                            not null,
    username        varchar(35),
    owner_id        integer                                                not null,
    blocked_ids     integer[],
    subscribers_ids integer[],
    admins_ids      integer[],
    created_at      timestamp default now()
);

alter table channels
    owner to admin;

create unique index channels_id_uindex
    on channels (id);

create table posts
(
    id         serial       not null
        constraint posts_pkey
            primary key,
    post_id    integer      not null,
    owner_id   integer      not null,
    from_id    integer      not null,
    created_at timestamp default now(),
    updated_at timestamp default now(),
    updated    boolean   default false,
    body       varchar(255) not null
);

alter table posts
    owner to admin;

create table users
(
    id         serial  not null,
    first_name varchar not null,
    last_name  varchar not null,
    email      varchar not null,
    username   varchar,
    password   varchar not null
);

alter table users
    owner to admin;

create unique index users_email_uindex
    on users (email);

create unique index users_username_uindex
    on users (username);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
