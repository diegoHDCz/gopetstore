package api

import (
	"errors"
	"github/diegoHDCz/gopet/internal/api/spec"
	"github/diegoHDCz/gopet/internal/mongostore"
	"github/diegoHDCz/gopet/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type store interface {
	SavePet(body *spec.StorePet) (*primitive.ObjectID, error)
	GetTagsByName(tagNames *[]string) (*[]spec.Tag, error)
	UpdatePet(body *spec.Pet) (*spec.Pet, error)
	DeleteDocumentById(petId *primitive.ObjectID) error
	GetOnePetById(petId *primitive.ObjectID) (*spec.Pet, error)
	FindPetsByTagsId(queryIDs *[]string) (*[]spec.Pet, error)
}

type PetAPI struct {
	store  store
	logger *zap.Logger
}

func NewPetAPI(db *mongo.Database, logger *zap.Logger) PetAPI {
	return PetAPI{mongostore.New(db), logger}
}

// Add a new pet to the store
// (POST /pet)
func (api PetAPI) AddPet(c *gin.Context) {
	internalServerMessage := "sometthing went wrong"
	badRequestMessage := "request is not valid"
	var body spec.StorePet

	if err := c.ShouldBindJSON(&body); err != nil {
		api.logger.Error("error parsing json ", zap.Error(err))
		c.JSON(http.StatusBadRequest, spec.Error{Message: &badRequestMessage})
		return
	}
	api.logger.Debug("body received parsed ", zap.Any("json", body))

	tagsNames := make([]string, len(*body.Tags))

	for i, tag := range *body.Tags {
		api.logger.Debug("tag", zap.Any("json", tag))
		tagsNames[i] = *tag.Name
	}

	api.logger.Debug("tag ", zap.Any("json", tagsNames))

	body.Tags, _ = api.store.GetTagsByName(&tagsNames)

	result, err := api.store.SavePet(&body)

	if err != nil {
		api.logger.Error("error parsing json ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, spec.Error{Message: &internalServerMessage})
		return
	}
	strId := result.Hex()
	rs := spec.Pet{
		Id:        &strId,
		Category:  body.Category,
		Name:      *body.Category,
		Tags:      body.Tags,
		PhotoUrls: body.PhotoUrls,
		Status:    (*spec.PetStatus)(body.Status),
	}

	c.JSON(http.StatusCreated, &rs)

}

// Update an existing pet
// (PUT /pet)
func (api PetAPI) UpdatePet(c *gin.Context) {
	badRequestMessage := "request is not valid"
	internalServerMessage := "sometthing went wrong"
	documentNotFound := "document not found"
	var body spec.Pet
	if err := c.ShouldBindJSON(&body); err != nil {
		api.logger.Error("error parsing json ", zap.Error(err))
		c.JSON(http.StatusBadRequest, spec.Error{Message: &badRequestMessage})
		return
	}

	api.logger.Debug("body parsed", zap.Any("json", body))

	res, err := api.store.UpdatePet(&body)
	if err != nil {
		api.logger.Error("error parsing json ", zap.Error(err))
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusInternalServerError, spec.Error{Message: &documentNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, spec.Error{Message: &internalServerMessage})
		return
	}

	c.JSON(http.StatusOK, &res)
}

// Finds Pets by tags
// (GET /pet/findByTags)
func (api PetAPI) FindPetsByTags(c *gin.Context, params spec.FindPetsByTagsParams) {
	tagsParams := c.QueryArray("tags")

	notFoundMessage := "not found"
	internalServerError := "something went wrong"

	api.logger.Debug("olha ", zap.Any("query", tagsParams))

	tags, err := api.store.GetTagsByName(&tagsParams)

	if err != nil {
		api.logger.Error("Error retrieving data ", zap.Error(err))
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, spec.Error{Message: &notFoundMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, spec.Error{Message: &internalServerError})
	}

	queryIDs := make([]string, len(*tags))

	for i, docTag := range *tags {

		if docTag.Id == nil {
			api.logger.Error("Error retrieving data ", zap.String("error", "ID is null"))
		}

		queryIDs[i] = *docTag.Id

	}
	api.logger.Debug("queries", zap.Any("my ids", &tags))
	res, err := api.store.FindPetsByTagsId(&queryIDs)
	if err != nil {
		api.logger.Error("Error retrieving data ", zap.Error(err))
		c.JSON(http.StatusInternalServerError, spec.Error{Message: &internalServerError})
		return
	}

	c.JSON(http.StatusOK, &res)
}

// Deletes a pet
// (DELETE /pet/{petId})
func (api PetAPI) DeletePet(c *gin.Context, petId spec.ObjectId) {
	internalServerError := "something went wrong"
	notFoundMessage := "not found"

	oid, err := utils.StringToID(petId)

	api.logger.Debug("pet id ", zap.Any("query", oid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, spec.Error{Message: &internalServerError})
		return
	}
	if err := api.store.DeleteDocumentById(oid); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, spec.Error{Message: &notFoundMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, spec.Error{Message: &internalServerError})
		return
	}

	c.JSON(http.StatusNoContent, nil)

}

// Find pet by ID
// (GET /pet/{petId})
func (api PetAPI) GetPetById(c *gin.Context, petId spec.ObjectId) {
	internalServerError := "something went wrong"
	notFoundMessage := "not found"

	api.logger.Debug("pet id ", zap.Any("query", petId))

	oid, err := utils.StringToID(petId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, spec.Error{Message: &internalServerError})
		return
	}
	res, err := api.store.GetOnePetById(oid)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, spec.Error{Message: &notFoundMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, spec.Error{Message: &internalServerError})
		return
	}
	c.JSON(http.StatusOK, &res)
}
