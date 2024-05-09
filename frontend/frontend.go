package frontend

import (
	"embed"
	"regexp"

	"github.com/labstack/echo/v4"
)

var (
	//go:embed all:dist
	dist embed.FS
	//go:embed dist/index.html
	indexHTML embed.FS
	//go:embed src/routeTree.gen.ts
	routeTreeFS embed.FS

	distDirFS = echo.MustSubFS(dist, "dist")
	distIndexHtml = echo.MustSubFS(indexHTML, "dist")
)

func RegisterHandlers(e *echo.Echo) {
	file := "src/routeTree.gen.ts"
	routes := collectRoutes(routeTreeFS, file)
	newRoutes := modifyRoutes(routes)

	for _, r := range newRoutes {
		e.FileFS(r, "index.html", distIndexHtml)
		// fmt.Println("registering:", r)
	}

	// this must match the routes (pages) in react
	// e.FileFS("/", "index.html", distIndexHtml)
	// e.FileFS("/profile", "index.html", distIndexHtml)
	// e.FileFS("/pokemon/*", "index.html", distIndexHtml)
	
	e.StaticFS("/", distDirFS)
}

func collectRoutes (fileFS embed.FS, file string) []string {
	routes := []string{}

	dat, err := fileFS.ReadFile(file)
    check(err)
	
	// regexp '/route':
	// r := regexp.MustCompile(`"([^{}]*)":`)
	r := regexp.MustCompile(`'([^{}]*)':`)
	matches := r.FindAllStringSubmatch(string(dat), -1)

	for _, r := range matches {
		route := r[1]
		routes = append(routes, route)
	}

	return routes
}

func modifyRoutes(routes []string) []string {
	newRoutes := []string{}
	for _, route := range routes {
		newRoutes = append(newRoutes, removeParams(route))
	}
	return newRoutes
}

func removeParams(route string) string {
	reg_dollar := regexp.MustCompile(`\$`)
	found_dollar := reg_dollar.FindStringIndex(route)

	newRoute := route

	if len(found_dollar) > 0 {
		start := found_dollar[0]
		url_after_dolar := route[start : ]

		reg_slash := regexp.MustCompile(`/`)
		found_slash := reg_slash.FindStringIndex(url_after_dolar)

		// end := len(url_after_dolar)
		if len(found_slash) > 0 {
			end := found_slash[0]
			newRoute = route[ : start] + "*" +  route[start+end : ]
		} else {
			newRoute = route[ : start] + "*"
		}
			
		// if it contains $ run it again
		found_another_dollar := reg_dollar.FindStringIndex(newRoute)
		if len(found_another_dollar) > 0 {
			newRoute = removeParams(newRoute)
		}
	}

	return newRoute
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}
