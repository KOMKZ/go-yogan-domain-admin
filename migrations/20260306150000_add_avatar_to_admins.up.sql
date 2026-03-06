ALTER TABLE `admins`
    ADD COLUMN `avatar` VARCHAR(500) NOT NULL DEFAULT '' COMMENT '头像URL' AFTER `phone`,
    ADD COLUMN `avatar_storage_id` VARCHAR(255) NOT NULL DEFAULT '' COMMENT '头像存储ID (driver:storage@filename)' AFTER `avatar`;
