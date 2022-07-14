package Router

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var regexpComment = regexp.MustCompile("#.+")
var regexpStaticPath = regexp.MustCompile("(StaticPath[ ]+)(.+)")
var regexpRs = regexp.MustCompile("(GET|POST|PUT)[ ]+(.+)[ ]+(.+)")
var regexpParams = regexp.MustCompile(":([a-zA-Z]+)")
var regexpFaviconPath = regexp.MustCompile("(FaviconPath[ ]+)(.+)")

var StaticPath = ""
var FaviconPath = ""

type Router struct {}

type Route struct{
	Method string
	Path string // Request
	BindPath string // thing to response
	Params map[string]int
}

var Config []Route

//because you cannot declare map as const
var routeBinging = func () map[string] func(w http.ResponseWriter, req *http.Request) {
	return map[string]func(w http.ResponseWriter, req *http.Request){
		"index": Index,
		"user": User,
	}
}

func (r *Router) mainRouting(w http.ResponseWriter, req *http.Request) {

	if len(regexp.MustCompile(`static.+`).FindString(req.URL.Path)) > 0 {
		ManageStatic(w, req)
		return
	} else if len(regexp.MustCompile(`favicon\.ico+`).FindString(req.URL.Path)) > 0 {
		fileBytes, err := ioutil.ReadFile(FaviconPath)
		if err != nil {
			fmt.Println(err)
		}

		_, err = w.Write(fileBytes)
		if err != nil {
			fmt.Println(err)
		}
	}

	for _, route := range Config {
		if route.Method == req.Method {
			if route.Path == req.URL.Path {
				routeBinging()[route.BindPath](w, req)
			}
		}
	}
}

func (r *Router) parseParams(url string) map[string]int {
	tmpMap := make(map[string]int)
	tmp := regexpParams.FindAllStringSubmatch(url, -1)
	for i, match := range tmp {
		if len(match) > 1 {
			if _, exist := tmpMap[match[1]]; exist {
				panic("you cannot use same params in one URI")
			} else {
				tmpMap[match[1]] = i
			}
		}
	}
	return tmpMap
}

func (r *Router) parseConfig ()  {
	Config = []Route{}
	file, err := os.Open("./back/Router/routes.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		rs := scanner.Text()
		if idx := regexpComment.FindIndex([]byte(scanner.Text())); idx != nil {
			rs = rs[:idx[0]]
		}
		if len(rs) < 1 {
			continue
		}
		switch true {
			case len(regexpRs.FindString(rs)) > 0: {
				groups := regexpRs.FindStringSubmatch(rs)
				tmpMap := r.parseParams(groups[2])
				tmp := Route{
					Method:         groups[1],
					Path:           groups[2],
					BindPath:       groups[3],
					Params:        tmpMap,
				}
				Config = append(Config, tmp)
			}
			case len(regexpStaticPath.FindString(rs)) > 0: {
				StaticPath = regexpStaticPath.FindStringSubmatch(rs)[2]
				// add new case for some variables
			}
			case len(regexpFaviconPath.FindString(rs)) > 0: {
				FaviconPath = regexpFaviconPath.FindStringSubmatch(rs)[2]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func (r *Router) Manage ()  {
	f := new(Filter)
	r.parseConfig() // initialization
	fmt.Println(Config)
	fmt.Println(StaticPath)
	fmt.Println(FaviconPath)
	http.HandleFunc("/", f.Manage(r.mainRouting)) // handle all incoming requests
	// "/" means every path. Filter incoming request then Route.
}