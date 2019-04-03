SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS  `Condition`;
CREATE TABLE `Condition` (
  `id` int(32) unsigned zerofill NOT NULL AUTO_INCREMENT COMMENT '唯一ID，主键，自增，索引',
  `typeID` int(32) unsigned NOT NULL COMMENT '外键，ConditionType表中的ID',
  `value` int(32) unsigned NOT NULL COMMENT '条件值',
  `name` text NOT NULL COMMENT '与value对应，用于显示的名称，如：菲律宾老师',
  `probability` float NOT NULL DEFAULT '1' COMMENT '该条件的生效几率最大值为1，最小值为0，精确到小数点后两位。',
  `des` text NOT NULL COMMENT '针对value的描述',
  `remark` text NOT NULL COMMENT '备用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `condition_index` (`id`),
  KEY `conditionSubKey` (`typeID`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='条件表';

DROP TABLE IF EXISTS  `RDAuth`;
CREATE TABLE `RDAuth` (
  `id` smallint(16) unsigned zerofill NOT NULL AUTO_INCREMENT COMMENT '唯一ID，主键，自增，索引',
  `authName` tinytext NOT NULL COMMENT '权限名',
  `authDes` text NOT NULL COMMENT '权限描述',
  `remark` text NOT NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `rdauth_index` (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=7 DEFAULT CHARSET=utf8 COMMENT='RD后台用户权限表';

insert into `RDAuth`(`id`,`authName`,`authDes`,`remark`) values
(1,'account_cd','账号创建、删除权限',''),
(2,'strategy_cd','创建、删除“策略”权限',''),
(3,'account_s','账号查询权限',''),
(4,'strategy_e','编辑“策略”权限',''),
(5,'strategy_s','查询“策略”权限',''),
(6,'account_e','账号“权限编辑”权限','');
DROP TABLE IF EXISTS  `StrategyResultCache_0`;
CREATE TABLE `StrategyResultCache_0` (
  `id` int(32) unsigned zerofill NOT NULL AUTO_INCREMENT COMMENT '唯一ID，主键，自增，索引',
  `sid` int(32) unsigned NOT NULL COMMENT '外键，对应StrategyCategroy中的id',
  `userID` int(32) unsigned zerofill NOT NULL DEFAULT '00000000000000000000000000000000' COMMENT '用户id（未偏移过的）',
  `userRole` int(32) NOT NULL DEFAULT '-1' COMMENT '用户的角色类型ID，对应接口返回的roleStyle',
  `valuePath` text NOT NULL COMMENT '策略文件模板的存放路径，默认为“”空字符串',
  `process` text NOT NULL COMMENT '记录整个匹配过程（如：各条件怎样计算，最终怎样选择策略，k值趋近）',
  `remark` text NOT NULL COMMENT '备用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `sr_0_index` (`id`),
  KEY `sr_0_subkey` (`sid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='策略结果缓存表，用于缓存用户ID已经匹配的策略映射，及记录匹配“策略”时的匹配算法';

DROP TABLE IF EXISTS  `RDUser`;
CREATE TABLE `RDUser` (
  `id` int(32) unsigned zerofill NOT NULL AUTO_INCREMENT COMMENT '唯一ID，主键，自增，索引',
  `signName` varchar(32) NOT NULL DEFAULT '' COMMENT '用户登录名',
  `signPWD` varchar(32) NOT NULL COMMENT '用户登录密码',
  `roleID` smallint(16) unsigned NOT NULL DEFAULT '0' COMMENT 'RD用户角色，外键，对应RDRole表中id',
  `remark` text NOT NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idIndex` (`id`),
  KEY `roleIDSubKey` (`roleID`)
) ENGINE=MyISAM AUTO_INCREMENT=2 DEFAULT CHARSET=utf8 COMMENT='RD后台账号管理表';

insert into `RDUser`(`id`,`signName`,`signPWD`,`roleID`,`remark`) values
(1,'guominglong','543b10.',1,'');
DROP TABLE IF EXISTS  `Strategy`;
CREATE TABLE `Strategy` (
  `id` int(32) unsigned zerofill NOT NULL AUTO_INCREMENT COMMENT '唯一ID，主键，自增，索引',
  `sid` int(32) unsigned NOT NULL COMMENT '外键，对应StrategyCategroy中的id',
  `conditionGroup` text NOT NULL COMMENT '条件组：对应条件表中Condition表中的各id的组合：例如1，2，3',
  `valuePath` text NOT NULL COMMENT '策略文件模板的存放路径，默认为“”空字符串',
  `expireDate` int(32) unsigned NOT NULL COMMENT '过期时间的utc时间戳秒值',
  `enabled` bit(1) NOT NULL DEFAULT b'0' COMMENT '是否启用，0为关闭，1为启用',
  `remark` text NOT NULL COMMENT '备用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `strategy_index` (`id`),
  KEY `Strategy_subkey` (`sid`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='策略表';

DROP TABLE IF EXISTS  `RDRole`;
CREATE TABLE `RDRole` (
  `id` smallint(16) unsigned zerofill NOT NULL AUTO_INCREMENT COMMENT '唯一ID，主键，自增，索引',
  `roleName` tinytext NOT NULL COMMENT '角色名',
  `roleDes` text NOT NULL COMMENT '角色描述',
  `authGroup` text NOT NULL COMMENT '对应权限表(RDAuth)中的一系列权限id拼接后的字符串，例:1,2,3',
  `remark` text NOT NULL COMMENT '备注',
  PRIMARY KEY (`id`),
  UNIQUE KEY `rdrole_index` (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='RD后台用户角色表';

insert into `RDRole`(`id`,`roleName`,`roleDes`,`authGroup`,`remark`) values
(1,'最高管理员','拥有最高权限','1,2,3,4,5,6',''),
(2,'管理员','负责维护后台账号的日常行为','2,3,4,5,6',''),
(3,'策略开发者','负责维护策略模板及条件匹配项','2,4,5',''),
(4,'策略维护','负责维护现有策略和现有条件，使他们任意组合。以及修改策略模板中的各值','4,5','');
DROP TABLE IF EXISTS  `ConditionType`;
CREATE TABLE `ConditionType` (
  `id` int(32) unsigned zerofill NOT NULL AUTO_INCREMENT COMMENT '唯一ID，主键，自增，索引',
  `enName` tinytext NOT NULL COMMENT '英文显示名',
  `zhName` tinytext NOT NULL COMMENT '中文显示名',
  `des` text NOT NULL COMMENT '条件类型描述',
  `remark` text NOT NULL COMMENT '备用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `conditiontype_index` (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='条件类型表';

DROP TABLE IF EXISTS  `StrategyCategroy`;
CREATE TABLE `StrategyCategroy` (
  `id` int(32) unsigned zerofill NOT NULL AUTO_INCREMENT COMMENT '唯一ID，主键，自增，索引',
  `name` tinytext NOT NULL COMMENT '策略名',
  `des` text NOT NULL COMMENT '策略描述',
  `remark` text NOT NULL COMMENT '备用',
  PRIMARY KEY (`id`),
  UNIQUE KEY `sc_index` (`id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COMMENT='策略分类表';

SET FOREIGN_KEY_CHECKS = 1;

