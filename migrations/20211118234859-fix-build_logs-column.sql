-- +migrate Up
ALTER TABLE `build_logs`
MODIFY COLUMN `result` varchar(10) COLLATE 'utf8mb4_general_ci' NOT NULL COMMENT 'ビルド結果',
ADD CONSTRAINT `fk_build_logs_result` FOREIGN KEY (`result`) REFERENCES `build_status` (`status`);

-- +migrate Down
ALTER TABLE `build_logs`
MODIFY COLUMN `result` ENUM ('BUILDING', 'SUCCEEDED', 'FAILED', 'CANCELED') NOT NULL COMMENT 'ビルド結果',
DROP FOREIGN KEY `fk_build_logs_result`;
