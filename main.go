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
команды телеги: channel ID( + кнопка  на изменение, в дальнейшем несколько каналов)
перенести loc и текста в конфиг
время полследнего обновления в lastrss
перебрать ошибки
возможность иметь более пустой env файл
*/
