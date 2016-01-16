package main

import (
	//	"fmt"
	"html/template"
	"net/http"
	"strconv"

	//	"image"
	//	_ "image/gif"
	//	"image/jpeg"
	//	_ "image/png"
)

// структура задания с информацией по товару
type TaskerTovar struct {
	Url string // ссылка на источник данных
	Tovar
	Tasker
}

//// структура книги
type Tovar struct {
	name          string // название товара
	price         int    // цена для всех (обычная)
	pricediscount int    // цена со скидкой которая видна
}

// задание-триггер для срабатывания оповещения
type Tasker struct {
	uslovie string // условие < , > , =
	price   int    // цена триггера
	result  bool   // результат срабатывания триггера, если true , то триггер сработал
}

var taskT TaskerTovar

type page struct {
	Title string
	Msg   string
	Msg2  string
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	title := r.URL.Path[len("/"):]

	if title != "exec/" {
		t, _ := template.ParseFiles("template.html")
		t.Execute(w, &page{Title: "Создание триггера", Msg: "Задание триггера (условия) на срабатывание бота цен"})
	} else {
		shop := r.FormValue("shop")
		taskT.Url = r.FormValue("surl")
		taskT.uslovie = r.FormValue("uslovie")
		taskT.Tasker.price, _ = strconv.Atoi(r.FormValue("schislo"))
		ss1 := "Введенное условие для магазина " + shop
		ss := taskT.Url + "   " + taskT.uslovie + " " + r.FormValue("schislo")
		t1, _ := template.ParseFiles("template-result.html")
		t1.Execute(w, &page{Title: "Введенное условие для магазина " + shop, Msg: ss, Msg2: ss1})

	}
}

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(":7777", nil)
}
