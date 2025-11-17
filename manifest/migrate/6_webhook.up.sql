CREATE TABLE IF NOT EXISTS push_channel_webhook (
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    channel_id int(6) NOT NULL, -- 渠道 ID
    webhook_url text NOT NULL, -- webhook url
    message_type tinyint(1) NOT NULL DEFAULT 0, -- 消息类型：0 文本 1 富文本 2 卡片
    sign text, -- 签名
    at text, -- @用户
    created_at datetime(0) NOT NULL, -- 创建时间
    updated_at datetime(0) DEFAULT NULL -- 更新时间
);
