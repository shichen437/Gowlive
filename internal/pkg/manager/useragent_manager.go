package manager

import (
	"errors"
	"sync"

	"github.com/shichen437/gowlive/internal/pkg/alg/lru"
	"github.com/shichen437/gowlive/internal/pkg/alg/shuffle"
)

type UserAgentManager struct {
	sharedPool []string
	shuffled   []uint32
	platforms  map[string]*lru.LRUIds
	defaultCap int
	mu         sync.RWMutex
}

func GetUserAgentManager() *UserAgentManager {
	useragentOnce.Do(func() {
		poolCopy := PCUserAgentPool
		n := len(poolCopy)
		ids := make([]uint32, n)
		for i := range n {
			ids[i] = uint32(i)
		}
		shuffled := shuffle.ShuffleArrUint32(ids)
		userAgentManager = &UserAgentManager{
			sharedPool: poolCopy,
			shuffled:   shuffled,
			platforms:  make(map[string]*lru.LRUIds),
			defaultCap: len(poolCopy),
		}
	})
	return userAgentManager
}

func (m *UserAgentManager) GetUserAgent(platform string) (string, error) {
	m.mu.RLock()
	pool := m.sharedPool
	lru := m.platforms[platform]
	m.mu.RUnlock()

	if lru == nil {
		if err := m.ensurePlatform(platform); err != nil {
			return "", err
		}
		m.mu.RLock()
		pool = m.sharedPool
		lru = m.platforms[platform]
		m.mu.RUnlock()
	}
	id, ok := lru.MoveTailToHead()
	if !ok {
		if err := m.ensurePlatform(platform); err != nil {
			return "", err
		}
		id, ok = lru.MoveTailToHead()
		if !ok {
			return "", errors.New("platform LRU is empty")
		}
	}
	if int(id) >= len(pool) {
		return "", errors.New("id out of range")
	}
	return pool[int(id)], nil
}

func (m *UserAgentManager) ensurePlatform(platform string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.platforms[platform]; ok {
		return nil
	}
	lru := lru.NewLRUIds(m.defaultCap)
	for i := range m.defaultCap {
		lru.PushBack(m.shuffled[i])
	}
	m.platforms[platform] = lru
	return nil
}

var (
	userAgentManager *UserAgentManager
	useragentOnce    sync.Once
	PCUserAgentPool  = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36 Edg/114.0.1823.67",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 Edg/125.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36 Edg/139.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36 Edg/126.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 Edg/124.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36 OPR/110.0.0.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36 OPR/109.0.0.0",
		"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.5 Safari/605.1.15",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 12_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.6 Safari/605.1.15",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.5; rv:126.0) Gecko/20100101 Firefox/126.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 14.5; rv:125.0) Gecko/20100101 Firefox/125.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv,2.0.1) Gecko/20100101 Firefox/4.0.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:127.0) Gecko/20100101 Firefox/127.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0) Gecko/20100101 Firefox/125.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:102.0) Gecko/20100101 Firefox/102.0",
		"Mozilla/5.0 (X11; Fedora; Linux x86_64; rv:120.0) Gecko/20100101 Firefox/120.0",
		"Mozilla/5.0 (Windows NT 6.1; rv,2.0.1) Gecko/20100101 Firefox/4.0.1",
	}
)
