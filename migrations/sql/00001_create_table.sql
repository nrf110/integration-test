-- +goose up
CREATE TABLE test(
    id uuid not null,
    name text not null
);