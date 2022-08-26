create database bjfu;
use bjfu;

DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS safe_log;
DROP TABLE IF EXISTS safe_job;

CREATE TABLE IF NOT EXISTS `user` (
    `pk` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
    `id` varchar(64) NOT NULL COMMENT '唯一学号',
    `user_name` varchar(64) NOT NULL COMMENT '姓名',
    `salt` varchar(64) NOT NULL COMMENT '盐值',
    `password` varchar(64) NOT NULL COMMENT '密码',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否被删除',
    `klass` varchar(64) DEFAULT 'Unknown' COMMENT '班级',
    `is_man` tinyint(1) NOT NULL DEFAULT '0' COMMENT '性别',
    `end_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '毕业时间',
    PRIMARY KEY (`pk`),
    KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '用户表';

CREATE TABLE IF NOT EXISTS `safe_job` (
                                      `pk` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                      `id` varchar(64) NOT NULL COMMENT '唯一标识',
                                      `user_id` varchar(64) NOT NULL COMMENT '学号外键',
                                      `path` varchar(64) NOT NULL COMMENT '发送的消息被存储的文件的路径',
                                      `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
                                      `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                      `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否被删除',
                                      PRIMARY KEY (`pk`),
                                      UNIQUE KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '报平安任务表';

INSERT INTO `safe_job`(`id`,`user_id`,`path`) VALUES('sj-fvdxh7whs7mx9qx9','191002213','/saysafe/191002213.txt');

CREATE TABLE IF NOT EXISTS `safe_log` (
                                          `pk` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增主键',
                                          `id` varchar(64) NOT NULL COMMENT '唯一标识',
                                          `user_id` varchar(64) NOT NULL COMMENT '学号外键',
                                          `job_id` varchar(64) NOT NULL COMMENT '发送的消息被存储的文件的路径',
                                          `is_delete` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否被删除',
                                          `result` text COMMENT '请求返回结果',
                                          `success` tinyint(1) NOT NULL DEFAULT '0' COMMENT '请求是否成功',
                                          `create_time` DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL COMMENT '创建时间',
                                          `update_time` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
                                          PRIMARY KEY (`pk`),
                                          UNIQUE KEY (`id`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT '报平安日志表';

set @@sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
set @@global.sql_mode = 'ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';
delete from user where id='191002213';

