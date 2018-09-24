package main

import (
	"flag"
	"os"

	"github.com/astrocorp42/astroflow-go"
	"github.com/astrocorp42/astroflow-go/log"
	_ "github.com/astrolib/godotenv/autoload"
	"gitlab.com/astrocorp/signal/api/server"
)

func checkEnv() {
	required := []string{
		"DATABASE_URL",
	}

	for _, r := range required {
		if os.Getenv(r) == "" {
			log.With("var", r).Fatal("Missing environment variable")
		}
	}

	// optional env
	port := os.Getenv("PORT")
	if port == "" {
		os.Setenv("PORT", "9090")
	}
}

func main() {
	var err error
	flagMigrate := flag.Bool("migrate", false, "auto migrate")
	flagPort := flag.String("p", os.Getenv("PORT"), "port to listen to")
	flag.Parse()

	// init
	log.Config(
		astroflow.SetFormatter(astroflow.NewConsoleFormatter()),
	)

	checkEnv()

	srv, err := server.New(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err.Error())
	}

	// migrate if the -migrate flag is provided
	if *flagMigrate {
		err = srv.AutoMigrate()
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	// run server
	log.With("port", *flagPort).Info("listenning...")
	log.Fatalf("%s", srv.Run(*flagPort))
}