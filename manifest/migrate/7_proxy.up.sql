CREATE TABLE IF NOT EXISTS sys_proxy(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    platform varchar(20) NOT NULL unique, -- 平台
    proxy text NOT NULL, -- 代理
    remark varchar(255), -- 备注
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

ALTER TABLE live_manage ADD COLUMN monitor_only int(1) NOT NULL DEFAULT 0;
