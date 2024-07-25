package controllers

import (
	"TestProject/model"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetComputers returns pc with requested parameters
func GetComputers(c *fiber.Ctx) error {
	filter := bson.M{}

	if findId := c.Query("id", ""); findId != "" {
		filter["id"] = findId
	}
	if findName := c.Query("name", ""); findName != "" {
		filter["name"] = primitive.Regex{Pattern: findName, Options: ""}
	}
	if findPrice := c.Query("price", ""); findPrice != "" {
		priceFloat, err := strconv.ParseFloat(findPrice, 8)

		if err != nil {
			c.SendStatus(fiber.StatusBadRequest)
		}
		filter["price"] = priceFloat
	}
	if findRating := c.Query("rating", ""); findRating != "" {
		filter["rating"] = findRating
	}
	if findSpecifications := c.Query("specs", ""); findSpecifications != "" {
		specifications := make(map[string]interface{})
		err := json.Unmarshal([]byte(findSpecifications), &specifications)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Incorrect specs parameters")
		}

		for key, value := range specifications {
			filterKey := fmt.Sprintf("specifications.%s", key)
			if reflect.TypeOf(value) == reflect.TypeOf("") {
				filter[filterKey] = primitive.Regex{Pattern: value.(string), Options: ""}
			} else {
				filter[filterKey] = value
			}
		}
	}

	cursor, err := model.Database.Database.Collection("computers").Find(c.Context(), filter)

	if err != nil {
		return err
	}

	var computers = make([]model.Computer, 0)

	if err := cursor.All(c.Context(), &computers); err != nil {
		return err
	}

	return c.JSON(computers)
}

func PostComputer(c *fiber.Ctx) error {
	newComputer := new(model.Computer)

	if err := c.BodyParser(newComputer); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	newComputer.ID = ""

	insertionResult, err := model.Database.Database.Collection("computers").InsertOne(c.Context(), newComputer)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(201).JSON(insertionResult)
}

func PutComputer(c *fiber.Ctx) error {
	idParam := c.Params("id")
	computerId, err := primitive.ObjectIDFromHex(idParam)

	// the provided ID might be invalid ObjectID
	if err != nil {
		return c.SendStatus(400)
	}

	computer := new(model.Computer)

	if err := c.BodyParser(computer); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	query := bson.D{{Key: "_id", Value: computerId}}
	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "name", Value: computer.Name},
				{Key: "price", Value: computer.Price},
				{Key: "rating", Value: computer.Rating},
				{Key: "specifications", Value: computer.Specifications},
			},
		},
	}

	err = model.Database.Database.Collection("computers").FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.SendStatus(404)
		}

		return c.SendStatus(500)
	}

	computer.ID = idParam
	return c.Status(200).JSON(computer)
}

func DeleteComputer(c *fiber.Ctx) error {
	computerId, err := primitive.ObjectIDFromHex(c.Params("id"))

	if err != nil {
		return c.SendStatus(400)
	}

	query := bson.D{{Key: "_id", Value: computerId}}
	result, err := model.Database.Database.Collection("computers").DeleteOne(c.Context(), &query)

	if err != nil {
		return c.SendStatus(500)
	}

	if result.DeletedCount < 1 {
		return c.SendStatus(404)
	}

	return c.SendStatus(204)
}
