package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"personal-web/connection"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type Blog struct {
	ID       int
	Title    string
	Content  string
	Author   string
	PostDate time.Time
	Image    string
}

var dataBlog = []Blog{
	{
		Title:   "Ini adalah title",
		Content: "Heyy apakah saya ganteng ?",
		Author:  "Dandi Saputra",
		// PostDate: "09 Maret 2023",
	},
	{
		Title:   "Memang saya ganteng",
		Content: "VALIDDDDDD BANGET",
		Author:  "Dandi Saputra",
		// PostDate: "09 Maret 2023",
	},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	connection.DatabaseConnect()
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
	e.POST("/add-blog", addBlog)          //localhost:5000/add-blog
	e.GET("/delete-blog/:id", deleteBlog)

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
	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, content, image, post_date FROM tb_blog")

	var result []Blog
	for data.Next() {
		var each = Blog{}

		err := data.Scan(&each.ID, &each.Title, &each.Content, &each.Image, &each.PostDate)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}

		result = append(result, each)
	}

	blogs := map[string]interface{}{
		"Blogs": result,
	}
	return c.Render(http.StatusOK, "blog.html", blogs)
}

func blogDetail(c echo.Context) error {
	// http://localhost:5000/blog-detail/1
	// "1" => 1
	id, _ := strconv.Atoi(c.Param("id")) // url params | dikonversikan dari string menjadi int/integer

	var BlogDetail = Blog{}

	for i, data := range dataBlog {
		if id == i {
			BlogDetail = Blog{
				Title:    data.Title,
				Content:  data.Content,
				PostDate: data.PostDate,
				Author:   data.Author,
			}
		}
	}

	// data yang akan digunakan/dikirimkan ke html menggunakan map interface
	detailBlog := map[string]interface{}{
		"Blog": BlogDetail,
	}

	return c.Render(http.StatusOK, "blog-detail.html", detailBlog)
}

func formAddBlog(c echo.Context) error {
	return c.Render(http.StatusOK, "add-blog.html", nil)
}

func addBlog(c echo.Context) error {
	title := c.FormValue("inputTitle")
	content := c.FormValue("inputContent")

	println("Title: " + title)
	println("Content: " + content)

	var newBlog = Blog{
		Title:   title,
		Content: content,
		Author:  "Dandi Saputra",
		// PostDate: time.Now().String(),
	}

	dataBlog = append(dataBlog, newBlog)

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}

func deleteBlog(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	dataBlog = append(dataBlog[:id], dataBlog[id+1:]...)

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}

// 3
// [0, 1, 2, 3, 4, 5, 6]
// [0, 1, 2]
// [4, 5, 6]

// [0, 1, 2, 4, 5,6]
