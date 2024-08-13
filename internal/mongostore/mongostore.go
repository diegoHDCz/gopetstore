package mongostore

import (
	"context"
	"errors"
	"fmt"
	"github/diegoHDCz/gopet/internal/api/spec"
	"github/diegoHDCz/gopet/internal/utils"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	db *mongo.Database
}

func New(db *mongo.Database) *Store {
	return &Store{db: db}
}

func (s *Store) SavePet(body *spec.StorePet) (*primitive.ObjectID, error) {
	coll := s.db.Collection("pets")

	result, err := coll.InsertOne(context.TODO(), body)
	if err != nil {
		log.Fatal("Error inserting data: %w", err)
		return nil, err
	}
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Fatal("Error parsing id: %w", err)
		return nil, err
	}

	return &oid, nil
}

func (s *Store) GetTagsByName(tagNames *[]string) (*[]spec.Tag, error) {
	collection := s.db.Collection("tags")

	ctx := context.TODO()

	var tags []spec.Tag
	fmt.Println(tagNames)

	filter := bson.M{"name": bson.M{"$in": tagNames}}

	cursor, err := collection.Find(context.TODO(), filter, options.Find())

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var result spec.Tag
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		tags = append(tags, result)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return &tags, nil
}

func (s *Store) UpdatePet(body *spec.Pet) (*spec.Pet, error) {
	collection := s.db.Collection("pets")

	ctx := context.TODO()

	id, _ := utils.StringToID(*body.Id)

	filter := bson.D{{"_id", id}}

	update := bson.D{{"$set", bson.D{{"name", body.Name},
		{"photoUrls", body.PhotoUrls}, {"category", body.Category}, {"status", body.Status}}}}

	result, err := collection.UpdateOne(ctx, filter, update)

	fmt.Printf("Documents matched: %v\n", result.MatchedCount)
	fmt.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		return nil, err
	}
	if result.ModifiedCount > 0 {
		return body, nil
	}
	return nil, errors.New("something went wrong")
}

func (s *Store) DeleteDocumentById(petId *primitive.ObjectID) error {
	collection := s.db.Collection("pets")

	ctx := context.TODO()

	filter := bson.M{"_id": petId}

	doc, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		return err
	}

	if doc.DeletedCount <= 0 {
		return errors.New("Something went wrong")
	}
	return nil

}

func (s *Store) GetOnePetById(petId *primitive.ObjectID) (*spec.Pet, error) {
	collection := s.db.Collection("pets")

	ctx := context.TODO()

	filter := bson.M{"_id": petId}

	var r *spec.Pet
	err := collection.FindOne(ctx, filter).Decode(&r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (s *Store) FindPetsByTagsId(queryIDs *[]string) (*[]spec.Pet, error) {
	collection := s.db.Collection("pets")

	ctx := context.TODO()
	fmt.Println("Befor filter")
	petFilter := bson.M{"tags._id": bson.M{"$in": queryIDs}}
	petCursor, err := collection.Find(ctx, petFilter)

	if err != nil {
		return nil, err
	}
	defer petCursor.Close(ctx)

	var petsArr []spec.Pet
	for petCursor.Next(ctx) {
		var result spec.Pet
		if err := petCursor.Decode(&result); err != nil {
			return nil, err
		}
		petsArr = append(petsArr, result)
	}

	return &petsArr, nil
}
