package data

import (
  "context"
  "time"
  "errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
  "test4/models"
  "test4/common"
)

func CreateUser(user *models.User) (primitive.ObjectID, error) {
  var dbName string = common.AppConfig.Database // "test4db"
  var userCollection = common.AppConfig.UserCollection // "users"
  var mongoUrl = common.AppConfig.MongoUrl // "mongodb://localhost:27017"

  userid := primitive.NewObjectID()

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
    return userid, err
  }

  userid, ok := insertResult.InsertedID.(primitive.ObjectID)

  if !ok {
    return userid, errors.New("Ошибка при получении insertResult.InsertedID.(objectid.ObjectID)")
  }

  return userid, nil
}
