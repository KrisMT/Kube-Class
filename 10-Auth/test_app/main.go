package main

import (
  "fmt"
  "net/http"
)


func main() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, you're requested: %s\n", r.URL.Path)
  })
  http.HandleFunc("/notsecure", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "not secure")
  })
  http.HandleFunc("/secure", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "secure")
  })

  http.ListenAndServe(":4000", nil)
}
