
drop table config;
drop table endpoints;
drop table eventlist;
drop table eventlog;
drop table weightdata;
drop table weightproc;

create table config
(
  item     varchar(32)  null,
  value    varchar(16)  null,
  itemtype char         null,
  comment  varchar(128) null
);

INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('id_colliery', '140000000', 'S', '本煤矿的编码');
INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('interval_heartbeat_app', '600', 'N', '应用程序心跳间隔(秒)');
INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('timeout_endpoint_data', '600', 'N', '端点数据超时时长(秒)');
INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('duration_endpoint_data', '10', 'N', '端点数据持续多少时间才有效(秒)');
INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('min_endpoint_data', '1.5', 'N', '端点数据大于多少吨才有效(吨)');
INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('num_endpoint', '1', 'N', '端点的数量');
INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('interval_reconnect', '60', 'N', '短线重连间隔(秒)');
INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('event_epdata_change', 'true', 'T', '是否生成端点数据改变事件');
INSERT INTO bridge.config (item, value, itemtype, comment) VALUES ('save_epdata_zero', 'true', 'T', '是否存储值为零的数据');

create table endpoints
(
  ipaddr  varchar(15) null,
  tcpport int         null,
  unit    tinyint     null
  comment '0,单位为吨；1，单位为公斤',
  protnb  tinyint     null
  comment '协议编号，只能是1，2，3',
  epid    char(3)     null
);

INSERT INTO bridge.endpoints (ipaddr, tcpport, unit, protnb, epid) VALUES ('192.168.4.232', 6000, 1, 1, '001');
INSERT INTO bridge.endpoints (ipaddr, tcpport, unit, protnb, epid) VALUES ('192.168.4.232', 6001, 1, 2, '002');
INSERT INTO bridge.endpoints (ipaddr, tcpport, unit, protnb, epid) VALUES ('192.168.4.232', 6002, 1, 3, '003');


create table eventlist
(
  eventid   char(2)      null,
  eventtext varchar(255) null
);

INSERT INTO bridge.eventlist (eventid, eventtext) VALUES ('A0', '应用程序启动');
INSERT INTO bridge.eventlist (eventid, eventtext) VALUES ('A1', '应用程序心跳');
INSERT INTO bridge.eventlist (eventid, eventtext) VALUES ('E0', '端点连接成功');
INSERT INTO bridge.eventlist (eventid, eventtext) VALUES ('E1', '端点连接中');
INSERT INTO bridge.eventlist (eventid, eventtext) VALUES ('E2', '端点连接断开');
INSERT INTO bridge.eventlist (eventid, eventtext) VALUES ('E4', '端点数据超时');
INSERT INTO bridge.eventlist (eventid, eventtext) VALUES ('E5', '端点数据变动');
INSERT INTO bridge.eventlist (eventid, eventtext) VALUES ('E3', '端点连接心跳(不使用)');

create table eventlog
(
  recid      int auto_increment
    primary key,
  collieryid char(9)      null,
  bridgeid   char(3)      null,
  eventtime  timestamp    null,
  eventid    char(2)      null,
  eventinfo  varchar(255) null
);

create table weightdata
(
  RecId       varchar(32)   not null
    primary key,
  CollieryId  char(9)       not null,
  BridgeId    char(3)       not null,
  VehiNum     varchar(10)   not null,
  BeginTime   datetime      not null,
  ValTime     datetime      not null,
  EndTime     datetime      not null,
  WeightValue decimal(8, 2) not null
);

create table weightproc
(
  RecID       varchar(32)   not null
    primary key,
  WeightTime  datetime      null,
  WeightValue decimal(8, 2) null,
  VehNum      varchar(10)   null,
  CollieryID  char(9)       null,
  BridgeID    char(3)       null
);

create table users
(
  name   varchar(16) null,
  salt   varchar(16) null,
  passwd varchar(32) null
);

insert into users (name,salt,passwd) values ("admin","FvscE", md5(concat("FvscE","13910580009")));

