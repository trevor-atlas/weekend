package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	f "github.com/fauna/faunadb-go/faunadb"
)

// DBClient - create DB client using secret from env
func DBClient() *f.FaunaClient {
	clientSecret := os.Getenv("FAUNADB_SECRET_KEY")
	return f.NewFaunaClient(clientSecret)
}

// ExtractToken - extract token from headers
func ExtractToken(r *http.Request) string {
	header := GetHeader("Authorization", r)
	return strings.Split(header, " ")[1]
}

func Runtime(s string) (string, time.Time) {
	log.Println("Start:	", s)
	return s, time.Now()
}

func Track(s string, startTime time.Time) {
	endTime := time.Now()
	log.Println("End:	", s, "took", endTime.Sub(startTime))
}

func GetHeader(name string, r *http.Request) string {
	header := r.Header.Get(name)
	return header
}

func GetHeaderWithDefault(name, fallback string, r *http.Request) string {
	header := GetHeader(name, r)
	if header == "" {
		header = fallback
	}
	return header
}

func GetParam(name string, r *http.Request) string {
	param := r.URL.Query().Get(name)
	return param
}

func GetParamWithDefault(name, fallback string, r *http.Request) string {
	param := GetParam(name, r)
	if param == "" {
		param = fallback
	}
	return param
}

func Write(obj interface{}, w http.ResponseWriter) {
	j, _ := json.Marshal(obj)
	w.Write(j)
}

//func AuthenticationRequired(auths ...string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		session := sessions.Default(c)
//		user := session.Get("user")
//		if user == nil {
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "user needs to be signed in to access this service"})
//			c.Abort()
//			return
//		}
//		if len(auths) != 0 {
//			authType := session.Get("authType")
//			if authType == nil || !funk.ContainsString(auths, authType.(string)) {
//				c.JSON(http.StatusForbidden, gin.H{"error": "invalid request, restricted endpoint"})
//				c.Abort()
//				return
//			}
//		}
//		// add session verification here, like checking if the user and authType
//		// combination actually exists if necessary. Try adding caching this (redis)
//		// since this middleware might be called a lot
//		c.Next()
//	}
//}
