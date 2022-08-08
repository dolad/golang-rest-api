package main

import (
	"fmt"
	"net/http"

	"github.com/dolad/rest-api/internal/comment"
	"github.com/dolad/rest-api/internal/database"
	transportHttp "github.com/dolad/rest-api/internal/transport/http"
)

// App struct  whiich contains things like pointers to database
type App struct {
}

func (app *App) Run() error {
	fmt.Println("Setting Up our App")
	var err error
	db, err := database.NewDatabase()
	if err != nil {
		return err
	}
	err = database.MigrateDB(db)

	if err != nil {
		return err
	}
	commentService := comment.NewService(db)
	handler := transportHttp.NewHandler(commentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		fmt.Println("Failed to setup server")
		return err
	}
	return nil
}

func main() {
	fmt.Println("Go rest api")
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting the rest api")
	}

}
