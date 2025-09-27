CREATE TABLE IF NOT EXISTS sys_notify(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    title varchar(50) NOT NULL, -- 通知标题
    content text NOT NULL, -- 通知内容
    level varchar(50) NOT NULL, -- 通知级别
    status tinyint(1) NOT NULL DEFAULT 0, -- 阅读状态
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

ALTER TABLE live_history ADD COLUMN is_delete int(1) NOT NULL DEFAULT 0;