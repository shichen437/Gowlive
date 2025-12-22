package manager

import "sync"

type LogItem struct {
	Time  int64  `json:"time"`
	Level string `json:"level"`
	Msg   string `json:"msg"`
}

type LogBuffer struct {
	size int
	data []LogItem
	idx  int
	mu   sync.RWMutex
}

var (
	logBufferManager *LogBuffer
	logBufferOnce    sync.Once
)

func GetLogBufferManager() *LogBuffer {
	logBufferOnce.Do(func() {
		logBufferManager = &LogBuffer{
			size: 2000,
			data: make([]LogItem, 2000),
		}
	})
	return logBufferManager
}

func (b *LogBuffer) Append(item LogItem) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data[b.idx%b.size] = item
	b.idx++
}

func (b *LogBuffer) List(since int64, limit int) (list []LogItem, next int64) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	count := min(b.idx, b.size)
	if count == 0 {
		return nil, since
	}
	startIdx := 0
	if b.idx > b.size {
		startIdx = b.idx % b.size
	}

	list = make([]LogItem, 0, limit)

	for i := range count {
		currentIdx := (startIdx + i) % b.size
		item := b.data[currentIdx]

		if item.Time > since {
			list = append(list, item)
			if len(list) >= limit {
				break
			}
		}
	}

	if len(list) > 0 {
		next = list[len(list)-1].Time
	} else {
		next = since
	}
	return
}
