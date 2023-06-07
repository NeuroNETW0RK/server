/*
 Navicat Premium Data Transfer

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 80030 (8.0.30)
 Source Host           : 127.0.0.1:3306
 Source Schema         : neuronet

 Target Server Type    : MySQL
 Target Server Version : 80030 (8.0.30)
 File Encoding         : 65001

 Date: 27/05/2023 14:32:48
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_permission
-- ----------------------------
DROP TABLE IF EXISTS `admin_permission`;
CREATE TABLE `admin_permission` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  `name` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限名字',
  `resource` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限资源',
  `parent_id` bigint DEFAULT NULL COMMENT '父权限id',
  `root_id` bigint DEFAULT NULL COMMENT '根权限id',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_resource_deleted` (`resource`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_permission
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_permission_role
-- ----------------------------
DROP TABLE IF EXISTS `admin_permission_role`;
CREATE TABLE `admin_permission_role` (
  `role_id` bigint DEFAULT NULL,
  `permission_id` bigint DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_permission_role
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_role
-- ----------------------------
DROP TABLE IF EXISTS `admin_role`;
CREATE TABLE `admin_role` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '角色名字',
  `parent_id` bigint DEFAULT NULL COMMENT '父角色id',
  `root_id` bigint DEFAULT NULL COMMENT '根角色id',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_deleted` (`name`,`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_role
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for admin_user
-- ----------------------------
DROP TABLE IF EXISTS `admin_user`;
CREATE TABLE `admin_user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  `account` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账号',
  `name` varchar(64) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '姓名',
  `password` varchar(256) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '密码',
  `deleted_at` bigint unsigned DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_account_deleted` (`account`,`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_user
-- ----------------------------
BEGIN;
INSERT INTO `admin_user` (`id`, `add_time`, `update_time`, `account`, `name`, `password`, `deleted_at`) VALUES (1, '2023-05-27 14:32:01.239', '2023-05-27 14:32:01.239', 'admin', 'admin', '$pbkdf2-sha512$n8GIXL1Q7CS5uQxv$10bef60dcd2aec074d28d4ecd947be1bc5b77d876e5a2b7585d730998e9d6133', 0);
COMMIT;

-- ----------------------------
-- Table structure for admin_user_role
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_role`;
CREATE TABLE `admin_user_role` (
  `user_id` bigint DEFAULT NULL,
  `role_id` bigint DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of admin_user_role
-- ----------------------------
BEGIN;
INSERT INTO `admin_user_role` (`user_id`, `role_id`) VALUES (1, 2);
COMMIT;

-- ----------------------------
-- Table structure for cluster
-- ----------------------------
DROP TABLE IF EXISTS `cluster`;
CREATE TABLE `cluster` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `add_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '集群名称',
  `config_path` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '集群配置文件路径',
  `description` longtext COLLATE utf8mb4_unicode_ci COMMENT '集群描述',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of cluster
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for image_build
-- ----------------------------
DROP TABLE IF EXISTS `image_build`;
CREATE TABLE `image_build` (
                               `id` bigint NOT NULL AUTO_INCREMENT,
                               `image_id` bigint NOT NULL DEFAULT '0',
                               `build_status` tinyint NOT NULL DEFAULT '0',
                               `build_log` longtext COLLATE utf8mb4_general_ci NOT NULL,
                               `build_user` bigint NOT NULL DEFAULT '0',
                               `add_time` bigint NOT NULL,
                               `update_time` bigint NOT NULL,
                               PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for image_info
-- ----------------------------
CREATE TABLE `image_info` (
                              `id` bigint NOT NULL AUTO_INCREMENT,
                              `image_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL,
                              `build_type` tinyint NOT NULL DEFAULT '0' COMMENT '构建方式 1:from self 2docker file 3 harbor',
                              `image_ext` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
                              `image_user_for` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
                              `create_user` bigint NOT NULL DEFAULT '0',
                              `update_user` bigint NOT NULL DEFAULT '0',
                              `create_time` bigint NOT NULL DEFAULT '0',
                              `update_time` bigint NOT NULL DEFAULT '0',
                              PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Table structure for image_tag
-- ----------------------------
CREATE TABLE `image_tag` (
                             `id` bigint NOT NULL AUTO_INCREMENT,
                             `tag_name` varchar(255) COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
                             `tag_desc` varchar(2000) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '',
                             `image_id` bigint NOT NULL DEFAULT '0',
                             `create_user` bigint NOT NULL DEFAULT '0',
                             `add_time` bigint NOT NULL DEFAULT '0',
                             `update_time` bigint NOT NULL DEFAULT '0',
                             PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
