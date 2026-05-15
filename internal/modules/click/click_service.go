package click

import (
	"kira-url/internal/database/models"
	"sync"
)

type ClickService struct {
	mutex  sync.Mutex
	clicks map[string]int64
}

func NewClickService() *ClickService {
	clicks := make(map[string]int64)
	return &ClickService{
		clicks: clicks,
	}
}

func (cs *ClickService) IncrementClicks(code string) {
	cs.mutex.Lock()
	cs.clicks[code]++
	cs.mutex.Unlock()
}

func (cs *ClickService) FlushClicks() []models.URL {
	// In here we block the sync to prevent race conditions
	cs.mutex.Lock()
	// Get the current map  with a temporal variable
	temp := cs.clicks
	// Assign a new map
	cs.clicks = make(map[string]int64)
	// Unlock
	cs.mutex.Unlock()

	// Is necessary do this beacuse if we do it in side the mutex.Lock
	// Is a lot of time, is a better idea do it outside. Diference of miliseconds but is important
	var urls []models.URL
	for code, amount := range temp {
		urls = append(urls, models.URL{
			ShortURL:   code,
			VisitCount: amount,
		})
	}

	return urls
}
