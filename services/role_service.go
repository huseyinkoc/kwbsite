package services

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var rolesCollection *mongo.Collection

func InitRolesService(client *mongo.Client) {
	rolesCollection = client.Database("admin_panel").Collection("roles")
}

// GetRolePermissions fetches permissions for a specific role and module
func GetRolePermissions(ctx context.Context, role string, module string) ([]string, error) {
	var roleData struct {
		Permissions map[string][]string `bson:"permissions"`
	}

	filter := bson.M{"_id": role}
	err := rolesCollection.FindOne(ctx, filter).Decode(&roleData)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	permissions, exists := roleData.Permissions[module]
	if !exists {
		return nil, nil
	}

	return permissions, nil
}
