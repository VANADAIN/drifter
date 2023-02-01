package cache

type CacheManager struct {
	Barterc *BarterCache
	TJc     *TJCache
	MSGc    *MessageCache
	// arrays of cache types
}

func NewCacheManager() *CacheManager {
	return &CacheManager{
		Barterc: NewBarterCache(),
		TJc:     NewTJCache(),
		MSGc:    NewMessageCache(),
	}
}
