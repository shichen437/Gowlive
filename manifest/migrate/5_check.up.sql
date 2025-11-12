CREATE TABLE IF NOT EXISTS file_check_task(
    id INTEGER PRIMARY KEY AUTOINCREMENT,  -- ID
    path text NOT NULL, -- 文件路径
    filename text NOT NULL, -- 文件名称
    duration int NOT NULL DEFAULT 0, -- 耗时
    progress tinyint(1) NOT NULL DEFAULT 0, -- 检查进度
    file_status tinyint(1) NOT NULL DEFAULT 0, -- 文件状态
    created_at datetime(0) NOT NULL, -- Created Time
    updated_at datetime(0) DEFAULT NULL -- Updated Time
);

CREATE INDEX idx_file_check_path ON file_check_task (
    path,
    filename
);
CREATE INDEX idx_file_check_status ON file_check_task (
    progress,
    file_status
);
