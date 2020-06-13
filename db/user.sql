DROP TABLE IF EXISTS `sea_users`;
CREATE TABLE `sea_users` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL DEFAULT '' COMMENT '登录用户名',
  `password` varchar(100) NOT NULL DEFAULT '' COMMENT '登录密码',
  `name` varchar(50) DEFAULT '' COMMENT '用户真实姓名',
  `number` varchar(20) DEFAULT '' COMMENT '用户工号',
  `email` varchar(64) NOT NULL DEFAULT '' COMMENT '邮箱地址',
  `mobile` varchar(11) NOT NULL DEFAULT '' COMMENT '手机号',
  `avatar` varchar(200) NOT NULL DEFAULT '' COMMENT '头像',
  `created_at` datetime NOT NULL COMMENT '创建记录时间',
  `updated_at` datetime NOT NULL COMMENT '更新记录时间',
  `status` int(2) NOT NULL DEFAULT '1' COMMENT '1启用 2禁用',
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表';
