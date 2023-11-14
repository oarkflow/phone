build:
	mkdir -p functions
	cd cmd/phoneserver && go build -ldflags "-X main.Version=`git describe --tags`" -o ../../functions/phoneserver .

fetch:
	cd cmd/update-network && yarn install && yarn fetch

update-network:
	cd cmd/update-network && go run main.go

move-file:
	mv cmd/update-network/networks.go .

network: fetch update-network move-file