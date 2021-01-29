CREATE DATABASE cab;

USE cab;

CREATE TABLE `users` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `full_name` varchar(255) NOT NULL
);

CREATE TABLE `drivers` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `cab_id` int,
  `driving_licence` varchar(255) UNIQUE NOT NULL,
  `full_name` varchar(255) NOT NULL
);

CREATE TABLE `cabs` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `license_plate` varchar(255) NOT NULL,
  `latitude` int,
  `longitude` int,
  `availability` boolean DEFAULT (false)
);

CREATE TABLE `bookings` (
  `id` int PRIMARY KEY AUTO_INCREMENT,
  `user_id` int,
  `driver_id` int,
  `cab_id` int,
  `source_latitude` int,
  `source_longitude` int,
  `destination_latitude` int,
  `destination_longitude` int,
  `time` datetime DEFAULT (now())
);

ALTER TABLE `drivers` ADD FOREIGN KEY (`cab_id`) REFERENCES `cabs` (`id`);

ALTER TABLE `bookings` ADD FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

ALTER TABLE `bookings` ADD FOREIGN KEY (`driver_id`) REFERENCES `drivers` (`id`);

ALTER TABLE `bookings` ADD FOREIGN KEY (`cab_id`) REFERENCES `cabs` (`id`);
