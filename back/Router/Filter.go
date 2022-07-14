package Router

import (
	"fmt"
	"net/http"
)

type Filter struct {}

var AllowedURL map[string]string

func (f *Filter) Manage(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return f.panicFilter(
			f.allowedUrl(
				handlerFunc))
}

func (f *Filter) panicFilter(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// prevent panic and return error page
		next.ServeHTTP(w, r)
	})
}

func (f *Filter) headerFilter (next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// add some security headers
		next.ServeHTTP(w, r)
	})
}

func (f *Filter) allowedUrl (next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		askedURL := r.URL.Path
		fmt.Println(askedURL)
		if _, exist := AllowedURL[askedURL]; !exist {
			//fmt.Println("Not allowed url")
		}

		next.ServeHTTP(w, r)
	})
}