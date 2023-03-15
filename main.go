package main

import (
	"Personal-web/connection"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type Projects struct {
	Id          int
	Title       string
	Sdate       time.Time
	Edate       time.Time
	Duration    string
	Descript    string
	Technologys []string
	Tech1       bool
	Tech2       bool
	Tech3       bool
	Tech4       bool
	Image       string
}

var dataProject = []Projects{
	{
		Title: "Mobile App 2019",
		// Sdate:    "26 January 2019",
		// Edate:    "05 March 2019",
		Duration: "3 Month",
		Descript: "App that used for dumbways student, it was deployed and can downloaded on playstore.<br />Happy download",
		Tech1:    true,
		Tech2:    false,
		Tech3:    true,
		Tech4:    true,
	},
	{
		Title: "Web App 2020",
		// Sdate:    "26 August 2020",
		// Edate:    "05 December 2020",
		Duration: "2 Month",
		Descript: "App that used for dumbways student, it was deployed and can downloaded on playstore.<br />Happy download",
		Tech1:    false,
		Tech2:    true,
		Tech3:    true,
		Tech4:    true,
	},
	{
		Title: "Web App 2023",
		// Sdate:    "26 March 2023",
		// Edate:    "05 September 2023",
		Duration: "3 Month",
		Descript: "App that used for dumbways student, it was deployed and can downloaded on playstore.<br />Happy download",
		Tech1:    true,
		Tech2:    true,
		Tech3:    false,
		Tech4:    true,
	},
	{
		Title: "Mobile App 2024",
		// Sdate:    "26 Jully 2024",
		// Edate:    "05 November 2024",
		Duration: "5 Month",
		Descript: "App that used for dumbways student, it was deployed and can downloaded on playstore.<br />Happy download",
		Tech1:    true,
		Tech2:    true,
		Tech3:    true,
		Tech4:    true,
	},
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	connection.ConnectDB()
	e := echo.New()

	// root statis untuk mengakses folder public
	e.Static("/public", "public") //public

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}

	// renderer
	e.Renderer = t

	// routing
	e.GET("/", home)
	e.GET("/contact", contactMe)
	e.GET("/project", myProject)
	e.GET("/project-detail/:id", projectDetail) //:id => url params
	e.POST("/add-project", addProject)
	e.GET("/delete/:id", deleteProject)     //:id => url params
	e.GET("/edit-project/:id", editProject) //:id => url params
	e.POST("/edit/:id", edit)               //:id => url params

	fmt.Println("localhost: 5000 sucssesfully")
	e.Logger.Fatal(e.Start("localhost: 5000"))
}

func home(c echo.Context) error {

	data, _ := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, technologys, description, image FROM public.tb_project;")

	var result []Projects
	for data.Next() {
		var each = Projects{}

		err := data.Scan(&each.Id, &each.Title, &each.Sdate, &each.Edate, &each.Technologys, &each.Descript, &each.Image)
		if err != nil {
			fmt.Println(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]string{"Message ": err.Error()})
		}
		// Duration
		// formatDate := "2006/01/02"

		durasi := each.Edate.Sub(each.Sdate)
		var Durations string

		if durasi.Hours()/24 < 7 {
			Durations = strconv.FormatFloat(durasi.Hours()/24, 'f', 0, 64) + " Days"
		} else if durasi.Hours()/24/7 < 4 {
			Durations = strconv.FormatFloat(durasi.Hours()/24/7, 'f', 0, 64) + " Weeks"
		} else if durasi.Hours()/24/30 < 12 {
			Durations = strconv.FormatFloat(durasi.Hours()/24/30, 'f', 0, 64) + " Months"
		} else {
			Durations = strconv.FormatFloat(durasi.Hours()/24/30/12, 'f', 0, 64) + " Years"
		}

		each.Duration = Durations

		result = append(result, each)
	}
	project := map[string]interface{}{
		"Projects": result,
	}
	return c.Render(http.StatusOK, "index.html", project)
}

func contactMe(c echo.Context) error {
	return c.Render(http.StatusOK, "contact-me.html", nil)
}

func myProject(c echo.Context) error {
	return c.Render(http.StatusOK, "myProject.html", nil)
}

func addProject(c echo.Context) error {
	title := c.FormValue("project-name")
	sDate := c.FormValue("start-date")
	eDate := c.FormValue("end-date")
	tech := c.Request().Form["check"]
	desc := c.FormValue("description")
	image := "MobileApp1.jpg"

	_, err := connection.Conn.Exec(context.Background(), "INSERT INTO tb_project (title, start_date, end_date, technologys, description, image) VALUES ($1, $2, $3, $4, $5, $6)", title, sDate, eDate, tech, desc, image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}
	return c.Redirect(http.StatusMovedPermanently, "/")
}

func projectDetail(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var ProjectDetail = Projects{}

	err := connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, technologys, description, image FROM tb_project WHERE id=$1", id).Scan(&ProjectDetail.Id, &ProjectDetail.Title, &ProjectDetail.Sdate, &ProjectDetail.Edate, &ProjectDetail.Technologys, &ProjectDetail.Descript, &ProjectDetail.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	// convert date ke string
	sDateFormat := ProjectDetail.Sdate.Format("02 January 2006")
	eDateFormat := ProjectDetail.Edate.Format("02 January 2006")

	// SET Durations
	durasi := ProjectDetail.Edate.Sub(ProjectDetail.Sdate)
	var Durations string

	if durasi.Hours()/24 < 7 {
		Durations = strconv.FormatFloat(durasi.Hours()/24, 'f', 0, 64) + " Days"
	} else if durasi.Hours()/24/7 < 4 {
		Durations = strconv.FormatFloat(durasi.Hours()/24/7, 'f', 0, 64) + " Weeks"
	} else if durasi.Hours()/24/30 < 12 {
		Durations = strconv.FormatFloat(durasi.Hours()/24/30, 'f', 0, 64) + " Months"
	} else {
		Durations = strconv.FormatFloat(durasi.Hours()/24/30/12, 'f', 0, 64) + " Years"
	}

	detailProject := map[string]interface{}{
		"Projects": ProjectDetail,
		"Duration": Durations,
		"StartD":   sDateFormat,
		"EndD":     eDateFormat,
	}
	return c.Render(http.StatusOK, "projectDetail.html", detailProject)
}

func deleteProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_project WHERE id=$1", id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}

func editProject(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	edit := Projects{}
	err := connection.Conn.QueryRow(context.Background(), "SELECT id, Title, start_date, end_date, technologys, description, image FROM tb_project WHERE id=$1;", id).Scan(&edit.Id, &edit.Title, &edit.Sdate, &edit.Edate, &edit.Technologys, &edit.Descript, &edit.Image)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	var Python, Js, React, Node bool

	for _, techno := range edit.Technologys {
		if techno == "python" {
			Python = true
		}
		if techno == "js" {
			Js = true
		}
		if techno == "react" {
			React = true
		}
		if techno == "node" {
			Node = true
		}
	}
	StartFormat := edit.Sdate.Format("2006-01-02")
	EndFormat := edit.Edate.Format("2006-01-02")

	editResult := map[string]interface{}{
		"Edit":   edit,
		"Id":     id,
		"StartD": StartFormat,
		"EndD":   EndFormat,
		"Tech1":  Python,
		"Tech2":  Js,
		"Tech3":  React,
		"Tech4":  Node,
	}

	return c.Render(http.StatusOK, "updateProject.html", editResult)
}

func edit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	title := c.FormValue("project-name")
	SDate := c.FormValue("start-date")
	EDate := c.FormValue("end-date")
	descript := c.FormValue("description")
	technologys := c.Request().Form["check"]
	image := "MobileApp1.jpg"

	sdate, _ := time.Parse("2006-01-02", SDate)
	edate, _ := time.Parse("2006-01-02", EDate)

	_, err := connection.Conn.Exec(context.Background(), "UPDATE public.tb_project SET title=$1, start_date=$2, end_date=$3, description=$4, technologys=$5, Image=$6 WHERE id=$7;", title, sdate, edate, descript, technologys, image, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message ": err.Error()})
	}

	return c.Redirect(http.StatusMovedPermanently, "/")
}
