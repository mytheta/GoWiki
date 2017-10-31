package main

import (
  "net/http"
  "io/ioutil"
  "text/template"
)

type Page struct {
  Title string
  Body []byte
}

func handler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[6:]
  p, _ := loadPage(title)

  t, _ := template.ParseFiles("view.html")
  t.Execute(w,p)
  // fmt.Fprintf(w, "Hi %s!", r.URL.Path[1:])
}

func editHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[6:]
  p, err := loadPage(title)
  if err != nil {
    p = &Page{Title:title}
  }

  t, _ := template.ParseFiles("edit.html")
  t.Execute(w,p)
  // fmt.Fprintf(w, "Hi %s!", r.URL.Path[1:])
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
  title := r.URL.Path[6:]
  body := r.FormValue("body")
  p := &Page{Title:title, Body:[]byte(body)}
  p.save()

  http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func (p *Page) save() error {
  filename := p.Title + ".txt"
  return ioutil.WriteFile(filename, p.Body,0600)
}

func loadPage(title string) (*Page, error) {
  filename := title + ".txt"
  body, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  return &Page{Title:title, Body:body}, nil
}

func main(){
  http.HandleFunc("/",handler)
  http.HandleFunc("/edit/",editHandler)
  http.HandleFunc("/save/",saveHandler)
  http.ListenAndServe(":8080",nil)
}
