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

