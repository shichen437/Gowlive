package consts

const (
	MultiLogin        = true
	ServerName        = "Gowlive Api"
	TokenType         = "Bearer"
	Timeout           = 86400
	GTokenAdminPrefix = "SYS-ADMIN-"

	ContextKey     = "ContextKey"
	UserSessionKey = "UserSessionKey"

	Success = "success"
	Error   = "error"

	CtxAdminId   = "CtxAdminId"
	CtxAdminName = "CtxAdminName"

	DefaultAdminId  = 1
	DefaultPassword = "gowlive"

	StatusActive  = 1
	StatusDisable = 0

	StorageThreshold = 90

	SKFilenameTemplate    = "sk_filename_template"
	SKArchiveStrategy     = "sk_archive_strategy"
	SKLiveEndNotify       = "sk_live_end_notify"
	SKDiskProtection      = "sk_disk_protection"
	SKAutoCleanLittleFile = "sk_auto_clean_little_file"
	SKFixedResolution     = "sk_fixed_resolution"

	SKDataSyncEnable      = "sk_data_sync_enable"
	SKDataSyncFailedRetry = "sk_data_sync_failed_retry"
	SKDataSyncAutoDelete  = "sk_data_sync_auto_delete"
)
