package main

import (
	"flag"
	"fmt"
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
	//	taskT       TaskerTovar
	tekuser     string // текущий пользователь который задает условия на срабатывания
	pathcfg     string // адрес где находятся папки пользователей, если пустая строка, то текущая папка
	pathcfguser string
)

type page struct {
	Title  string
	Msg    string
	Msg2   string
	TekUsr string
	Shops  map[string]string
}

type TTasker struct {
	Url     string // урл товара в магазине
	Uslovie string // условие < , > , =
	Price   string // цена для всех (обычная)
	Shop    string // название магазина
	Nstr    string // номер строки который редактируется или удаляется, если -1 то новое условие
}

//------------ END Объявление типов и глобальных переменных

// сохранить в новый файл
func SaveNewstrtofile(namef string, str string) int {
	file, err := os.Create(namef)
	if err != nil {
		// handle the error here
		return -1
	}
	defer file.Close()

	file.WriteString(str)
	return 0
}

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

//сохранение данных из TTasker товара в файл с именем namef
func savetofilecfg(namef string, t TTasker) {
	str := t.Url + ";" + t.Uslovie + ";" + t.Price + ";" + "\n"
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
	pp := page{Title: "Создание триггера", Msg: "Задание триггера (условия) на срабатывание бота цен", TekUsr: "Текущий пользователь: " + string(user)}
	pp.Shops = map[string]string{
		"labirint":   "Лабиринт",
		"ulmart":     "Юлмарт",
		"citilink":   "Ситилинк",
		"dns":        "ДНС",
		"mvideo":     "МВидео",
		"eldorado":   "Эльдорадо",
		"aliexpress": "Алиэкспресс",
	}
	rr.HTML(200, "addtask", &pp)
}

// обработка редактирования задания
func EditTaskHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	nstr, _ := strconv.Atoi(params["nstr"])
	var tt TTasker
	shop := params["shop"]
	fmt.Println(shop)

	namef := pathcfguser + string(user) + string(os.PathSeparator) + shop + "-url.cfg"

	s := readfiletxt(namef)
	ss := strings.Split(s, "\n")

	ts := strings.Split(ss[nstr], ";")
	if len(ts) == 4 {
		tt = TTasker{Url: ts[0], Uslovie: ts[1], Price: ts[2], Shop: shop, Nstr: params["nstr"]}
	}

	rr.HTML(200, "edit", &tt)
}

// обработка удаления задания
func DelTaskHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	nstr, _ := strconv.Atoi(params["nstr"])
	//	var tt TTasker
	shop := params["shop"]
	//	fmt.Println(shop)
	namef := pathcfguser + string(user) + string(os.PathSeparator) + shop + "-url.cfg"
	s := readfiletxt(namef)
	ss := strings.Split(s, "\n")
	//	if (nstr >= 0) && (nstr < len(ss)) {
	str := ""
	for i, v := range ss {
		if i != nstr {
			str += v + "\n"
		}
	}
	fmt.Println("str= ", str)
	SaveNewstrtofile(namef, str)
	//	}
	rr.Redirect("/")
}

// выбор магазина который будет выбран для вывода содержимого cfg файла
func clickViewTaskHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	pp := page{TekUsr: string(user)}
	pp.Shops = map[string]string{
		"labirint":   "Лабиринт",
		"ulmart":     "Юлмарт",
		"citilink":   "Ситилинк",
		"dns":        "ДНС",
		"mvideo":     "МВидео",
		"eldorado":   "Эльдорадо",
		"aliexpress": "Алиэкспресс",
	}
	rr.HTML(200, "clickview", &pp)
}

// просмотр заданий выбранного магазина
func ViewTaskHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request) {
	shop := r.FormValue("shop")
	tt := make([]TTasker, 0)
	s := readfiletxt(pathcfguser + string(user) + string(os.PathSeparator) + shop + "-url.cfg")
	ss := strings.Split(s, "\n")
	for _, v := range ss {
		ts := strings.Split(v, ";")
		if len(ts) == 4 {
			tt = append(tt, TTasker{Url: ts[0], Uslovie: ts[1], Price: ts[2], Shop: shop})
		}
	}
	rr.HTML(200, "view", &tt)
}

func ExecHandler(user auth.User, rr render.Render, w http.ResponseWriter, r *http.Request, params martini.Params) {
	var tt TTasker
	nstr, _ := strconv.Atoi(params["nstr"])
	//	shop := params["shop"]
	shop := r.FormValue("shop")

	tt.Shop = r.FormValue("shop")
	tt.Url = r.FormValue("surl")
	tt.Uslovie = r.FormValue("uslovie")
	tt.Price = r.FormValue("schislo")

	if _, err := os.Stat(pathcfguser + string(user)); os.IsNotExist(err) {
		os.Mkdir(pathcfguser+string(user), 0776)
	}
	namef := pathcfguser + string(user) + string(os.PathSeparator) + shop + "-url.cfg"

	if nstr == -1 {
		savetofilecfg(namef, tt)
	} else {
		s := readfiletxt(namef)
		ss := strings.Split(s, "\n")
		if nstr <= (len(ss) - 1) {
			ss[nstr] = tt.Url + ";" + tt.Uslovie + ";" + tt.Price + ";"
		}
		str := ""
		for _, v := range ss {
			if v != "" {
				str += v + "\n"
			}
		}
		SaveNewstrtofile(namef, str)
	}

	ss1 := "Введенное условие для магазина " + shop
	ss := tt.Url + "   " + tt.Uslovie + " " + tt.Price
	rr.HTML(200, "exec", &page{Title: "Введенное условие для магазина " + shop, Msg: ss, Msg2: ss1})
}

func authFunc(username, password string) bool {
	return (auth.SecureCompare(username, "admin") && auth.SecureCompare(password, "!!!!VVjhsdsajdbabjd1")) || (auth.SecureCompare(username, "mars") && auth.SecureCompare(password, "Verbat1mert")) || (auth.SecureCompare(username, "oilnur") && auth.SecureCompare(password, "Verbat1mqwe"))
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

	m.Use(render.Renderer(render.Options{
		Directory:  "templates", // Specify what path to load the templates from.
		Layout:     "layout",    // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}}))

	m.Use(auth.BasicFunc(authFunc))

	m.Get("/", indexHandler)
	m.Get("/addtask", AddTaskHandler)
	m.Get("/edit/:shop/:nstr", EditTaskHandler)
	m.Get("/del/:shop/:nstr", DelTaskHandler)
	m.Post("/exec/:shop/:nstr", ExecHandler)
	m.Post("/view", ViewTaskHandler)
	m.Get("/clickview", clickViewTaskHandler)
	m.Get("/", indexHandler)
	m.RunOnAddr(":7777")

}
