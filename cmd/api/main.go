package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/williamnoble/kube-botany/db/sqlc"
	"github.com/williamnoble/kube-botany/pkg/plant"
	"github.com/williamnoble/kube-botany/pkg/server"
	"log"
	"reflect"
	"time"
)

func main() {
	svr := server.NewServer(populatePlants())
	//_ = connectPostgres()
	go svr.BackgroundTasks()
	if err := svr.Start(8090); err != nil {
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

func connectPostgres() error {
	ctx := context.Background()

	connPool, err := pgxpool.New(ctx, "postgresql://postgres:postgres@localhost:5432/botany?sslmode=disable")
	if err != nil {
		panic(err)
	}

	queries := db.New(connPool)

	// list all authors
	authors, err := queries.Listplants(ctx)
	if err != nil {
		return err
	}
	log.Println(authors)

	params := db.CreatePlantParams{
		Name:                 "foo",
		CanDie:               pgtype.Bool{Bool: false},
		WaterConsumptionRate: pgtype.Numeric{Exp: 100},
		MinimumWaterLevel:    pgtype.Numeric{Exp: 20},
		WaterLevel:           pgtype.Numeric{Exp: 100},
		LastWatered:          pgtype.Timestamptz{Time: time.Now()},
		GrowthRate:           pgtype.Numeric{Exp: 100},
		Growth:               pgtype.Numeric{Exp: 100},
		GrowthStage:          pgtype.Text{String: "seeding"},
		LastUpdated:          pgtype.Timestamptz{Time: time.Now()},
		Backdrop:             pgtype.Text{String: "foo"},
		Mascot:               pgtype.Text{String: "foo"},
	}
	// create an author
	insertedPlant, err := queries.CreatePlant(ctx, params)
	if err != nil {
		return err
	}
	log.Println("something worked", insertedPlant.Name)
	log.Println(insertedPlant)

	// get the author we just inserted
	fetchedPlant, err := queries.GetPlant(ctx, insertedPlant.ID)
	if err != nil {
		return err
	}

	// prints true
	log.Println(reflect.DeepEqual(insertedPlant, fetchedPlant))
	return nil
}
