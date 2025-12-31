CREATE TABLE IF NOT EXISTS file_sync_task(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    path text NOT NULL, -- 文件路径
    filename text NOT NULL, -- 文件名称
    sync_path text NOT NULL, -- 同步路径
    duration int NOT NULL DEFAULT 0, -- 耗时
    status tinyint(1) NOT NULL DEFAULT 0, -- 同步状态:0 待同步 1 同步中 2 同步成功 3 同步失败
    remark text, -- 备注
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE INDEX idx_file_sync_path ON file_sync_task (
    path,
    filename
);
CREATE INDEX idx_file_sync_status ON file_sync_task (
    status
);

ALTER TABLE live_manage ADD COLUMN sync_path text;
