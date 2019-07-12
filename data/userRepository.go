package data

import (
  "context"
  "time"
  "errors"
  "strings"
  "fmt"
  "golang.org/x/crypto/bcrypt"
  "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
  "github.com/smukhalov/test4/models"
  "github.com/smukhalov/test4/common"
)

var dbName string = common.AppConfig.Database // "test4db"
var userCollection = common.AppConfig.UserCollection // "users"
var mongoUrl = common.AppConfig.MongoUrl // "mongodb://localhost:27017"

func CreateUser(user *models.User) (primitive.ObjectID, error) {
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

func Login(email, password string) (*models.User, error) {

  client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
  if err != nil {
    return nil, err
  }

  ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
  defer cancel()

  err = client.Connect(ctx)
  if err != nil {
     return nil, err
  }

  collection := client.Database(dbName).Collection(userCollection)

  var user models.User
  filter := bson.D{{"email", email}}

  err = collection.FindOne(ctx, filter).Decode(&user)
  if err != nil {
    return nil, err
  }

	// Validate password
	err = bcrypt.CompareHashAndPassword(user.HashPassword, []byte(password))
	return &user, err
}
