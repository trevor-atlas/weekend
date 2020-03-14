package utils

import (
	"github.com/gin-gonic/gin"
	"os"
	"strings"

	f "github.com/fauna/faunadb-go/faunadb"
)

// DBClient - create DB client using secret from env
func DBClient() *f.FaunaClient {
	clientSecret := os.Getenv("FAUNADB_SECRET_KEY")
	return f.NewFaunaClient(clientSecret)
}

// ExtractToken - extract token from headers
func ExtractToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	return strings.Split(header, " ")[1]
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
