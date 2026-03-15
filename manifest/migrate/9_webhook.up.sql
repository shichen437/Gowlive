CREATE TABLE IF NOT EXISTS push_channel_custom_webhook (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    channel_id int(6) NOT NULL, -- 渠道 ID
    webhook_url text NOT NULL, -- webhook url
    request_method tinyint(1) NOT NULL DEFAULT 0, -- 请求方法：0 GET 1 POST
    request_headers text, -- 请求头 
    request_body text, -- 请求体
    created_at datetime(0) NOT NULL, -- 创建时间
    updated_at datetime(0) DEFAULT NULL -- 更新时间
);