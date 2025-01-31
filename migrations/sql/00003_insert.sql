-- +goose up
insert into test (id, name, description) values (gen_random_uuid(), '1', 'test');
