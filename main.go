package main

import "context"

func main() {
	app := New()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	defer app.Stop(context.Background())
}

//TODO добавить поиск id по нику
//TODO динамический id, права доступа
