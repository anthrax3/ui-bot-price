package main

import (
	"html/template"
	"net/http"

	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
)

type page struct {
	Title string
	Msg   string
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	title := r.URL.Path[len("/"):]

	if title != "exec/" {
		t, _ := template.ParseFiles("template.html")
		t.Execute(w, &page{Title: "Convert Image"})
	} else {
		imgfile, fhead, _ := r.FormFile("imgfile")

		img, ext, _ := image.Decode(imgfile)

		w.Header().Set("Content-type", "image/jpeg")
		w.Header().Set("Content-Disposition", "filename=\""+fhead.Filename+"."+ext+"\"")
		jpeg.Encode(w, img, &jpeg.Options{0})
	}
}

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(":7777", nil)
}
