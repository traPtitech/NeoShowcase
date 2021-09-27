-- +migrate Up
CREATE TABLE `providers` (
  `id` CHAR(26) NOT NULL COMMENT 'プロバイダID',
  `name` VARCHAR(16) NOT NULL COMMENT 'プロバイダ名',
  `secret` BINARY(60) NOT NULL COMMENT 'Webhookシークレット',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB 
  DEFAULT CHARACTER SET = utf8mb4 
    COMMENT = 'プロバイダテーブル';

CREATE TABLE `users` (
  `id` CHAR(26) NOT NULL COMMENT 'ユーザID',
  `name` VARCHAR(255) NOT NULL COMMENT 'ユーザ名',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB 
  DEFAULT CHARACTER SET = utf8mb4 
    COMMENT = 'ユーザテーブル';

CREATE TABLE `owners` (
  `user_id` CHAR(26) NOT NULL COMMENT 'ユーザID',
  `app_id` CHAR(26) NOT NULL COMMENT 'アプリID',
  PRIMARY KEY (`user_id`, `app_id`),
  CONSTRAINT `fk_owners_user_id` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT,
  CONSTRAINT `fk_owners_app_id` FOREIGN KEY (`app_id`) REFERENCES `applications`(`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB 
  DEFAULT CHARACTER SET = utf8mb4 
    COMMENT = 'アプリケーション所有者テーブル';

-- +migrate Down
DROP TABLE `providers`;

DROP TABLE `owners`;

DROP TABLE `users`;