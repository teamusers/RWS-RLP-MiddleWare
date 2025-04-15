-- Create the `users` table
CREATE TABLE `users` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `external_id` VARCHAR(50) NOT NULL,
  `opted_in` BOOLEAN NOT NULL DEFAULT FALSE,
  `external_id_type` VARCHAR(50),
  `email` VARCHAR(50),
  `dob` DATE,
  `country` VARCHAR(50),
  `first_name` VARCHAR(255),
  `last_name` VARCHAR(255),
  `burn_pin` INT(4), 
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB;

-- Create the `user_phone_numbers` table with a foreign key referencing the `users` table
CREATE TABLE `users_phone_numbers` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `user_id` BIGINT NOT NULL,
  `phone_number` VARCHAR(20),
  `phone_type` VARCHAR(20),
  `preference_flags` VARCHAR(50), 
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_user_id`
    FOREIGN KEY (`user_id`)
    REFERENCES `users`(`id`)
    ON DELETE CASCADE
) ENGINE=InnoDB;

CREATE TABLE `sys_channel` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `app_id` varchar(100) NOT NULL,
  `app_key` varchar(100) NOT NULL,
  `status` char(2) NOT NULL DEFAULT '10',
  `sig_method` varchar(100) NOT NULL DEFAULT 'SHA256',
  `create_time` datetime DEFAULT NULL,
  `update_time` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

