package main

import (
	//	"fmt"
	"html/template"
	"net/http"
	"os"
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

var tekuser string // текущий пользователь который задает условия на срабатывания

type page struct {
	Title  string
	Msg    string
	Msg2   string
	TekUsr string
}

func Savestrtofile(namef string, str string) int {
	file, err := os.Create(namef)
	if err != nil {
		// handle the error here
		return -1
	}
	defer file.Close()

	file.WriteString(str)
	return 0
}

//http://www.labirint.ru/books/408438/;<;1534;
//сохранение данных из Tasker товара в файл с именем namef
func savetofilecfg(namef string, t TaskerTovar) {
	str := t.Url + ";" + t.uslovie + ";" + strconv.Itoa(t.Tasker.price) + ";" + "\n"
	Savestrtofile(namef, str)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	title := r.URL.Path[len("/"):]

	if title != "exec/" {
		t, _ := template.ParseFiles("template.html")
		t.Execute(w, &page{Title: "Создание триггера", Msg: "Задание триггера (условия) на срабатывание бота цен", TekUsr: "Текущий пользователь: " + tekuser})
	} else {
		shop := r.FormValue("shop")
		taskT.Url = r.FormValue("surl")
		taskT.uslovie = r.FormValue("uslovie")
		taskT.Tasker.price, _ = strconv.Atoi(r.FormValue("schislo"))

		savetofilecfg(shop+"-url.cfg", taskT)

		ss1 := "Введенное условие для магазина " + shop
		ss := taskT.Url + "   " + taskT.uslovie + " " + r.FormValue("schislo")
		t1, _ := template.ParseFiles("template-result.html")
		t1.Execute(w, &page{Title: "Введенное условие для магазина " + shop, Msg: ss, Msg2: ss1})

	}
}

func main() {
	tekuser = "testuser"

	http.HandleFunc("/", index)

	http.ListenAndServe(":7777", nil)
}
