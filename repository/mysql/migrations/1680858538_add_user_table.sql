-- +migrate Up
-- please read this article to understand why we use VARCHAR(191)
-- https://www.grouparoo.com/blog/varchar-191#why-varchar-and-not-text
CREATE TABLE `users` (
                       `id` INT PRIMARY KEY AUTO_INCREMENT,
                       `name` VARCHAR(191) NOT NULL,
                       `phone_number` VARCHAR(191) NOT NULL UNIQUE,
                       `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +migrate Down
DROP TABLE `users`;