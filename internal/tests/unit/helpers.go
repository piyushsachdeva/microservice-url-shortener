package unit

import (
	"context"

	"github.com/itsbaivab/url-shortener/internal/adapters/cache"
	"github.com/itsbaivab/url-shortener/internal/core/domain"
)

func FillCache(cache *cache.RedisCache, links []domain.Link) error {
	for _, link := range links {
		err := cache.Set(context.Background(), link.Id, link.OriginalURL)
		if err != nil {
			return err
		}
	}
	return nil
}
