package services

import (
	"admin-panel/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var rolesCollection *mongo.Collection

// InitRolesService initializes the roles collection
func InitRolesService(client *mongo.Client) {
	rolesCollection = client.Database("admin_panel").Collection("roles")
}

// GetModulePermissions fetches permissions for a specific role and module
func GetModulePermissions(ctx context.Context, role string, module string) ([]string, error) {
	var roleData models.Role
	filter := bson.M{"_id": role}

	err := rolesCollection.FindOne(ctx, filter).Decode(&roleData)
	if err != nil {
		return nil, err
	}

	permissions, exists := roleData.Modules[module]
	if !exists {
		return nil, nil
	}

	return permissions, nil
}
