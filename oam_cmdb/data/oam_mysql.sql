CREATE DATABASE IF NOT EXISTS oam default character set utf8 collate utf8_general_ci;

use oam;

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for account
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
  `account_id` int(11) NOT NULL AUTO_INCREMENT,
  `account_name` varchar(150) NOT NULL,
  `type_id` int(11) NOT NULL DEFAULT '0',
  `create_by` varchar(30) NOT NULL,
  `update_by` varchar(30) DEFAULT NULL,
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `field_user` varchar(150) DEFAULT '' COMMENT '账号',
  `field_pwd` varchar(250) DEFAULT NULL COMMENT '账号密码',
  `field_url` varchar(255) DEFAULT NULL COMMENT '账号地址',
  `field_remark` varchar(1000) CHARACTER SET utf8mb4 DEFAULT NULL COMMENT '账号备注',
  `field_other` text CHARACTER SET utf8mb4 COMMENT '账号动态字段内容(JSON)',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `host_id` int(11) DEFAULT NULL,
  `rel_account_ids` varchar(300) DEFAULT NULL,
  PRIMARY KEY (`account_id`),
  KEY `idx_type_id` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for account_field
-- ----------------------------
DROP TABLE IF EXISTS `account_field`;
CREATE TABLE `account_field` (
  `field_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `field_name` varchar(100) NOT NULL,
  `field_key` varchar(50) NOT NULL,
  `type_id` int(11) NOT NULL,
  `is_required` tinyint(1) NOT NULL DEFAULT '0',
  `is_ciphertext` tinyint(1) NOT NULL DEFAULT '0',
  `max_len` int(10) unsigned NOT NULL,
  `value_type` smallint(6) NOT NULL DEFAULT '0' COMMENT '值类型:0=文本,1=布尔,2=数字,3=可选值',
  `value_rule` varchar(300) DEFAULT NULL COMMENT '取值规则json',
  `sort` smallint NOT NULL DEFAULT 0 COMMENT '排序',
  PRIMARY KEY (`field_id`),
  KEY `idx_fieldtype_id` (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for account_type
-- ----------------------------
DROP TABLE IF EXISTS `account_type`;
CREATE TABLE `account_type` (
  `type_id` int(11) NOT NULL AUTO_INCREMENT,
  `type_name` varchar(50) NOT NULL,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for app_info
-- ----------------------------
DROP TABLE IF EXISTS `app_info`;
CREATE TABLE `app_info` (
  `app_id` int(11) NOT NULL AUTO_INCREMENT,
  `proj_id` int(11) NOT NULL DEFAULT '0',
  `app_name` varchar(50) NOT NULL,
  `app_url` varchar(255) DEFAULT NULL COMMENT '外部访问地址',
  `app_port` int(9) DEFAULT NULL COMMENT '占用端口',
  `app_dir` varchar(255) DEFAULT NULL COMMENT '安装目录',
  `sourcecode_repo` varchar(255) DEFAULT NULL COMMENT '源码地址',
  `desc` varchar(600) DEFAULT NULL,
  `create_time` datetime NOT NULL,
  `dev_lang` varchar(30) DEFAULT NULL COMMENT '开发语言',
  `app_type` varchar(30) NOT NULL DEFAULT '' COMMENT '业务应用, web服务器, 关系型数据库, 非关系数据库, 缓存, 消息队列, 文件存储服务, 搜索引擎, DevOps, 其它',
  PRIMARY KEY (`app_id`),
  KEY `idx_proj_id` (`proj_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for dict_item
-- ----------------------------
DROP TABLE IF EXISTS `dict_item`;
CREATE TABLE `dict_item` (
  `item_id` varchar(50) NOT NULL,
  `item_name` varchar(100) NOT NULL,
  `item_value` varchar(5000) DEFAULT NULL,
  `item_type` int NOT NULL DEFAULT '0',
  `update_time` datetime NOT NULL,
  `update_by` varchar(30) NOT NULL,
  PRIMARY KEY (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for document
-- ----------------------------
DROP TABLE IF EXISTS `document`;
CREATE TABLE `document` (
  `doc_id` int(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL,
  `doc_type` varchar(5) NOT NULL COMMENT 'pdf,md...',
  `proj_id` int(11) NOT NULL DEFAULT '0',
  `author_id` int(30) NOT NULL,
  `create_by` varchar(50) NOT NULL,
  `update_by` varchar(30) NOT NULL,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `content` text,
  PRIMARY KEY (`doc_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for function
-- ----------------------------
DROP TABLE IF EXISTS `sys_fun`;
CREATE TABLE `sys_fun` (
  `fun_id` int(11) NOT NULL AUTO_INCREMENT,
  `fun_name` varchar(30) NOT NULL COMMENT '功能名',
  `fun_code` varchar(100) NOT NULL DEFAULT '' COMMENT '功能权限标识',
  `parent_id` int(11) NOT NULL DEFAULT '0',
  `fun_type` tinyint(4) NOT NULL COMMENT '1=菜单 2=功能 3=数据权限',
  `fun_order` smallint(6) unsigned NOT NULL DEFAULT '1',
  `fun_level` tinyint(6) unsigned NOT NULL DEFAULT '1' COMMENT '树层级',
  `create_time` datetime NOT NULL COMMENT '建立时间',
  `create_by` varchar(40) NOT NULL COMMENT '建立人',
  `update_time` datetime NOT NULL COMMENT '更新时间',
  `update_by` varchar(40) NOT NULL COMMENT '更新人',
  `menu_class` varchar(50) DEFAULT NULL COMMENT '菜单css类',
  `fun_url` varchar(1000) DEFAULT NULL COMMENT '功能代表的url,多个以逗号分隔',
  PRIMARY KEY (`fun_id`),
  UNIQUE KEY `uq_index_funcode` (`fun_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='系统功能表';

-- ----------------------------
-- Table structure for host_info
-- ----------------------------
DROP TABLE IF EXISTS `host_info`;
CREATE TABLE `host_info` (
  `host_id` int(11) NOT NULL AUTO_INCREMENT,
  `host_name` varchar(50) NOT NULL,
  `public_ip` varchar(32) DEFAULT NULL COMMENT '公网ip',
  `internal_ip` varchar(32) DEFAULT NULL COMMENT '内网ip',
  `ssh_port` int(11) DEFAULT NULL COMMENT '远程连接端口',
  `os_name` varchar(30) DEFAULT NULL COMMENT '操作系统',
  `create_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0',
  `host_type` tinyint(4) NOT NULL DEFAULT '0' COMMENT '服务器类型:0=通用服务器,1=数据库,2=Web应用,3=文件服务器',
  `desc` varchar(500) DEFAULT NULL COMMENT '说明备注',
  PRIMARY KEY (`host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for project
-- ----------------------------
DROP TABLE IF EXISTS `project`;
CREATE TABLE `project` (
  `proj_id` int(11) NOT NULL AUTO_INCREMENT,
  `proj_name` varchar(50) NOT NULL,
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `creator_id` int(11) NOT NULL COMMENT '创建人ID',
  `proj_desc` varchar(600) DEFAULT NULL,
  PRIMARY KEY (`proj_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rel_host_app
-- ----------------------------
DROP TABLE IF EXISTS `rel_host_app`;
CREATE TABLE `rel_host_app` (
  `host_id` int(11) NOT NULL,
  `app_id` int(11) NOT NULL,
  PRIMARY KEY (`host_id`,`app_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rel_proj_account
-- ----------------------------
DROP TABLE IF EXISTS `rel_proj_account`;
CREATE TABLE `rel_proj_account` (
  `proj_id` int(11) NOT NULL,
  `account_id` int(11) NOT NULL,
  PRIMARY KEY (`proj_id`,`account_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rel_proj_user
-- ----------------------------
DROP TABLE IF EXISTS `rel_proj_user`;
CREATE TABLE `rel_proj_user` (
  `proj_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`proj_id`,`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

DROP TABLE IF EXISTS `rel_proj_host`;
CREATE TABLE `rel_proj_host` (
  `proj_id` int(11) NOT NULL,
  `host_id` int(11) NOT NULL,
  PRIMARY KEY (`proj_id`,`host_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- ----------------------------
-- Table structure for rel_role_function
-- ----------------------------
DROP TABLE IF EXISTS `rel_role_function`;
CREATE TABLE `rel_role_function` (
  `role_code` varchar(30) NOT NULL,
  `fun_id` int(11) NOT NULL,
  PRIMARY KEY (`role_code`,`fun_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='角色功能关联表';

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role` (
  `role_code` varchar(30) NOT NULL COMMENT '角色CODE',
  `role_name` varchar(30) NOT NULL COMMENT '角色名称',
  `role_status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '0=禁用1=启用',
  update_time datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  PRIMARY KEY (`role_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='系统角色表';

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(30) NOT NULL,
  `passwd` varchar(250) NOT NULL,
  `real_name` varchar(50) CHARACTER SET utf8mb4 DEFAULT NULL,
  `user_status` tinyint(4) NOT NULL DEFAULT '1',
  `create_time` datetime NOT NULL,
  `role_code` varchar(30) DEFAULT NULL,
  `create_by` varchar(50) NOT NULL,
  `update_by` varchar(30) NOT NULL,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`user_id`),
  KEY `idx_user_name` (`user_name`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- 默认密码2022@00
INSERT INTO `user_info` (`user_id`, `user_name`, `passwd`, `real_name`, `user_status`, `create_time`, `role_code`, `create_by`, `update_by`, `update_time`) VALUES (1, 'root', 'a3595a8039559093a5108ec841e7fe4794f1301337e79109e88ed93ea3929a2b', '超级管理员', 1, '2022-02-11 14:46:06', 'ROLE_SUPERVISOR', 'root', 'root', '2024-10-10 14:39:14');

INSERT INTO `sys_role` (`role_code`, `role_name`, `role_status`,update_time) VALUES ('ROLE_SUPERVISOR', '超级管理员', '1','2023-09-06 16:06:09');
INSERT INTO `sys_role` (`role_code`, `role_name`, `role_status`,update_time) VALUES ('ROLE_OPS', '运维', '1','2023-09-06 16:06:09');

INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (1, '资源管理', 'RESOUCE_MGT', 0, 1, 1, 1, '2024-07-12 14:36:08', 'root', '2024-07-12 14:36:15', 'root', 'fa fa-tasks', NULL);
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (2, '系统管理', 'SYS_MGT', 0, 1, 2, 1, '2024-07-12 14:37:20', 'root', '2024-07-12 14:37:27', 'root', 'fa fa-cog', NULL);
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (3, '项目', 'PROJ_MGT', 1, 1, 11, 2, '2024-07-12 14:38:37', 'root', '2024-07-12 14:38:39', 'root', NULL, '/project/list');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (4, '新增或修改项目', 'PROJ_MGT:SAVE', 3, 2, 12, 3, '2024-07-12 15:19:59', 'root', '2024-07-12 15:20:04', 'root', NULL, '/project/save');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (5, '删除项目', 'PROJ_MGT:DEL', 3, 2, 15, 3, '2024-07-12 16:07:38', 'root', '2024-08-28 14:32:16', 'root', '', '/project/delproject');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (6, '项目成员管理', 'PROJ_MGT:MEMBER', 3, 2, 13, 3, '2024-07-12 16:07:38', 'root', '2024-08-28 14:55:26', 'root', '', '/project/delmember,/project/addmember,/project/selectablemembers');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (8, '主机', 'HOST_MGT', 1, 1, 12, 2, '2024-08-28 11:08:41', 'root', '2024-08-28 11:08:41', 'root', '', '/host/tohostlist,/host/hostpage,/host/basehostpage');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (9, '应用', 'APP_MGT', 1, 1, 13, 2, '2024-08-28 11:45:22', 'root', '2024-08-28 11:45:22', 'root', '', '/appinfo/toapplist');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (10, '账号', 'ACCOUNT_MGT', 1, 1, 16, 2, '2024-08-28 11:48:26', 'root', '2024-08-28 14:30:55', 'root', '', '/account/toaccountlist,/account/accountpage');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (11, '文档', 'DOC_MGT', 1, 1, 14, 2, '2024-08-28 11:51:34', '', '2024-09-14 11:37:53', 'root', '', '/doc/todoclist,/doc/docpage');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (12, '用户管理', 'USER_MGT', 2, 1, 21, 2, '2024-08-28 11:55:26', '', '2024-09-14 11:27:45', 'root', '', '/user/touserlist');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (13, '角色管理', 'ROLE_MGT', 2, 1, 20, 2, '2024-08-28 12:07:14', 'root', '2024-08-28 14:56:00', 'root', '', '/role/toRole.html,/role/detail,/role/rolelist');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (14, '权限管理', 'FUN_MGT', 2, 1, 22, 2, '2024-08-28 12:31:28', 'root', '2024-08-28 12:31:28', 'root', '', '/fun/toFun.html,/fun/funtree,/fun/detail');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (15, '新增', 'ROLE_MGT:SAVE', 13, 2, 133, 3, '2024-09-06 15:52:58', 'root', '2024-09-06 15:53:05', 'root', NULL, '/role/save');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (16, '删除角色', 'ROLE_MGT:DEL', 13, 2, 134, 3, '2024-09-06 15:56:02', 'root', '2024-09-06 15:56:07', 'root', NULL, '/role/deleterole');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (17, '分配权限', 'ROLE_MGT:ASSIGN_PERM', 13, 2, 135, 3, '2024-09-06 16:15:37', 'root', '2024-09-06 16:15:41', 'root', NULL, '/role/saverolefun,/role/rolefuntree');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (18, '查看', 'PROJ_MGT:VIEW', 3, 2, 14, 3, '2024-09-09 13:41:51', 'root', '2024-09-09 13:41:55', 'root', NULL, '/project/detail,/project/findhostandapptree');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (19, '新增或修改主机', 'HOST_MGT:SAVE', 8, 2, 121, 3, '2024-09-09 13:49:28', 'root', '2024-09-09 13:49:34', 'root', NULL, '/host/save,/host/recoverhost');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (20, '删除主机', 'HOST_MGT:DEL', 8, 2, 122, 3, '2024-09-09 13:49:28', 'root', '2024-09-09 13:49:34', 'root', NULL, '/host/delhost,/host/recoverhost');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (21, '账号类型', 'ACCOUNT_MGT:TYPE_MGT', 10, 2, 161, 3, '2024-09-09 13:53:35', 'root', '2024-09-09 13:53:40', 'root', NULL, '/account/toAccount-type.html,/account/typelist,/account/gettype');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (23, '查看', 'HOST_MGT:VIEW', 8, 2, 123, 3, '2024-09-09 14:04:22', 'root', '2024-09-09 14:04:25', 'root', NULL, '/host/hostdetail');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (24, '新增/编辑', 'ACCOUNT_MGT:TYPE_MGT:EDIT', 21, 2, 124, 3, '2024-09-09 14:28:07', 'root', '2024-09-09 14:28:12', 'root', NULL, '/account/savetype,/account/gettypefields');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (25, '删除', 'ACCOUNT_MGT:TYPE_MGT:DEL', 21, 2, 125, 3, '2024-09-09 14:28:07', 'root', '2024-09-09 14:28:12', 'root', NULL, '/account/deltype');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (26, '系统参数', 'SYS_DICT_MGT', 2, 1, 25, 2, '2024-09-09 14:38:02', 'root', '2024-09-09 14:38:10', 'root', NULL, '/dict/toDict.html,/dict/dictlist,/dict/savedict');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (27, '新增账号', 'ACCOUNT_MGT:ADD', 10, 2, 160, 3, '2024-09-09 14:46:57', 'root', '2024-09-09 14:47:00', 'root', NULL, '/account/toaddaccount,/account/saveaccount');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (28, '删除/作废账号', 'ACCOUNT_MGT:DEL', 10, 2, 161, 3, '2024-09-09 14:46:57', 'root', '2024-09-09 14:47:00', 'root', NULL, '/account/delaccount');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (29, '查看账号信息', 'ACCOUNT_MGT:VIEW', 10, 2, 162, 3, '2024-09-09 15:38:21', 'root', '2024-09-09 15:38:26', 'root', NULL, '/account/toviewaccount,/account/getaccountdetail');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (30, '查看或复制密码', 'ACCOUNT_MGT:COPYPWD', 10, 2, 163, 3, '2024-09-09 15:38:21', 'root', '2024-09-09 15:38:26', 'root', NULL, '');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (32, '编辑账号', 'ACCOUNT_MGT:EDIT', 10, 2, 160, 3, '2024-09-09 14:46:57', 'root', '2024-09-09 14:47:00', 'root', NULL, '/account/toeditaccount,/account/saveaccount,/account/recoveraccount,,/account/getaccountdetail');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (33, '新增/编辑', 'FUN_MGT:SAVE', 14, 2, 220, 3, '2024-09-13 17:27:06', 'root', '2024-09-13 17:27:11', 'root', NULL, '/fun/save');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (34, '删除', 'FUN_MGT:DEL', 14, 2, 221, 3, '2024-09-13 17:27:06', 'root', '2024-09-13 17:27:11', 'root', NULL, '/fun/deletefun');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (35, '新增/修改', 'USER_MGT:SAVE', 12, 2, 210, 3, '2024-09-14 11:25:40', '', '2024-09-14 11:27:56', 'root', '', '/usr/save,/user/detail');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (36, '重置密码', 'USER_MGT:RESETPWD', 12, 2, 212, 3, '2024-09-14 11:26:55', 'root', '2024-09-14 11:26:55', 'root', '', '/user/resetuserpasswd');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (37, '新增/导入', 'DOC_MGT:ADD', 11, 2, 140, 3, '2024-09-14 11:32:22', '', '2024-09-14 11:34:09', 'root', '', '/doc/toadddoc,/doc/savedoc,/doc/importdoc');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (38, '编辑', 'DOC_MGT:EDIT', 11, 2, 141, 3, '2024-09-14 11:33:27', 'root', '2024-09-14 11:33:27', 'root', '', '/doc/toeditdoc,/doc/savedoc');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (39, '查看文档', 'DOC_MGT:VIEW', 11, 2, 142, 3, '2024-09-14 11:35:15', 'root', '2024-09-14 11:35:15', 'root', '', '/doc/view,/doc/preview');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (40, '删除文档', 'DOC_MGT:DEL', 11, 2, 143, 3, '2024-09-14 11:36:01', 'root', '2024-09-14 11:36:01', 'root', '', '/doc/deletedoc');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (41, '新增/修改', 'APP_MGT:SAVE', 9, 2, 130, 3, '2024-09-14 12:13:12', '', '2024-09-14 13:41:43', 'root', '', '/appinfo/appdetail,/appinfo/saveapp,/host/basehostpage');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (43, '删除应用', 'APP_MGT:DEL', 9, 2, 132, 3, '2024-09-14 12:15:50', '', '2024-09-14 13:40:03', 'root', '', '/appinfo/delapp');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (44, '应用列表', 'APP_MGT:LIST', 9, 2, 133, 3, '2024-09-14 12:19:23', '', '2024-09-14 13:42:37', 'root', '', '/appinfo/apppage');

INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 1);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 3);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 4);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 6);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 8);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 9);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 10);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 11);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 12);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 18);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 19);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 20);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 23);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 27);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 28);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 29);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 30);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 32);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 35);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 37);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 38);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 39);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 40);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 41);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 43);
INSERT INTO `rel_role_function` (`role_code`, `fun_id`) VALUES ('ROLE_OPS', 44);


INSERT INTO `account_type` (`type_id`, `type_name`, `update_time`, `create_time`) VALUES (1, '微信公众号', '2023-09-06 16:06:09', '2022-04-13 10:44:43');
INSERT INTO `account_type` (`type_id`, `type_name`, `update_time`, `create_time`) VALUES (3, '微信商户账号', '2023-09-06 16:08:33', '2022-06-24 10:08:52');
INSERT INTO `account_type` (`type_id`, `type_name`, `update_time`, `create_time`) VALUES (4, '支付宝商户账号', '2022-06-24 10:43:02', '2022-06-24 10:43:02');
INSERT INTO `account_type` (`type_id`, `type_name`, `update_time`, `create_time`) VALUES (5, '支付宝应用账号', '2023-09-06 16:28:14', '2022-06-24 10:47:25');
INSERT INTO `account_type` (`type_id`, `type_name`, `update_time`, `create_time`) VALUES (6, '微信小程序', '2023-09-06 16:28:46', '2022-06-24 10:53:41');
INSERT INTO `account_type` (`type_id`, `type_name`, `update_time`, `create_time`) VALUES (7, '数据库账号', '2023-09-06 11:54:11', '2022-06-24 10:53:41');
INSERT INTO `account_type` (`type_id`, `type_name`, `update_time`, `create_time`) VALUES (8, '操作系统账号', '2022-06-24 10:53:41', '2022-06-24 10:53:41');

INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (12, '商户号', 'wxpay_mch_id', 2, 1, 0, 32, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (13, '应用appid', 'wxpay_app_id', 2, 1, 0, 32, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (14, '子应用appid', 'wxpay_sub_appid', 2, 0, 0, 32, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (15, '密钥app_secret', 'wxpay_app_secret', 2, 0, 0, 250, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (16, '密钥apikey', 'wxpay_api_key', 2, 0, 0, 250, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (28, '商户号PID', 'ali_pid', 4, 0, 0, 20, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (29, '操作密码', 'ali_optpasswd', 4, 0, 1, 30, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (30, '支付密码', 'ali_paypasswd', 4, 0, 1, 250, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (31, '绑定人', 'ali_binduser', 4, 0, 0, 50, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (32, '公司全称', 'ali_compay', 4, 0, 0, 250, 0, NULL, 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (66, 'appid', 'wx_app_id', 1, 1, 0, 50, 0, '', 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (67, '密钥appSecret', 'wx_app_secret', 1, 0, 0, 500, 0, '', 1);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (68, '账号绑定人', 'wx_bind_user', 1, 0, 0, 250, 0, '', 3);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (69, '主体信息', 'wx_company', 1, 0, 0, 50, 0, '', 4);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (70, 'API秘钥', 'mch_apikey', 3, 0, 0, 50, 0, '', 3);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (71, '操作密码', 'mch_optpasswd', 3, 0, 1, 250, 0, '', 5);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (72, 'API证书', 'mch_apicert', 3, 0, 0, 250, 4, '', 2);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (73, '绑定人', 'mch_binduser', 3, 0, 0, 250, 0, '', 6);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (74, '公司全称', 'mch_compay', 3, 0, 0, 250, 0, '', 1);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (75, '商户类型', 'mch_usertype', 3, 0, 0, 0, 3, '[{\"text\":\"服务商\",\"id\":1},{\"text\":\"普通商户\",\"id\":0},{\"text\":\"子商户\",\"id\":2}]', 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (88, '应用名', 'alipay_appname', 5, 0, 0, 100, 0, '', 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (89, '应用appid', 'alipay_appid', 5, 1, 0, 250, 0, '', 1);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (90, '支付宝公钥', 'alipay_publickey', 5, 0, 0, 4000, 0, '', 2);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (91, '私钥', 'alipay_private_key', 5, 0, 0, 4000, 0, '', 3);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (92, '接口密钥(AES)', 'alipay_aeskey', 5, 0, 0, 30, 0, '', 5);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (93, 'appid', 'wx_miniappid', 6, 1, 0, 32, 0, '', 0);
INSERT INTO `account_field` (`field_id`, `field_name`, `field_key`, `type_id`, `is_required`, `is_ciphertext`, `max_len`, `value_type`, `value_rule`, `sort`) VALUES (94, '密钥appsecret', 'wx_miniapp_secret', 6, 0, 0, 32, 0, '', 0);

INSERT INTO `dict_item` (`item_id`, `item_name`, `item_value`, `item_type`, `update_time`, `update_by`) VALUES ('service_softwares', '服务软件列表', '[\r\n{\"appName\":\"Mysql\",\"appPort\":3306,\"appType\":\"关系型数据库\",\"devLang\":\"C++\"},\r\n{\"appName\":\"SqlServer\",\"appPort\":1433,\"appType\":\"关系型数据库\"},\r\n{\"appName\":\"Oracle\",\"appPort\":1521,\"appType\":\"关系型数据库\"},\r\n{\"appName\":\"Postgresql\",\"appPort\":5432,\"appType\":\"关系型数据库\"},\r\n{\"appName\":\"MongoDB\",\"appPort\":27017,\"appType\":\"非关系型数据库\"},\r\n{\"appName\":\"HBase\",\"appPort\":16020,\"appType\":\"非关系型数据库\"},\r\n{\"appName\":\"Cassandra\",\"appPort\":7199,\"appType\":\"非关系型数据库\"},\r\n{\"appName\":\"ClickHouse\",\"appPort\":9000,\"appType\":\"非关系型数据库\"},\r\n{\"appName\":\"Nginx\",\"appPort\":80,\"appType\":\"web服务器\"},\r\n{\"appName\":\"Tomcat\",\"appPort\":8080,\"appType\":\"web服务器\"},\r\n{\"appName\":\"IIS\",\"appPort\":80,\"appType\":\"web服务器\"},\r\n{\"appName\":\"Redis\",\"appPort\":6379,\"appType\":\"缓存\"},\r\n{\"appName\":\"Memcached\",\"appPort\":11211,\"appType\":\"缓存\"},\r\n{\"appName\":\"FastDFS\",\"appPort\":22122,\"appType\":\"文件存储服务\"},\r\n{\"appName\":\"MinIO\",\"appPort\":9000,\"appType\":\"文件存储服务\"},\r\n{\"appName\":\"HDFS\",\"appPort\":50070 ,\"appType\":\"文件存储服务\"},\r\n{\"appName\":\"RabbitMQ\",\"appPort\":5672 ,\"appType\":\"消息队列\"},\r\n{\"appName\":\"Kafka\",\"appPort\":9092 ,\"appType\":\"消息队列\"},\r\n{\"appName\":\"RocketMQ\",\"appPort\":9876  ,\"appType\":\"消息队列\"},\r\n{\"appName\":\"ActiveMQ\",\"appPort\":61616 ,\"appType\":\"消息队列\"},\r\n{\"appName\":\"Apache Pulsar\",\"appPort\":6650 ,\"appType\":\"消息队列\"},\r\n{\"appName\":\"ElasticSearch\",\"appPort\":9200 ,\"appType\":\"搜索引擎\"},\r\n{\"appName\":\"Solr\",\"appPort\":8983 ,\"appType\":\"搜索引擎\"},\r\n{\"appName\":\"Jenkins\",\"appPort\":9090 ,\"appType\":\"DevOps\"},\r\n{\"appName\":\"Docker\",\"appType\":\"DevOps\"},\r\n{\"appName\":\"K8S\",\"appType\":\"DevOps\"},\r\n{\"appName\":\"Nacos\",\"appPort\":8848 ,\"appType\":\"其它\"},\r\n{\"appName\":\"Zookeeper\",\"appPort\":2181 ,\"appType\":\"其它\"}\r\n]', '4', '2022-05-27 15:57:52', 'sys');
INSERT INTO `dict_item` (`item_id`, `item_name`, `item_value`, `item_type`, `update_time`, `update_by`) VALUES ('app_type', '常用应用类型', '业务应用\r,web服务器\r,关系型数据库\r,NoSQL\r,缓存\r,消息队列\r,文件存储服务\r,搜索引擎\r,DevOps\r', 3, '2023-11-20 14:19:56', 'root');
INSERT INTO `dict_item` (`item_id`, `item_name`, `item_value`, `item_type`, `update_time`, `update_by`) VALUES ('dev_lang', '常用编程语言', 'Java,C,C++,C#,JavaScript,Go,PHP,Python,VB.net,Swift,Kotlin,Rust,Scala,Erlang', 3, '2023-11-20 09:51:27', 'sys');
INSERT INTO `dict_item` (`item_id`, `item_name`, `item_value`, `item_type`, `update_time`, `update_by`) VALUES ('os_names', '操作系统名称', 'CentOS,Ubuntu,RedHat,OpenEuler,Debian,SUSE,Linux,Windows server 2012,Windows server 2016,Windows server 2019,Windows server', 3, '2022-05-26 16:41:04', 'sys');
