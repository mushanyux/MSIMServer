-- +migrate Up

ALTER TABLE `app` ADD COLUMN app_name VARCHAR(40) NOT NULL DEFAULT '' COMMENT 'app名字'; 

ALTER TABLE `app` ADD COLUMN app_logo VARCHAR(400) NOT NULL DEFAULT '' COMMENT 'app logo'; 