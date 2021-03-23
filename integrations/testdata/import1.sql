SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";

-- 基本用户信息表
CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  KEY `uid` (`id`),
  `user_name` varchar(128) CHARACTER SET utf8mb4 NOT NULL,
  KEY `user_name` (`user_name`),
  `email` varchar(32) NOT NULL,
  KEY `email` (`email`),
  `pass` varchar(256) NOT NULL,
  `passwd` varchar(16) NOT NULL,
  `uuid` TEXT NULL DEFAULT NULL COMMENT 'uuid',
  `t` int(11) NOT NULL DEFAULT '0',
  `u` bigint(20) NOT NULL,
  `d` bigint(20) NOT NULL,
  `plan` varchar(2) CHARACTER SET utf8mb4 NOT NULL DEFAULT 'A',
  `node_group` INT NOT NULL DEFAULT '0',
  `auto_reset_day` INT NOT NULL DEFAULT '0',
  `auto_reset_bandwidth` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `transfer_enable` BIGINT(20) NOT NULL,
  `port` int(11) NOT NULL,
  `protocol_param` VARCHAR(128) NULL DEFAULT NULL,
  `obfs_param` VARCHAR(128) NULL DEFAULT NULL,
  `switch` tinyint(4) NOT NULL DEFAULT '1',
  `enable` tinyint(4) NOT NULL DEFAULT '1',
  `type` tinyint(4) NOT NULL DEFAULT '1',
  `last_get_gift_time` int(11) NOT NULL DEFAULT '0',
  `last_check_in_time` int(11) NOT NULL DEFAULT '0',
  `last_rest_pass_time` int(11) NOT NULL DEFAULT '0',
  `reg_date` datetime NOT NULL,
  `invite_num` int(8) NOT NULL,
  `money` decimal(12,2) NOT NULL,
  `ref_by` int(11) NOT NULL DEFAULT '0',
  `expire_time` int(11) NOT NULL DEFAULT '0',
  `is_email_verify` tinyint(4) NOT NULL DEFAULT '0',
  `reg_ip` varchar(128) NOT NULL DEFAULT '127.0.0.1',
  `node_speedlimit` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `node_connector` int(11) NOT NULL DEFAULT '0',
  `forbidden_ip` LONGTEXT NULL DEFAULT '',
  `forbidden_port` LONGTEXT NULL DEFAULT '',
  `disconnect_ip` LONGTEXT NULL DEFAULT '',
  `is_hide` INT NOT NULL DEFAULT '0',
  `last_detect_ban_time` datetime DEFAULT '1989-06-04 00:05:00',
  `all_detect_number` int(11) NOT NULL DEFAULT '0',
  `is_multi_user` INT NOT NULL DEFAULT '0',
  `telegram_id` BIGINT NULL,
  `is_admin` int(2) NOT NULL DEFAULT '0',
  `im_type` int(11) DEFAULT '1',
  `im_value` text,
  `last_day_t` bigint(20) NOT NULL DEFAULT '0',
  `mail_notified` int(11) NOT NULL DEFAULT '0',
  `class` int(11) NOT NULL DEFAULT '0',
  `class_expire` datetime NOT NULL DEFAULT '1989-06-04 00:05:00',
  `expire_in` datetime NOT NULL DEFAULT '2099-06-04 00:05:00',
  `theme` text NOT NULL,
  `ga_token` text NOT NULL,
  `ga_enable` int(11) NOT NULL DEFAULT '0',
  `pac` LONGTEXT,
  `remark` text
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户流量信息表
-- TODO: 重写流量信息提取逻辑
CREATE TABLE IF NOT EXISTS `user_traffic_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  `user_id` int(11) NOT NULL,
  `u` BIGINT(20) NOT NULL,
  `d` BIGINT(20) NOT NULL,
  `node_id` int(11) NOT NULL,
  `rate` float NOT NULL,
  `traffic` varchar(32) NOT NULL,
  `log_time` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户订阅 TOKEN 信息表
CREATE TABLE IF NOT EXISTS `user_token` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  `token` varchar(256) NOT NULL,
  `user_id` int(11) NOT NULL,
  `create_time` int(11) NOT NULL,
  `expire_time` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 充值码使用信息表
CREATE TABLE IF NOT EXISTS `charge_code` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  `code` text NOT NULL,
  `type` int(11) NOT NULL,
  `number` DECIMAL(11,2) NOT NULL,
  `isused` int(11) NOT NULL DEFAULT '0',
  `userid` bigint(20) NOT NULL,
  `usedatetime` datetime NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 邀请码使用信息表
CREATE TABLE IF NOT EXISTS `invite_code` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  `code` varchar(128) NOT NULL,
  KEY `user_id` (`user_id`),
  `user_id` int(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT '2016-06-01 00:00:00'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 公告信息表
CREATE TABLE IF NOT EXISTS `announcement` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  `date` datetime NOT NULL,
  `content` LONGTEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `markdown` LONGTEXT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 节点信息表
CREATE TABLE IF NOT EXISTS `node` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  `name` varchar(128) NOT NULL,
  `type` int(3) NOT NULL,
  `online_user` int(11) NOT NULL,
  `mu_only` INT NULL DEFAULT '0',
  `online` BOOLEAN NOT NULL DEFAULT TRUE,
  `server` varchar(128) NOT NULL,
  `method` varchar(64) NOT NULL,
  `info` varchar(128) NOT NULL,
  `status` varchar(128) NOT NULL,
  `node_group` INT NOT NULL DEFAULT '0',
  `sort` int(3) NOT NULL,
  `custom_method` tinyint(1) NOT NULL DEFAULT '0',
  `traffic_rate` float NOT NULL DEFAULT '1',
  `node_class` int(11) NOT NULL DEFAULT '0',
  `node_speedlimit` DECIMAL(12,2) NOT NULL DEFAULT '0.00',
  `node_connector` int(11) NOT NULL DEFAULT '0',
  `node_bandwidth` bigint(20) NOT NULL DEFAULT '0',
  `node_bandwidth_limit` bigint(20) NOT NULL DEFAULT '0',
  `bandwidthlimit_resetday` int(11) NOT NULL DEFAULT '0',
  `node_heartbeat` bigint(20) NOT NULL DEFAULT '0',
  `node_ip` text
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
-- TODO: 修改 VPN 节点的结算说明

-- 商店数据表
CREATE TABLE `shop` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` TEXT NOT NULL,
  `price` DECIMAL(12,2) NOT NULL,
  `content` TEXT NOT NULL,
  `auto_renew` INT NOT NULL,
  `status` INT NOT NULL DEFAULT '1',
  `auto_reset_bandwidth` INT NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 优惠券数据表
CREATE TABLE `coupon` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `code` TEXT NOT NULL,
  `onetime` INT NOT NULL,
  `expire` BIGINT NOT NULL,
  `shop` TEXT NOT NULL,
  `credit` INT NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 购买记录数据表
CREATE TABLE `bought` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `userid` BIGINT NOT NULL,
  `shopid` BIGINT NOT NULL,
  `coupon` TEXT NOT NULL,
  `datetime` BIGINT NOT NULL,
  `renew` BIGINT(11) NOT NULL,
  `price` DECIMAL(12,2) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 工单数据表
CREATE TABLE `ticket` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `title` LONGTEXT NOT NULL,
  `status` INT NOT NULL DEFAULT '1',
  `content` LONGTEXT NOT NULL,
  `rootid` BIGINT NOT NULL,`userid` BIGINT NOT NULL,
  `datetime` BIGINT NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 返利记录数据表
CREATE TABLE `payback` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `total` DECIMAL(12,2) NOT NULL,
  `userid` BIGINT NOT NULL,
  `ref_by` BIGINT NOT NULL,
  `ref_get` DECIMAL(12,2) NOT NULL,
  `datetime` BIGINT NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 审计规则数据表
CREATE TABLE `detect_list` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `name` LONGTEXT NOT NULL,
  `type` INT NOT NULL,
  `text` LONGTEXT NOT NULL,
  `regex` LONGTEXT NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 审计记录数据表
CREATE TABLE `detect_log` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `node_id` INT NOT NULL,
  `list_id` BIGINT NOT NULL,
  `datetime` BIGINT NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 中转规则数据表
CREATE TABLE IF NOT EXISTS `relay` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  `user_id` bigint(20) NOT NULL,
  `source_node_id` bigint(20) NOT NULL,
  `dist_node_id` bigint(20) NOT NULL,
  `dist_ip` text NOT NULL,
  `port` int(11) NOT NULL,
  `priority` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 用户订阅日志
CREATE TABLE IF NOT EXISTS `user_subscribe_log` (
  `id`                 int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name`          varchar(128)     NOT NULL COMMENT '用户名',
  `user_id`            int(11)          NOT NULL COMMENT '用户 ID',
  `email`              varchar(32)      NOT NULL COMMENT '用户邮箱',
  `subscribe_type`     varchar(20)      NOT NULL COMMENT '获取的订阅类型',
  `request_ip`         varchar(128)     NOT NULL COMMENT '请求 IP',
  `request_time`       datetime         NOT NULL COMMENT '请求时间',
  `request_user_agent` text                      COMMENT '请求 UA 信息',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户订阅日志';

-- 审计封禁日志
CREATE TABLE IF NOT EXISTS `detect_ban_log` (
  `id`                int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name`         varchar(128)     NOT NULL COMMENT '用户名',
  `user_id`           int(11)          NOT NULL COMMENT '用户 ID',
  `email`             varchar(32)      NOT NULL COMMENT '用户邮箱',
  `detect_number`     int(11)          NOT NULL COMMENT '本次违规次数',
  `ban_time`          int(11)          NOT NULL COMMENT '本次封禁时长',
  `start_time`        bigint(20)       NOT NULL COMMENT '统计开始时间',
  `end_time`          bigint(20)       NOT NULL COMMENT '统计结束时间',
  `all_detect_number` int(11)          NOT NULL COMMENT '累计违规次数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='审计封禁日志';

-- 管理员操作记录
CREATE TABLE IF NOT EXISTS `gconfig` (
  `id`             int(11) unsigned NOT NULL AUTO_INCREMENT,
  `key`            varchar(128)     NOT NULL COMMENT '配置键名',
  `type`           varchar(32)      NOT NULL COMMENT '值类型',
  `value`          text             NOT NULL COMMENT '配置值',
  `oldvalue`       text             NOT NULL COMMENT '之前的配置值',
  `name`           varchar(128)     NOT NULL COMMENT '配置名称',
  `comment`        text             NOT NULL COMMENT '配置描述',
  `operator_id`    int(11)          NOT NULL COMMENT '操作员 ID',
  `operator_name`  varchar(128)     NOT NULL COMMENT '操作员名称',
  `operator_email` varchar(32)      NOT NULL COMMENT '操作员邮箱',
  `last_update`    bigint(20)       NOT NULL COMMENT '修改时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='网站配置';