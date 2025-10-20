export interface Overview {
    liveRoomCount: number;
    recordingRoomCount: number;
    recordTimeCount: number;
    unreadMessageCount: number;
    parseMediaCount: number;
}

export interface MonitorInfo {
    cpu: CpuInfo;
    mem: MemoryInfo;
    disk: DiskInfo;
}

export interface CpuInfo {
    cpu: number;
    cores: number;
    modelName: string;
    mhz: number;
    percent: number;
}

export interface MemoryInfo {
    total: number;
    used: number;
    available: number;
    usedPercent: number;
}

export interface DiskInfo {
    path: string;
    fstype: string;
    total: number;
    free: number;
    used: number;
    usedPercent: number;
}

export interface MetricsResult {
    data: Map<string, MetricsData>;
}

export interface MetricsData {
    platform: string;
    totalRequests: number;
    totalErrors: number;
    totalPercent: number;
    mainRequests: number;
    mainErrors: number;
    mainPercent: number;
}
