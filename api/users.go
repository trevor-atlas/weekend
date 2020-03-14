package handler

import (
	f "github.com/fauna/faunadb-go/faunadb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	utils "github.com/trevor-atalas/weekend/api/utils"
)

type User struct {
	ID string `fauna:"id"`
	Token string `fauna:"token"`
	Name string `fauna:"name"`
}

type LoginRequest struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}

type UserCreateRequest struct {
	Name string `form:"name" json:"name" xml:"name"  binding:"required"`
	Token string `form:"token" json:"name" xml:"name"  binding:"required"`
}

func Login(c *gin.Context) {
	var json LoginRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if json.User != "trevor" || json.Password != "12345" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "you are logged in", "token": "12345"})
}

// CreateEntry - returns entry object from data send in request
func CreateUser(c *gin.Context) {
	token := utils.ExtractHeader()
	if token != "12345" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid or missing authorization token"})
		return
	}
	var json UserCreateRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := User{
		Token: json.Token,
		ID: uuid.New().String(),
		Name: json.Name,
	}

	client := utils.DBClient()

	_, err := user.DBCreate(client)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

// DBGetAllRefs - get all elements
func DBGetAllRefs(client *f.FaunaClient, token string) (refs []f.RefV, err error) {
	value, err := client.Query(
		f.Paginate(
			f.MatchTerm(
				f.Index("entries_with_token"),
				token,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	value.At(f.ObjKey("data")).Get(&refs)
	return refs, nil
}

// DBGetFromRefs - get all elements
func DBGetFromRefs(client *f.FaunaClient, refs []f.RefV) (entries []Entry, err error) {
	request := mapRefV(refs, func(ref f.RefV) interface{} {
		return f.Get(ref)
	})
	value, err := client.Query(f.Arr(request))

	if err != nil {
		return nil, err
	}

	var elements f.ArrayV
	value.Get(&elements)

	results := make([]Entry, len(elements))
	for index, element := range elements {
		var object f.ObjectV
		element.At(f.ObjKey("data")).Get(&object)
		var entry Entry
		object.Get(&entry)
		results[index] = entry
	}

	return results, nil
}

// DBGet - get existing element from database
func (entry User) DBGet(client *f.FaunaClient) (value f.Value, err error) {
	return client.Query(
		f.Get(
			f.MatchTerm(
				f.Index("entry_with_token_user_item"),
				f.Arr{entry.Token, entry.UserID, entry.ItemID},
			),
		),
	)
}

// DBCreate - create new Entry object
func (entry User) DBCreate(client *f.FaunaClient) (value f.Value, err error) {
	return client.Query(
		f.Create(
			f.Collection("user"),
			f.Obj{"data": entry},
		),
	)
}

// DBUpdate - update existing object provided in result parameter
func (entry User) DBUpdate(client *f.FaunaClient, result f.Value) (value f.Value, err error) {
	var ref f.RefV
	result.At(f.ObjKey("ref")).Get(&ref)
	return client.Query(
		f.Update(
			ref,
			f.Obj{"data": entry},
		),
	)
}

// DBCreateOrUpdate - combine DBGet, DBCreate and DBUpdate to make uperation easier
func (entry User) DBCreateOrUpdate(client *f.FaunaClient) (value f.Value, err error) {
	value, _ = entry.DBGet(client)

	if value == nil {
		value, err = entry.DBCreate(client)
	} else {
		value, err = entry.DBUpdate(client, value)
	}
	return value, err
}

// DBDelete - remove Entry object from database
func (entry User) DBDelete(client *f.FaunaClient) (value f.Value, err error) {
	result, err := entry.DBGet(client)
	if result != nil {
		var ref f.RefV
		result.At(f.ObjKey("ref")).Get(&ref)
		return client.Query(f.Delete(ref))
	}
	return result, err
}

func mapRefV(vs []f.RefV, f func(f.RefV) interface{}) []interface{} {
	vsm := make([]interface{}, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
