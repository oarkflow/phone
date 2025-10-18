fetch:
	cd cmd/update-network && yarn install && yarn fetch

update-network:
	cd cmd/update-network && go run main.go

move-file:
	mv cmd/update-network/networks.go .

network: fetch update-network move-file
