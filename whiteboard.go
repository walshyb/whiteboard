package main

import (
  "time"
	"context"

	"github.com/google/uuid"
  "go.mongodb.org/mongo-driver/v2/mongo"
  "go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Canvas struct {
  ID        string     `bson:"_id"`
  Name      string     `bson:"name"`
  CreatedAt time.Time  `bson:"createdAt"`
  UpdatedAt time.Time  `bson:"updatedAt"`
  Shapes    []Shape    `bson:"shapes"`
}

type Shape struct {
  Id          string        `bson:"id"`
  Type        string        `bson:"type"` // rect, ellipse
  X           float64       `bson:"x,omitempty"`
  Y           float64       `bson:"y,omitempty"`
  Width       float64       `bson:"width,omitempty"`
  Height      float64       `bson:"height,omitempty"`
  RX          float64       `bson:"rx,omitempty"`
  RY          float64       `bson:"ry,omitempty"`
  Color       string        `bson:"color,omitempty"`
  StrokeWidth float64       `bson:"strokeWidth,omitempty"`
}

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

func AddShape(ctx context.Context, col *mongo.Collection, shape Shape) error {
  filter := bson.M{"_id": "demo"}
	shapeId := uuid.NewString()
	shape.Id = shapeId

  update := bson.M{
    "$push": bson.M{
      "shapes": shape,
    },
  }

  _, err := col.UpdateOne(ctx, filter, update)
  return err
}

func (hub *Hub) GetCanvas() *mongo.SingleResult {
  collection := hub.mongo.Database("whiteboards").Collection("demo")
  res := collection.FindOne(hub.ctx, bson.D{{ "name", "demo" }})
  return res
}
