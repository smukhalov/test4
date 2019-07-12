package data

import (
  "context"
  "time"
  "errors"
  "strings"
  "fmt"
  "golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
  "github.com/smukhalov/test4/models"
  "github.com/smukhalov/test4/common"
)

func CreateUser(user *models.User) (primitive.ObjectID, error) {
  var dbName string = common.AppConfig.Database // "test4db"
  var userCollection = common.AppConfig.UserCollection // "users"
  var mongoUrl = common.AppConfig.MongoUrl // "mongodb://localhost:27017"

  userid := primitive.NewObjectID()

  hpass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return userid, err
	}
	user.HashPassword = hpass
	user.Password = ""

  client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
  if err != nil {
    return userid, err
  }

  ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
  defer cancel()

  err = client.Connect(ctx)
  if err != nil {
     return userid, err
  }

  collection := client.Database(dbName).Collection(userCollection)

  insertResult, err := collection.InsertOne(ctx, user)
  if err != nil {
    if strings.Contains(err.Error(),"duplicate key error") {
      return userid, errors.New(fmt.Sprintf("Пользователь %s уже существует", user.Email))
    }

    return userid, err
  }

  userid, ok := insertResult.InsertedID.(primitive.ObjectID)

  if !ok {
    return userid, errors.New("Ошибка при получении insertResult.InsertedID.(objectid.ObjectID)")
  }

  return userid, nil
}
