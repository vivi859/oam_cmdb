ALTER TABLE sys_role 
ADD COLUMN update_time datetime not null DEFAULT '2024-01-01 01:01:01';

update sys_role set role_code='ROLE_SUPERVISOR',role_name='超级管理员' where role_code='root';
update user_info set role_code='ROLE_SUPERVISOR' where user_name='root';

DROP TABLE "function";
CREATE TABLE "sys_fun" (
  "fun_id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "fun_name" nvarchar(30) NOT NULL,
  "fun_code" nvarchar(100) NOT NULL,
  "parent_id" integer NOT NULL DEFAULT 0,
  "fun_type" integer NOT NULL,
  "fun_order" integer NOT NULL DEFAULT 1,
  "fun_level" integer NOT NULL DEFAULT 1,
  "menu_class" nvarchar(200),
  "fun_url" nvarchar(1000),
  "create_time" datetime NOT NULL,
  "create_by" nvarchar(40) NOT NULL,
  "update_time" datetime NOT NULL,
  "update_by" nvarchar(40) NOT NULL
);

CREATE UNIQUE INDEX "uq_index_funcode_function"
ON "sys_fun" (
  "fun_code" ASC
);

ALTER TABLE user_info ADD COLUMN create_by1 nvarchar(50);
ALTER TABLE user_info ADD COLUMN update_by1 nvarchar(30);
ALTER TABLE user_info ADD COLUMN update_time1 datetime;

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
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (24, '新增/编辑', 'ACCOUNT_TYPE_MGT:EDIT', 21, 2, 124, 3, '2024-09-09 14:28:07', 'root', '2024-09-09 14:28:12', 'root', NULL, '/account/savetype,/account/gettypefields');
INSERT INTO sys_fun (fun_id, fun_name, fun_code, parent_id, fun_type, fun_order, fun_level, create_time, create_by, update_time, update_by, menu_class, fun_url) VALUES (25, '删除', 'ACCOUNT_TYPE_MGT:DEL', 21, 2, 125, 3, '2024-09-09 14:28:07', 'root', '2024-09-09 14:28:12', 'root', NULL, '/account/deltype');
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

INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 1);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 3);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 4);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 6);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 8);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 9);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 10);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 11);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 12);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 18);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 19);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 20);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 23);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 27);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 28);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 29);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 30);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 32);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 35);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 37);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 38);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 39);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 40);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 41);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 43);
INSERT INTO rel_role_function (role_code, fun_id) VALUES ('ROLE_OPS', 44);