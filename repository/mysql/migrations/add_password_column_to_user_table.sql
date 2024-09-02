-- +migrate Up
ALTER TABLE users add column password varchar(255) not null;

-- +migrate Down
ALTER TABLE users drop column password;