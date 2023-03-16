package main

import (
	"Personal-Web/connection"
	"Personal-Web/middleware"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"text/template"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}

type SessionData struct {
	IsLogin bool
	Name    string
}

var userData = SessionData{}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	connection.DatabaseConnect()
	e := echo.New()

	// route statis untuk mengakses folder public
	e.Static("/public", "public") // /public
	e.Static("/upload", "upload")
	// to use sessions using echo
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("session"))))

	// renderer
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	e.Renderer = t

	// Routing
	e.GET("/", home)                                    //localhost:5000
	e.GET("/contact", contact)                          //localhost:5000/contact
	e.GET("/blog", blog)                                //localhost:5000/blog
	e.GET("/blog-detail/:id", blogDetail)               //localhost:5000/blog-detail/0 | :id = url params
	e.GET("/form-blog", formAddBlog)                    //localhost:5000/form-blog
	e.POST("/add-blog", middleware.UploadFile(addBlog)) //localhost:5000/add-blog
	e.GET("/delete-blog/:id", deleteBlog)
	e.GET("/form-login", formLogin)
	e.POST("/login", login)
	e.GET("/form-register", formRegister)
	e.POST("/register", addRegister)
	e.GET("/logout", logout)

	fmt.Println("Server berjalan di port 5000")
	e.Logger.Fatal(e.Start("localhost:5000"))
}

func home(c echo.Context) error {
	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["isLogin"],
		"FlashMessage": sess.Values["message"],
		"FlashName":    sess.Values["name"],
	}
	delete(sess.Values, "message")
	delete(sess.Values, "status")
	sess.Save(c.Request(), c.Response())
	return c.Render(http.StatusOK, "index.html", flash)
}

func contact(c echo.Context) error {
	c.Response().Header().Set("Content-type", "text/html; charset=UTF-8")
	return c.Render(http.StatusOK, "contactme.html", nil)
}

func blog1(c echo.Context) error {
	return c.Render(http.StatusOK, "blog.html", nil)
}

func blog(c echo.Context) error {
	sess, _ := session.Get("session", c)

	if sess.Values["isLogin"] != true {
		userData.IsLogin = false
	} else {
		userData.IsLogin = sess.Values["isLogin"].(bool)
		userData.Name = sess.Values["name"].(string)
	}
	data, _ := connection.Conn.Query(context.Background(), "SELECT tb_blog.id, title, content, image, post_date, tb_user.name as author FROM tb_blog LEFT JOIN tb_user ON tb_blog.author = tb_user.id ORDER BY id DESC")

	var result []Blog
	for data.Next() {
		var each = Blog{}

		err := data.Scan(&each.ID, &each.Title, &each.Content, &each.Image, &each.PostDate, &each.Author)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}

		each.FormatDate = each.PostDate.Format("07 February 2006")

		result = append(result, each)
	}

	blogs := map[string]interface{}{
		"Blogs":       result,
		"DataSession": userData,
	}
	return c.Render(http.StatusOK, "blog.html", blogs)
}

func blogDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id")) // url params | dikonversikan dari string menjadi int/integer
	var BlogDetail = Blog{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT tb_blog.id, title, content, image, post_date, tb_user.name as author FROM tb_blog LEFT JOIN tb_user ON tb_blog.author = tb_user.id WHERE tb_blog.id=$1", id).Scan(&BlogDetail.ID, &BlogDetail.Title, &BlogDetail.Content, &BlogDetail.Image, &BlogDetail.PostDate, &BlogDetail.Author)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

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
	title := c.FormValue("title")
	content := c.FormValue("content")
	image := c.Get("dataFile").(string) // => image-982349187nfjka.png

	sess, _ := session.Get("session", c)
	authorId := sess.Values["id"]

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_blog (title, content, image, post_date, author) VALUES ($1, $2, $3, $4, $5)", title, content, image, time.Now(), authorId)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/blog")
}

func formRegister(c echo.Context) error {
	tmpl, err := template.ParseFiles("views/register.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), nil)
}

func addRegister(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	name := c.FormValue("name")
	email := c.FormValue("email")
	password := c.FormValue("password")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(), "INSERT INTO tb_user (name, email, password) VALUES ($1, $2, $3)", name, email, passwordHash)
	if err != nil {
		fmt.Println(err)
		redirectWithMessage(c, "Register failed, please try again", false, "/form-register")
	}

	return redirectWithMessage(c, "Register success", true, "/form-login")
}

func formLogin(c echo.Context) error {
	sess, _ := session.Get("session", c)

	flash := map[string]interface{}{
		"FlashStatus":  sess.Values["status"],
		"FlashMessage": sess.Values["message"],
	}

	delete(sess.Values, "message")
	delete(sess.Values, "status")
	tmpl, err := template.ParseFiles("views/login.html")

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return tmpl.Execute(c.Response(), flash)
}

func login(c echo.Context) error {
	err := c.Request().ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := c.FormValue("email")
	password := c.FormValue("password")

	user := User{}
	err = connection.Conn.QueryRow(context.Background(), "SELECT * FROM tb_user WHERE email=$1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return redirectWithMessage(c, "Email Salah !", false, "/form-login")
	}

	fmt.Println(user)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return redirectWithMessage(c, "Password Salah !", false, "/form-login")
	}

	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = 10800 //3 jam
	sess.Values["message"] = "Login Success !"
	sess.Values["status"] = true //show alert
	sess.Values["name"] = user.Name
	sess.Values["id"] = user.ID
	sess.Values["isLogin"] = true //akses login
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func logout(c echo.Context) error {
	sess, _ := session.Get("session", c)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusTemporaryRedirect, "/")
}

func redirectWithMessage(c echo.Context, message string, status bool, path string) error {
	sess, _ := session.Get("session", c)
	sess.Values["message"] = message
	sess.Values["status"] = status
	sess.Save(c.Request(), c.Response())
	return c.Redirect(http.StatusMovedPermanently, path)
}
