package utils

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StringToID(mongoStr string) (*primitive.ObjectID, error) {

	objID, err := primitive.ObjectIDFromHex(mongoStr)
	if err != nil {
		return nil, errors.New("could not parse string to mongodb id")
	}

	return &objID, nil

}
