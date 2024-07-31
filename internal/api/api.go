package api

import (
	"github/diegoHDCz/gopet/internal/api/spec"
	"github/diegoHDCz/gopet/internal/mongostore"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type store interface{}

type PetAPI struct {
	store  store
	logger *zap.Logger
}

func NewPetAPI(store *mongo.Database, logger *zap.Logger) PetAPI {
	return PetAPI{mongostore.New(store), logger}
}

// Add a new pet to the store
// (POST /pet)
func (api *PetAPI) AddPet(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// Update an existing pet
// (PUT /pet)
func (api *PetAPI) UpdatePet(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// Finds Pets by status
// (GET /pet/findByStatus)
func (api *PetAPI) FindPetsByStatus(c *gin.Context, params spec.FindPetsByStatusParams) {
	panic("not implemented") // TODO: Implement
}

// Finds Pets by tags
// (GET /pet/findByTags)
func (api *PetAPI) FindPetsByTags(c *gin.Context, params spec.FindPetsByTagsParams) {
	panic("not implemented") // TODO: Implement
}

// Deletes a pet
// (DELETE /pet/{petId})
func (api *PetAPI) DeletePet(c *gin.Context, petId int64, params spec.DeletePetParams) {
	panic("not implemented") // TODO: Implement
}

// Find pet by ID
// (GET /pet/{petId})
func (api *PetAPI) GetPetById(c *gin.Context, petId int64) {
	panic("not implemented") // TODO: Implement
}

// Updates a pet in the store with form data
// (POST /pet/{petId})
func (api *PetAPI) UpdatePetWithForm(c *gin.Context, petId int64, params spec.UpdatePetWithFormParams) {
	panic("not implemented") // TODO: Implement
}

// uploads an image
// (POST /pet/{petId}/uploadImage)
func (api *PetAPI) UploadFile(c *gin.Context, petId int64, params spec.UploadFileParams) {
	panic("not implemented") // TODO: Implement
}

// Returns pet inventories by status
// (GET /store/inventory)
func (api *PetAPI) GetInventory(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// Place an order for a pet
// (POST /store/order)
func (api *PetAPI) PlaceOrder(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// Delete purchase order by ID
// (DELETE /store/order/{orderId})
func (api *PetAPI) DeleteOrder(c *gin.Context, orderId int64) {
	panic("not implemented") // TODO: Implement
}

// Find purchase order by ID
// (GET /store/order/{orderId})
func (api *PetAPI) GetOrderById(c *gin.Context, orderId int64) {
	panic("not implemented") // TODO: Implement
}

// Create user
// (POST /user)
func (api *PetAPI) CreateUser(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// Creates list of users with given input array
// (POST /user/createWithList)
func (api *PetAPI) CreateUsersWithListInput(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// Logs user into the system
// (GET /user/login)
func (api *PetAPI) LoginUser(c *gin.Context, params spec.LoginUserParams) {
	panic("not implemented") // TODO: Implement
}

// Logs out current logged in user session
// (GET /user/logout)
func (api *PetAPI) LogoutUser(c *gin.Context) {
	panic("not implemented") // TODO: Implement
}

// Delete user
// (DELETE /user/{username})
func (api *PetAPI) DeleteUser(c *gin.Context, username string) {
	panic("not implemented") // TODO: Implement
}

// Get user by user name
// (GET /user/{username})
func (api *PetAPI) GetUserByName(c *gin.Context, username string) {
	panic("not implemented") // TODO: Implement
}

// Update user
// (PUT /user/{username})
func (api *PetAPI) UpdateUser(c *gin.Context, username string) {
	panic("not implemented") // TODO: Implement
}
