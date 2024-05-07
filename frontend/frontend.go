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
	// routes := collectRoutes("frontend/src/routeTree.gen.ts")
	routes := collectRoutes(routeTreeFS)

	for _, r := range routes {
		e.FileFS(r, "index.html", distIndexHtml)
		// fmt.Println("registering", r)
	}

	// this must match the routes (pages) in react
	// e.FileFS("/", "index.html", distIndexHtml)
	// e.FileFS("/profile", "index.html", distIndexHtml)
	// e.FileFS("/pokemon/*", "index.html", distIndexHtml)
	
	e.StaticFS("/", distDirFS)
}

// func collectRoutes (file string) []string {
func collectRoutes (fileFS embed.FS) []string {
	routes := []string{}

	dat, err := fileFS.ReadFile("src/routeTree.gen.ts")
    check(err)

	// regexp '/route':
	// r := regexp.MustCompile(`'([^{}]*)':`)
	r := regexp.MustCompile(`"([^{}]*)":`)
	matches := r.FindAllStringSubmatch(string(dat), -1)

	reg := regexp.MustCompile(`\$`)
	for _, v := range matches {
		route := v[1]
		found := reg.FindStringIndex(route)
		start := -1
		if len(found) > 0 {
			start = found[0]
		}

		if start != -1 {
			route = route[0 : start] + "*"
		}

		routes = append(routes, route)
	}
	return routes
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

