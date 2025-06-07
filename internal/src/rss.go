package src

import (
	"context"

	"github.com/MRXTOOR/RustFunBook/internal/model"
	"github.com/SlyMarbo/rss"
	"github.com/samber/lo"
)

type RSSSource struct {
	URL        string
	SourceName string
	SourceID   int64
}

func NewRSSSource(m model.Source) RSSSource {
	return RSSSource{
		URL:        m.FeedURL,
		SourceName: m.Name,
		SourceID:   m.ID,
	}
}

func (s RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errCh <- err
			return
		}
		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}
}

func (s RSSSource) Fetch(ctx context.Context) ([]model.Item, error) {
	items, err := func() ([]model.Item, error) {
		feed, err := s.loadFeed(ctx, s.URL)
		if err != nil {
			return nil, err
		}
		return lo.Map(feed.Items, func(item *rss.Item, _ int) model.Item {
			return model.Item{
				Title:      item.Title,
				Categories: item.Categories,
				Link:       item.Link,
				Date:       item.Date,
				Summary:    item.Summary,
				SourceName: s.SourceName,
			}
		}), nil
	}()
	return items, err
}
