CREATE TABLE IF NOT EXISTS sys_user(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    username varchar(50) NOT NULL unique, -- 用户名
    password varchar(255) NOT NULL, -- 用户密码
    nickname varchar(50) NOT NULL, -- 昵称
    sex tinyint(1) DEFAULT 1, -- 性别
    status tinyint(1) DEFAULT 1, -- 用户状态
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE TABLE IF NOT EXISTS sys_logs(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    type tinyint(2) DEFAULT 1, -- 日志类型: 1 用户 2 直播 3 推送 4 解析
    content text, -- 日志内容
    status tinyint(1) DEFAULT 1, -- 日志状态: 0 错误 1 成功
    created_at datetime(0) NOT NULL -- Created Time
);

CREATE INDEX idx_logs_type ON sys_logs (
    type
);
CREATE INDEX idx_logs_status ON sys_logs (
    status
);

CREATE TABLE IF NOT EXISTS live_manage(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    room_url varchar(255) NOT NULL unique, -- 房间 url
    interval int(6) NOT NULL DEFAULT 30, -- 轮询间隔
    format varchar(10) NOT NULL DEFAULT 'flv', -- 导出视频格式
    monitor_type int(1) NOT NULL DEFAULT 0, -- 监控类型: 0 停止监控 1 实时监控 2 定时监控
    monitor_start_at varchar(20), -- 监控开始时间: 定时监控必填
    monitor_stop_at varchar(20), -- 监控结束时间: 定时监控必填
    remark varchar(255), -- 房间备注
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE TABLE IF NOT EXISTS live_room_info (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    live_id int(6) NOT NULL, -- 房间 ID
    room_name text NOT NULL, -- 房间名称
    anchor text NOT NULL, -- 主播
    platform varchar(100) NOT NULL, -- 直播平台
    status int(2) NOT NULL DEFAULT 1, -- 状态
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE TABLE IF NOT EXISTS live_cookie(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    platform varchar(100) NOT NULL unique, -- 平台
    cookie text NOT NULL, -- cookie
    remark varchar(255), -- 备注
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE TABLE IF NOT EXISTS live_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    live_id int(8) NOT NULL, -- 直播ID
    anchor text NOT NULL, -- 主播
    started_at datetime NOT NULL, -- 直播开始时间
    ended_at datetime DEFAULT NULL, -- 直播结束时间
    duration float DEFAULT NULL, -- 直播时长
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE TABLE IF NOT EXISTS push_channel (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    name varchar(100) NOT NULL, -- 渠道名称
    type varchar(100) NOT NULL, -- 渠道类型
    status int(1) NOT NULL DEFAULT 1, -- 状态：0 禁用 1 启用
    url text DEFAULT NULL, -- webhook
    remark varchar(255) DEFAULT NULL, -- 备注
    created_at datetime(0) NOT NULL, -- 创建时间
    updated_at datetime(0) DEFAULT NULL -- 更新时间
);

CREATE TABLE IF NOT EXISTS push_channel_email (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    channel_id int(6) NOT NULL, -- 渠道 ID
    sender varchar(255) NOT NULL, -- 发送人
    receiver text NOT NULL, -- 接收人
    server varchar(255) NOT NULL, -- 发送服务器地址
    port int(6) NOT NULL, -- 发送端口
    auth_code varchar(255) NOT NULL, -- 授权码
    created_at datetime(0) NOT NULL, -- 创建时间
    updated_at datetime(0) DEFAULT NULL -- 更新时间
);
