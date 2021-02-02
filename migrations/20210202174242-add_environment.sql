-- +migrate Up
CREATE TABLE `environments`
(
    `id`             VARCHAR(22)              NOT NULL COMMENT '環境ID',
    `application_id` VARCHAR(22)              NOT NULL COMMENT 'アプリケーションID',
    `branch_name`    VARCHAR(100)             NOT NULL COMMENT 'Gitブランチ・タグ名',
    `build_type`     ENUM ('image', 'static') NOT NULL COMMENT 'ビルドタイプ',
    `created_at`     DATETIME(6)              NOT NULL COMMENT '作成日時',
    `updated_at`     DATETIME(6)              NOT NULL COMMENT '更新日時',
    PRIMARY KEY (`id`),
    UNIQUE (`application_id`, `branch_name`),
    CONSTRAINT fk_environments_application_id FOREIGN KEY (`application_id`) REFERENCES applications (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
    COMMENT = 'アプリ環境テーブル';

ALTER TABLE `repositories`
    DROP COLUMN `refs`;

ALTER TABLE `applications`
    DROP COLUMN `build_type`;

ALTER TABLE `websites`
    DROP FOREIGN KEY `fk_websites_application_id`;
ALTER TABLE `websites`
    DROP COLUMN `application_id`;
ALTER TABLE `websites`
    ADD COLUMN `environment_id` VARCHAR(22) NOT NULL COMMENT '環境ID';
ALTER TABLE `websites`
    ADD CONSTRAINT `fk_websites_environment_id` FOREIGN KEY (`environment_id`) REFERENCES `environments` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE `websites`
    ADD UNIQUE `uk_environment_id` (`environment_id`);


ALTER TABLE `build_logs`
    DROP FOREIGN KEY `fk_build_logs_application_id`;
ALTER TABLE `build_logs`
    DROP COLUMN `application_id`;
ALTER TABLE `build_logs`
    ADD COLUMN `environment_id` VARCHAR(22) COMMENT '環境ID';
ALTER TABLE `build_logs`
    ADD CONSTRAINT `fk_build_logs_environment_id` FOREIGN KEY (`environment_id`) REFERENCES `environments` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;

-- +migrate Down

ALTER TABLE `build_logs`
    DROP FOREIGN KEY `fk_build_logs_environment_id`;
ALTER TABLE `build_logs`
    DROP COLUMN `environment_id`;
ALTER TABLE `build_logs`
    ADD COLUMN `application_id` VARCHAR(22) COMMENT 'アプリケーションID';
ALTER TABLE `build_logs`
    ADD CONSTRAINT `fk_build_logs_application_id` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;

ALTER TABLE `websites`
    DROP FOREIGN KEY `fk_websites_environment_id`;
ALTER TABLE `websites`
    DROP COLUMN `environment_id`;
ALTER TABLE `websites`
    ADD COLUMN `application_id` VARCHAR(22) NOT NULL COMMENT 'アプリケーションID';
ALTER TABLE `websites`
    ADD CONSTRAINT `fk_websites_application_id` FOREIGN KEY (`application_id`) REFERENCES `applications` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
ALTER TABLE `websites`
    ADD UNIQUE `uk_application_id` (`application_id`);

ALTER TABLE `applications`
    ADD COLUMN build_type ENUM ('image', 'static') NOT NULL COMMENT 'ビルドタイプ';

ALTER TABLE `repositories`
    ADD COLUMN `refs` TEXT NOT NULL COMMENT '使用するGit Ref';

DROP TABLE `environments`;