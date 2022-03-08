package speedtest_lib

//go:generate mockgen -source=./domain/domain.go         -destination=./mocks/domain.generated.go         -package=mocks
//go:generate mockgen -source=./domain/netflix_client.go -destination=./mocks/netflix_client.generated.go -package=mocks
//go:generate mockgen -source=./domain/ookla_client.go   -destination=./mocks/ookla_client.generated.go   -package=mocks
