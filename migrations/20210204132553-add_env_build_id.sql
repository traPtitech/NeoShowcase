-- +migrate Up
ALTER TABLE `websites`
    DROP FOREIGN KEY `fk_websites_build_id`,
    DROP COLUMN `build_id`;

ALTER TABLE `environments`
    ADD COLUMN `build_id` VARCHAR(22) COMMENT '稼働中のビルドID',
    ADD CONSTRAINT `fk_environments_build_id` FOREIGN KEY (`build_id`) REFERENCES `build_logs` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;

-- +migrate Down
ALTER TABLE `environments`
    DROP FOREIGN KEY `fk_environments_build_id`,
    DROP COLUMN `build_id`;

ALTER TABLE `websites`
    ADD COLUMN `build_id` VARCHAR(22) COMMENT '稼働中のサイトのビルドID',
    ADD CONSTRAINT `fk_websites_build_id` FOREIGN KEY (`build_id`) REFERENCES `build_logs` (`id`) ON UPDATE RESTRICT ON DELETE RESTRICT;
