module github.com/brunoscheufler/atlas/sdk/atlas-sdk-go

go 1.19

require (
	github.com/brunoscheufler/atlas v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.9.0
	google.golang.org/grpc v1.49.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20220715151400-c0bba94af5f8 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/brunoscheufler/atlas => ../../
