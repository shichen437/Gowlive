package entity

type CpuInfo struct {
	Cpu       int32   `json:"cpu"`
	Cores     int32   `json:"cores"`
	ModelName string  `json:"modelName"`
	Mhz       float64 `json:"mhz"`
	Percent   float64 `json:"percent"`
}

type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"usedPercent"`
}

type DiskInfo struct {
	Path        string  `json:"path"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}
