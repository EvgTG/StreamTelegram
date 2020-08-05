package main

import "context"

func main() {
	app := New()
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
	defer app.Stop(context.Background())
}

/*
TODO:
перенести loc и текста в конфиг
перебрать ошибки
возможность иметь более пустой env файл
*/
