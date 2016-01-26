package main

import (
	"flag"
	//	"fmt"
	//	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/render"

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

type TTasker struct {
	Url     string
	Uslovie string // условие < , > , =
	Price   string // цена для всех (обычная)
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

//// чтение файла с именем namefи возвращение содержимое файла, иначе текст ошибки
func readfiletxt(namef string) string {
	file, err := os.Open(namef)
	if err != nil {
		return "handle the error here"
	}
	defer file.Close()
	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return "error here"
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return "error here"
	}
	return string(bs)
}

func indexHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	rr.HTML(200, "index", &page{Title: "Йоу Начало", Msg: "Начальная страница", TekUsr: "Текущий пользователь: " + string(user)})
}

func AddTaskHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	rr.HTML(200, "addtask", &page{Title: "Создание триггера", Msg: "Задание триггера (условия) на срабатывание бота цен", TekUsr: "Текущий пользователь: " + string(user)})
}

// выбор магазина который будет выбран для вывода содержимого cfg файла
func clickViewTaskHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	//	tt := make([]TTasker, 0)
	//	s := readfiletxt(pathcfguser + string(user) + string(os.PathSeparator) + "labirint-url.cfg")
	//	ss := strings.Split(s, "\n")
	//	for _, v := range ss {
	//		ts := strings.Split(v, ";")
	//		if len(ts) == 4 {
	//			tt = append(tt, TTasker{Url: ts[0], Uslovie: ts[1], Price: ts[2]})
	//		}
	//	}

	rr.HTML(200, "clickview", &page{TekUsr: string(user)})
}

// просмотр
func ViewTaskHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	shop := r.FormValue("shop")
	tt := make([]TTasker, 0)
	s := readfiletxt(pathcfguser + string(user) + string(os.PathSeparator) + shop + "-url.cfg")
	ss := strings.Split(s, "\n")
	for _, v := range ss {
		ts := strings.Split(v, ";")
		if len(ts) == 4 {
			tt = append(tt, TTasker{Url: ts[0], Uslovie: ts[1], Price: ts[2]})
		}
	}

	rr.HTML(200, "view", &tt)
}

func ExecHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	shop := r.FormValue("shop")
	taskT.Url = r.FormValue("surl")
	taskT.uslovie = r.FormValue("uslovie")
	taskT.Tasker.price, _ = strconv.Atoi(r.FormValue("schislo"))

	if _, err := os.Stat(pathcfguser + string(user)); os.IsNotExist(err) {
		os.Mkdir(pathcfguser+string(user), 0776)
	}

	savetofilecfg(pathcfguser+string(user)+string(os.PathSeparator)+shop+"-url.cfg", taskT)

	ss1 := "Введенное условие для магазина " + shop
	ss := taskT.Url + "   " + taskT.uslovie + " " + r.FormValue("schislo")
	//	fmt.Println(ss1)
	//	fmt.Println(ss)
	rr.HTML(200, "template-result", &page{Title: "Введенное условие для магазина " + shop, Msg: ss, Msg2: ss1})
}

func authFunc(username, password string) bool {
	return (auth.SecureCompare(username, "admin") && auth.SecureCompare(password, "1")) || (auth.SecureCompare(username, "mars") && auth.SecureCompare(password, "2")) || (auth.SecureCompare(username, "oilnur") && auth.SecureCompare(password, "oilnur"))
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
	m := martini.Classic()

	if !parse_args() {
		return
	}

	if pathcfg == "" {
		pathcfguser = ""
	} else {
		pathcfguser = pathcfg + string(os.PathSeparator)
	}

	//	if pathcfg == "" {
	//		pathcfguser = tekuser
	//	} else {
	//		pathcfguser = pathcfg + string(os.PathSeparator) + tekuser
	//	}

	//	fmt.Println(pathcfguser)

	m.Use(render.Renderer(render.Options{
		Directory:  "templates", // Specify what path to load the templates from.
		Layout:     "layout",    // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}}))

	m.Use(auth.BasicFunc(authFunc))

	m.Get("/", indexHandler)
	m.Get("/addtask", AddTaskHandler)
	m.Post("/exec", ExecHandler)
	m.Post("/view", ViewTaskHandler)
	m.Get("/clickview", clickViewTaskHandler)
	m.Get("/", indexHandler)
	m.RunOnAddr(":7777")

}
