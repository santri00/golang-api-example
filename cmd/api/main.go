package main

import (
	"context"
	"github.com/santri00/golang-api-example/pkg/api"
	"github.com/santri00/golang-api-example/pkg/services/vehicle"
)

func main() {
	apiServer := api.NewServer(context.Background(),
		vehicle.NewService())
	apiServer.Run()
}
