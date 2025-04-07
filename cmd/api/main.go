package main

import (
	"github.com/williamnoble/kube-botany/pkg/plant"
	"github.com/williamnoble/kube-botany/pkg/server"
	"time"
)

func main() {
	svr := server.NewServer(populatePlants())
	if err := svr.Start(8080); err != nil {
		panic(err)
	}
}

func populatePlants() []*plant.Plant {
	var plants []*plant.Plant
	plants = append(plants, plant.NewPlant(
		"DefaultBonsai123",
		"my-bonsai",
		plant.Bonsai,
		time.Now(),
		false))
	plants = append(plants, plant.NewPlant(
		"DefaultSunflower234",
		"my-sunflower",
		plant.Sunflower,
		time.Now(),
		false))
	return plants
}
