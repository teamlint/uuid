DROP TABLE IF EXISTS `photos`;
CREATE TABLE `photos` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ksuid` varchar(36) NOT NULL DEFAULT '',
  `ksuid_raw` varchar(100) DEFAULT NULL,
  `ksuid_time` datetime DEFAULT NULL,
  `xid` varchar(36) NOT NULL DEFAULT '',
  `xid_raw` varchar(100) DEFAULT NULL,
  `xid_time` datetime DEFAULT NULL,
  `ulid` varchar(36) DEFAULT NULL,
  `ulid_raw` varchar(100) DEFAULT NULL,
  `ulid_time` datetime DEFAULT NULL,
  `type` int(11) NOT NULL DEFAULT '0' COMMENT '照片类型: 0未设置,1:闯关,2对战,3每日竞猜,4群',
  `created_at` datetime(6) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB CHARSET=utf8mb4;
