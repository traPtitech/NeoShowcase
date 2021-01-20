-- +migrate Up
ALTER TABLE `build_logs`
    MODIFY application_id VARCHAR(22) COMMENT 'アプリケーションID';


-- +migrate Down
ALTER TABLE `build_logs`
    MODIFY application_id VARCHAR(22) NOT NULL COMMENT 'アプリケーションID';
