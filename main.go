package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	// route statis untuk mengakses folder public
	e.Static("/public", "public") // /public

	// renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t

	// Routing
	e.GET("/hello", helloWorld)           //localhost:5000/hello
	e.GET("/", home)                      //localhost:5000
	e.GET("/contact", contact)            //localhost:5000/contact
	e.GET("/blog", blog)                  //localhost:5000/blog
	e.GET("/blog-detail/:id", blogDetail) //localhost:5000/blog-detail/0 | :id = url params
	e.GET("/form-blog", formAddBlog)      //localhost:5000/form-blog

	fmt.Println("Server berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func home(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func contact(c echo.Context) error {
	return c.Render(http.StatusOK, "contact.html", nil)
}

func blog(c echo.Context) error {
	return c.Render(http.StatusOK, "blog.html", nil)
}

func blogDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) // url params | dikonversikan dari string menjadi int/integer

	data := map[string]interface{}{ // data yang akan digunakan/dikirimkan ke html menggunakan map interface
		"Id":      id,
		"Title":   "Pasar Coding di Indonesia Dinilai Masih Menjanjikan",
		"Content": "Ketimpangan sumber daya manusia (SDM) di sektor digital masih menjadi isu yang belum terpecahkan. Berdasarkan penelitian ManpowerGroup, ketimpangan SDM global, termasuk Indonesia, meningkat dua kali lipat dalam satu dekade terakhir. Lorem ipsum, dolor sit amet consectetur adipisicing elit. Quam, molestiae numquam! Deleniti maiores expedita eaque deserunt quaerat! Dicta, eligendi debitis?",
	}

	return c.Render(http.StatusOK, "blog-detail.html", data)
}

func formAddBlog(c echo.Context) error {
	return c.Render(http.StatusOK, "add-blog.html", nil)
}
