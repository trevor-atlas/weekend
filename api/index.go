package handler

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"time"
)


func buildRouter(w http.ResponseWriter, r *http.Request) *gin.Engine {
	defer track(runtime("build router"))
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	// pass our ResponseWriter and Request into gin
	// otherwise it has no knowledge of them coming in via now serverless
	router.ServeHTTP(w, r)

	return router
}

// Handler is the function that Now calls for every request
func Handler(w http.ResponseWriter, req *http.Request) {
	defer track(runtime("handler"))
	log.Println("Request url: ", req.URL.Path)

	r := buildRouter(w, req)
	r.GET("/greet", greet)

	// By default it serves on :8080 unless a
	// PORT environment variable was defined.
	r.Run()
	// router.Run(":3000") for a hard coded port
}

func greet(c *gin.Context) {
	defer track(runtime("greet"))
	name := c.DefaultQuery("name", "Guest")
	c.String(http.StatusOK,"<h1>Hello, "+name+" from Go on Now!</h1>")
}

func runtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}
