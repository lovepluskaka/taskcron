module github.com/lovepluskaka/taskcron

go 1.12

require (
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/golang/protobuf v1.2.0
	github.com/onsi/ginkgo v1.10.1 // indirect
	github.com/onsi/gomega v1.7.0 // indirect
	google.golang.org/grpc v0.0.0-00010101000000-000000000000
)

replace google.golang.org/grpc => github.com/grpc/grpc-go v1.23.1
