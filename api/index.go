package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func buildRouter(w http.ResponseWriter, r *http.Request) *gin.Engine {
	defer track(runtime("build router"))
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())
	// otherwise it has no knowledge of them coming in via now serverless
	router.ServeHTTP(w, r)

	return router
}

// Handler is the function that Now calls for every request
func Handler(w http.ResponseWriter, r *http.Request) {

	defer track(runtime("handler"))
	log.Println("Request url: ", r.URL.Path)
	defer track(runtime("greet"))
	keys, ok := r.URL.Query()["name"]

	name := "Guest"

	if ok || len(keys[0]) > 1 {
		log.Println("Url Param 'name' is missing")
		name = keys[0]
		return
	}

	fmt.Fprintf(w, `<h1>Hello, %s from Go on Now!</h1>`, name)

}

func greet(c *gin.Context) {
	defer track(runtime("greet"))
	name := c.DefaultQuery("name", "Guest")
	c.String(http.StatusOK, "<h1>Hello, "+name+" from Go on Now!</h1>")
}

func runtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}
