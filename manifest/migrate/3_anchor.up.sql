CREATE TABLE IF NOT EXISTS anchor_info(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    anchor_name text NOT NULL, -- 主播名称
    url text NOT NULL, -- 主页 URL
    signature text, -- 主播签名
    platform varchar(20) NOT NULL, -- 主播平台
    unique_id varchar(100), -- 主播唯一 ID
    follower_count INTEGER NOT NULL, -- 粉丝数量
    following_count INTEGER NOT NULL, -- 关注数量
    like_count INTEGER NOT NULL DEFAULT 0, -- 获赞数
    video_count INTEGER NOT NULL DEFAULT 0, -- 作品数量
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE TABLE IF NOT EXISTS anchor_info_history(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    anchor_id INTEGER NOT NULL, -- 主播 ID
    anchor_name text NOT NULL, -- 主播名称
    signature text, -- 主播签名
    follower_count INTEGER NOT NULL DEFAULT 0, -- 粉丝数量
    following_count INTEGER NOT NULL DEFAULT 0, -- 关注数量
    like_count INTEGER NOT NULL DEFAULT 0, -- 获赞数
    video_count INTEGER NOT NULL DEFAULT 0, -- 作品数量
    collected_date text NOT NULL, -- 数据采集日期
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE INDEX idx_anchor_history ON anchor_info_history (
    anchor_id,
    collected_date
);