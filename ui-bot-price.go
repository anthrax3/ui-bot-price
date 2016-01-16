package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

	//	"image"
	//	_ "image/gif"
	//	"image/jpeg"
	//	_ "image/png"
)

////------------ Объявление типов и глобальных переменных

var (
	hd   string
	user string
)

var (
	taskT       TaskerTovar
	tekuser     string // текущий пользователь который задает условия на срабатывания
	pathcfg     string // адрес где находятся папки пользователей, если пустая строка, то текущая папка
	pathcfguser string
)

type page struct {
	Title  string
	Msg    string
	Msg2   string
	TekUsr string
}

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

//------------ END Объявление типов и глобальных переменных

// сохранить файл
func Savestrtofile(namef string, str string) int {

	file, err := os.OpenFile(namef, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0776)
	if err != nil {
		// handle the error here
		return -1
	}
	defer file.Close()

	file.WriteString(str)
	return 0
}

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

		if _, err := os.Stat(pathcfguser); os.IsNotExist(err) {
			os.Mkdir(pathcfguser, 0776)
		}
		fmt.Println(pathcfguser + string(os.PathSeparator) + shop + "-url.cfg")

		savetofilecfg(pathcfguser+string(os.PathSeparator)+shop+"-url.cfg", taskT)

		ss1 := "Введенное условие для магазина " + shop
		ss := taskT.Url + "   " + taskT.uslovie + " " + r.FormValue("schislo")
		t1, _ := template.ParseFiles("template-result.html")
		t1.Execute(w, &page{Title: "Введенное условие для магазина " + shop, Msg: ss, Msg2: ss1})

	}
}

// функция парсинга аргументов программы
func parse_args() bool {
	flag.StringVar(&hd, "hd", "", "Рабочая папка где нах-ся папки пользователей для сохранения ")
	flag.StringVar(&user, "user", "", "Рабочая папка где нах-ся папки пользователей для сохранения ")
	flag.Parse()
	pathcfg = hd
	if user == "" {
		tekuser = "testuser"
	} else {
		tekuser = user
	}
	return true
}

func main() {

	if !parse_args() {
		return
	}

	if pathcfg == "" {
		pathcfguser = tekuser
	} else {
		pathcfguser = pathcfg + string(os.PathSeparator) + tekuser
	}

	fmt.Println(pathcfguser)

	http.HandleFunc("/", index)

	http.ListenAndServe(":7777", nil)
}
