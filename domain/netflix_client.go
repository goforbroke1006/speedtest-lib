package domain

type NetflixFastClient interface {
	GetUrls() (urls []string, err error)
	Measure(urls []string, KbpsChan chan<- float64) (err error)
}
