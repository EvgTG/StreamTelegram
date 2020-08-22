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
tg команды - время проверки (через сколько чекать rss), ch id, несколько rss(протестить не забанят ли)
перенести loc и текста в конфиг
перебрать ошибки, + errors.Wrap, case в tg повставлять, перебрать имена пм
возможность иметь более пустой env файл
*/
