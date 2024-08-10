package facades

import "github.com/mousav1/weiser/app/cache"

// Cache returns the singleton instance of CacheFacade.
func Cache() *cache.CacheFacade {
	return &cache.CacheFacade{}
}
