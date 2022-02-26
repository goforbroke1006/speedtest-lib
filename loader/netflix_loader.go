package loader

func NewNetflixLoader() *netflixLoader {
	return &netflixLoader{}
}

var (
	_ NetworkLoader = &netflixLoader{}
)

type netflixLoader struct {
}

func (n netflixLoader) LoadServersList() (uint, error) {
	//TODO implement me
	//panic("implement me")
	return 1, nil
}

func (n netflixLoader) Download() (bytesPerSecond float64, err error) {
	//TODO implement me
	//panic("implement me")
	return 0.1, nil
}

func (n netflixLoader) Upload() (bytesPerSecond float64, err error) {
	//TODO implement me
	//panic("implement me")
	return 0.1, nil
}
