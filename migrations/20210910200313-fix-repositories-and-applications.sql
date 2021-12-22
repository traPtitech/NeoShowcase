-- +migrate Up
ALTER TABLE `repositories` 
  CHANGE COLUMN `remote` `url` TEXT NOT NULL COMMENT 'Git remote URL',
  ADD COLUMN `owner` varchar(256) NOT NULL COMMENT 'リポジトリのオーナー' AFTER `id`,
  ADD COLUMN `name` varchar(256) NOT NULL  COMMENT 'リポジトリ名' AFTER `owner`,
  ADD COLUMN `provider_id` char(36) NOT NULL COMMENT 'プロバイダID' AFTER `url`,
  ADD CONSTRAINT `fk_repositories_provider_id` FOREIGN KEY (`provider_id`) REFERENCES `providers` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE `applications`
  DROP COLUMN `owner`,
  DROP COLUMN `name`;

-- +migrate Down

ALTER TABLE `applications`
  ADD COLUMN `owner` varchar(100) NOT NULL COMMENT 'アプリケーションのオーナー' AFTER `id`,
  ADD COLUMN `name` varchar(100) NOT NULL  COMMENT 'アプリケーション名' AFTER `owner`;

ALTER TABLE `repositories`
  DROP FOREIGN KEY `fk_repositories_provider_id`,
  DROP COLUMN `provider_id`,
  DROP COLUMN `owner`,
  DROP COLUMN `name`;

