CREATE DATABASE cobaApp;

USE cobaApp;

CREATE TABLE `cars` (
    id int not null primary key AUTO_INCREMENT,
    name varchar(255) not null ,
    price DECIMAL(20, 3) NOT NULL DEFAULT 0.000,
    release_date timestamp not null default current_timestamp
)engine=InnoDB;

INSERT INTO cars(name, price) VALUES ('Toyota Innova Zenix Q', 614000000);