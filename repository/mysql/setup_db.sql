CREATE TABLE users (
                         id int primary key AUTO_INCREMENT,
                         name varchar(255) not null,
                         phone_number varchar(255) not null unique,
                         password varchar(255) not null,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
); 