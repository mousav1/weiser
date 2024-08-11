package facades

// Storage returns the singleton instance of StorageFacade.
func Storage() *StorageFacade {
	return NewStorageFacade()
}

// Cache returns the singleton instance of CacheFacade.
func Cache() *CacheFacade {
	return NewCacheFacade()
}

// Cookie returns the singleton instance of CookieFacade.
func Cookie() *CookieFacade {
	return NewCookieFacade()
}
