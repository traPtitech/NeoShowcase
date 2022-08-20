-- +migrate Up
CREATE TABLE `build_status` (
  `status` varchar(10) NOT NULL COMMENT 'ビルドの状態',
  PRIMARY KEY (`status`)
) COMMENT='ビルドの状態' ENGINE='InnoDB' COLLATE 'utf8mb4_general_ci';

INSERT INTO `build_status` (`status`)
VALUES ('BUILDING'), ('SUCCEEDED'), ('FAILED'), ('CANCELED'), ('QUEUED'), ('SKIPPED');

-- +migrate Down
DROP TABLE `build_status`;
