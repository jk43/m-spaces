/*
 Navicat Premium Data Transfer

 Source Server         : K3S Mysql - DEV
 Source Server Type    : MySQL
 Source Server Version : 80031 (8.0.31)
 Source Host           : 192.168.1.55:3306
 Source Schema         : dev_auth

 Target Server Type    : MySQL
 Target Server Version : 80031 (8.0.31)
 File Encoding         : 65001

 Date: 14/12/2023 08:01:21
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for credentials
-- ----------------------------
DROP TABLE IF EXISTS `credentials`;
CREATE TABLE `credentials`  (
  `user_id` bigint UNSIGNED NULL DEFAULT NULL,
  `organization_id` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL,
  `first_name` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL,
  `last_name` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL,
  `password` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL,
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_credentials_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `fk_users_credentials`(`user_id` ASC) USING BTREE,
  CONSTRAINT `fk_users_credentials` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci;

-- ----------------------------
-- Records of credentials
-- ----------------------------
BEGIN;
INSERT INTO `credentials` (`user_id`, `organization_id`, `first_name`, `last_name`, `password`, `id`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, '64540b79cb7c89d64e01ff45', 'hyojun', 'kim', '$2a$10$W6LKVCpRihBVVGhW9taguO9gTbYNs/712d/BT6frcoIot6I8ppBFy', 1, '2023-09-20 19:58:53.925', '2023-09-20 19:58:53.925', NULL), (1, '64cc081d96abe03a3b060ea9', 'alex', 'kim', '$2a$10$56Vzzh3xg/JDdQ0dzbv/H.IHnt6lOLx/yYbqzfNIF8ZYog.ayuuMi', 2, '2023-09-20 19:59:21.209', '2023-09-20 19:59:21.209', NULL), (2, '64540b79cb7c89d64e01ff45', 'gogo', 'mama', '$2a$10$NwmBx49hgBJQI1s783Wvj.rkDmn7IotTRrH7opyMTZXGMV7VKpofW', 3, '2023-11-01 13:14:50.768', '2023-11-01 13:14:50.768', NULL), (3, '64540b79cb7c89d64e01ff45', 'james', 'kim', '$2a$10$7r/QDmKkMDJqeOsvVQp/puNgPcAwZzjQqzsX.2eSeUTkjYCWiZgNS', 4, '2023-12-10 20:29:06.339', '2023-12-10 20:29:06.339', NULL);
COMMIT;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `email` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL,
  `object_id` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL,
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_users_deleted_at`(`deleted_at` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 4 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci;

-- ----------------------------
-- Records of users
-- ----------------------------
BEGIN;
INSERT INTO `users` (`email`, `object_id`, `id`, `created_at`, `updated_at`, `deleted_at`) VALUES ('jk@jktech.net', '650b4efdf04d408b0538309f', 1, '2023-09-20 19:58:53.915', '2023-09-20 19:58:53.915', NULL), ('mama@mama.com', '65424f4a1687ab5d3810c74a', 2, '2023-11-01 13:14:50.749', '2023-11-01 13:14:50.749', NULL), ('jame@jktech.net', '65761f92d43e101ec019d778', 3, '2023-12-10 20:29:06.329', '2023-12-10 20:29:06.329', NULL);
COMMIT;

-- ----------------------------
-- Table structure for verification_tokens
-- ----------------------------
DROP TABLE IF EXISTS `verification_tokens`;
CREATE TABLE `verification_tokens`  (
  `token` longtext CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL,
  `verified` enum('Y','N') CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NULL DEFAULT 'N',
  `credentials_id` bigint UNSIGNED NULL DEFAULT NULL,
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) NULL DEFAULT NULL,
  `updated_at` datetime(3) NULL DEFAULT NULL,
  `deleted_at` datetime(3) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_verification_tokens_deleted_at`(`deleted_at` ASC) USING BTREE,
  INDEX `fk_credentials_verification_tokens`(`credentials_id` ASC) USING BTREE,
  CONSTRAINT `fk_credentials_verification_tokens` FOREIGN KEY (`credentials_id`) REFERENCES `credentials` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb3 COLLATE = utf8mb3_general_ci;

-- ----------------------------
-- Records of verification_tokens
-- ----------------------------
BEGIN;
INSERT INTO `verification_tokens` (`token`, `verified`, `credentials_id`, `id`, `created_at`, `updated_at`, `deleted_at`) VALUES ('c2a69fae-40fb-4b00-b9af-da3ae6071dc1', 'Y', 1, 1, '2023-09-20 19:58:53.931', '2023-09-20 19:58:53.931', NULL), ('ebfde703-619a-4a97-9a45-c216e8972f9c', 'Y', 2, 2, '2023-09-20 19:59:21.219', '2023-09-20 19:59:21.219', NULL), ('264c7aa6-a07d-4c3e-9470-1e35eea9b85a', 'Y', 3, 3, '2023-11-01 13:14:50.778', '2023-11-01 13:14:50.778', NULL), ('8a5fda45-1a86-4ef6-bd44-b07c584c4180', 'Y', 4, 4, '2023-12-10 20:29:06.348', '2023-12-10 20:29:06.348', NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
