package main

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	events "github.com/walshyb/whiteboard/proto"


)

type Shape struct {
	Id     string  `bson:"id" json:"id"`
	Type   int32   `bson:"type" json:"type"`
	Width  float64 `bson:"width" json:"width"`
	Height float64 `bson:"height" json:"height"`
	X      float64 `bson:"x" json:"x"`
	Y      float64 `bson:"y" json:"y"`
	Color  string  `bson:"color" json:"color"`
}

type Board struct {
	Id     string  `bson:"_id" json:"id"` 
	Shapes []Shape `bson:"shapes" json:"shapes"`
	Name   string  `bson:"name" json:"name"`
}

func EnsureDemoBoard(ctx context.Context, col *mongo.Collection) (*mongo.UpdateResult, error) {
	filter := bson.M{"_id": "demo"}

	update := bson.M{
		"$setOnInsert": bson.M{
			"_id": "demo",
			"name": "Demo",
			"shapes": []interface{}{},
		},
	}

	opts := options.UpdateOne().SetUpsert(true)

	return col.UpdateOne(ctx, filter, update, opts)
}

func AddShape(hub *Hub, shape *events.Shape ) (*events.Shape, error) {
	// Validate that passed in shape type exists
	if _, ok := events.ShapeType_value[shape.Type.String()]; !ok {
		println("Shape error")
		return shape, errors.New("Shape value does not exist")
	}

	shape.Id = uuid.NewString()

	dbShape := Shape{
		Id:     shape.Id,
		Type:   int32(shape.Type), // Cast Proto Enum to int32
		X:      shape.X,
		Y:      shape.Y,
		Width:  shape.Width,
		Height: shape.Height,
		Color:  shape.Color,
	}

	collection := hub.mongo.Database("whiteboards").Collection("demo")
	filter := bson.M{"_id": "demo"}
	update := bson.M{"$push": bson.M{"shapes": dbShape}}
	_, err := collection.UpdateOne(hub.ctx, filter, update)

	return shape, err
}

func (hub *Hub) GetBoard() (*Board, error) {
	collection := hub.mongo.Database("whiteboards").Collection("demo")

	var board Board

	err := collection.FindOne(hub.ctx, bson.D{{Key: "_id", Value: "demo"}}).Decode(&board)
	if err != nil {
		// Handle "no document found" if necessary
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	return &board, nil
}
