
-- +migrate Up
ALTER TABLE `repositories` 
  ADD COLUMN `owner` varchar(255) NOT NULL COMMENT 'レポジトリのオーナー' AFTER `id`,
  ADD COLUMN `name` varchar(255) NOT NULL  COMMENT 'レポジトリ名' AFTER `owner`,
  ADD COLUMN `provider_id` char(26) NOT NULL AFTER `url`
  ADD CONSTRAINTS `fk_repositories_provider_id` FOREIGN KEY (`provider_id`) REFERENCES `providers` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE `applications`
  DROP COLUMN `owner`,
  DROP COLUMN `name`;

-- +migrate Down

ALTER TABLE `repositories`
  DROP FOREIGN KEY `fk_repositories_provider_id`,
  DROP COLUMN `provider_id`,
  DROP COLUMN `owner`,
  DROP COLUMN `name`;

ALTER TABLE `applications`
  ADD COLUMN `owner` varchar(100) NOT NULL COMMENT 'アプリケーションのオーナー' AFTER `id`,
  ADD COLUMN `name` varchar(100) NOT NULL  COMMENT 'アプリケーション名' AFTER `owner`;
