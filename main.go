package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"personal-web/connection"
	"strconv"

	// "text/template"
	"html/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type Blog struct {
	ID         int
	Title      string
	Content    string
	Author     string
	PostDate   time.Time
	Image      string
	FormatDate string
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
	e.GET("/", home)                      //localhost:5000
	e.GET("/contact", contact)            //localhost:5000/contact
	e.GET("/blog", blog)                  //localhost:5000/blog
	e.GET("/blog-detail/:id", blogDetail) //localhost:5000/blog-detail/0 | :id = url params
	e.GET("/form-blog", formAddBlog)      //localhost:5000/form-blog
	e.POST("/add-blog", addBlog)          //localhost:5000/add-blog
	e.GET("/delete-blog/:id", deleteBlog)
	// e.GET("/test", testHome)

	fmt.Println("Server berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
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
		each.Author = "Dandi Saputra"
		each.FormatDate = each.PostDate.Format("07 February 2006")

		result = append(result, each)
	}

	blogs := map[string]interface{}{
		"Blogs": result,
	}
	return c.Render(http.StatusOK, "blog.html", blogs)
}

func blogDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) // url params | dikonversikan dari string menjadi int/integer
	var BlogDetail = Blog{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, content, image, post_date FROM tb_blog WHERE id=$1", id).Scan(&BlogDetail.ID, &BlogDetail.Title, &BlogDetail.Content, &BlogDetail.Image, &BlogDetail.PostDate)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	BlogDetail.Author = "Dandi Saputra"
	BlogDetail.FormatDate = BlogDetail.PostDate.Format("02 February 2006")

	data := map[string]interface{}{
		"Blog": BlogDetail,
	}

	return c.Render(http.StatusOK, "blog-detail.html", data)
}

func formAddBlog(c echo.Context) error {
	return c.Render(http.StatusOK, "add-blog.html", nil)
}

func deleteBlog(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) //id = 2

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_blog WHERE id=$1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}

func addBlog(c echo.Context) error {
	title := c.FormValue("inputTitle")
	content := c.FormValue("inputContent")
	image := "image.png"

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_blog (title, content, image, post_date) VALUES ($1, $2, $3, $4)", title, content, image, time.Now())

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}
