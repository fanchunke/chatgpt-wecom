CREATE DATABASE IF NOT EXISTS chatgpt DEFAULT CHARACTER SET utf8mb4;

CREATE TABLE IF NOT EXISTS `chatgpt.sessions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` varchar(50) NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `updated_at` timestamp NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `deleted_at` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `session_status_user_id` (`status`,`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS `chatgpt.messages` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `from_user_id` varchar(50) NOT NULL,
  `to_user_id` varchar(50) NOT NULL,
  `content` longtext NOT NULL,
  `created_at` timestamp NULL DEFAULT current_timestamp(),
  `spouse_id` bigint(20) DEFAULT NULL,
  `session_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `spouse_id` (`spouse_id`),
  KEY `message_session_id_from_user_id_created_at` (`session_id`,`from_user_id`,`created_at`),
  KEY `message_session_id_to_user_id_created_at` (`session_id`,`to_user_id`,`created_at`),
  CONSTRAINT `messages_messages_spouse` FOREIGN KEY (`spouse_id`) REFERENCES `messages` (`id`) ON DELETE SET NULL,
  CONSTRAINT `messages_sessions_messages` FOREIGN KEY (`session_id`) REFERENCES `sessions` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin