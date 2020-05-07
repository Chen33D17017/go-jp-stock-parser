DROP DATABASE IF EXISTS `stock_price`;
CREATE DATABASE `stock_price`;
USE `stock_price`;

SET NAMES utf8;
set character_set_client = utf8mb4;

create table `stock_category` (
	`id` INT AUTO_INCREMENT,
    `category_name` char(50) UNIQUE,
    PRIMARY KEY (`id`)
) ENGINE=INNODB;

INSERT INTO `stock_category`(`category_name`) VALUES("その他");

CREATE TABLE `stock_name` (
    `id` INT UNIQUE NOT NULL,
    `category_id` INT NOT NULL DEFAULT 1,
    `name` CHAR(50),
    PRIMARY KEY (`id`),
    FOREIGN KEY (`category_id`)
        REFERENCES `stock_category` (`id`)
)  ENGINE=INNODB;


create table `stock_data` (
	`stock_id` int NOT NULL,
    `price_at` date,
    `open` double NOT NULL,
    `high` double NOT NULL,
    `low` double NOT NULL,
    `close` double NOT NULL,
    `vol` double NOT NULL,
    CONSTRAINT `contacts_pk` PRIMARY KEY (`stock_id`, `price_at`),
    FOREIGN KEY (`stock_id`) REFERENCES `stock_name`(`id`)
) ENGINE=INNODB;




