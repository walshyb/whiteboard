package main

import (
	"context"
	"errors"

	"github.com/google/uuid"
  "go.mongodb.org/mongo-driver/v2/mongo"
  "go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	 events "github.com/walshyb/whiteboard/proto"
)

func EnsureDemoBoard(ctx context.Context, col *mongo.Collection) (*mongo.UpdateResult, error) {
  filter := bson.M{"_id": "demo"}

  update := bson.M{
    "$setOnInsert": bson.M{
      "_id": "demo",
			"shapes": []interface{}{},
    },
  }

	opts := options.UpdateOne().SetUpsert(true)

  return col.UpdateOne(ctx, filter, update, opts)
}

func AddShape(hub *Hub, shape *events.Shape ) (*events.Shape, error) {
	// Validate that passed in shape type exists
	if _, ok := events.ShapeType_value[shape.Type.String()]; ok {
		return shape, errors.New("Shape value does not exist")
	}

  filter := bson.M{"_id": "demo"}
	shapeId := uuid.NewString()
	shape.Id = shapeId

  update := bson.M{
    "$push": bson.M{
      "shapes": shape,
    },
  }

  collection := hub.mongo.Database("whiteboards").Collection("demo")
  _, err := collection.UpdateOne(hub.ctx, filter, update)
  return shape, err
}

func (hub *Hub) GetBoard() *mongo.SingleResult {
  collection := hub.mongo.Database("whiteboards").Collection("demo")
	res := collection.FindOne(hub.ctx, bson.D{{ Key: "_id", Value: "demo" }})
  return res
}
