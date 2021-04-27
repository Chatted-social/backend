-- +goose Up
-- +goose StatementBegin
create sequence channels_id_seq
    maxvalue -1000
    increment by -1
    minvalue -100000000;

alter sequence channels_id_seq owner to admin;

create sequence posts_id_seq
    as integer;

alter sequence posts_id_seq owner to admin;

create sequence users_id_seq
    as integer;

alter sequence users_id_seq owner to admin;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
