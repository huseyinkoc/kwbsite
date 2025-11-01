package services

import (
	"admin-panel/models"
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection

// BcryptCost okunur veya default 12
var bcryptCost = func() int {
	if v := os.Getenv("BCRYPT_COST"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i >= 4 {
			return i
		}
	}
	return 12
}()

// InitUserService initializes the user collection
func InitUserService(client *mongo.Client) {
	userCollection = client.Database("admin_panel").Collection("users")
}

// CreateUser inserts a new user into the database
func CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID()
	return userCollection.InsertOne(ctx, user)
}

// GetAllUsers retrieves all users from the database (excludes password)
func GetAllUsers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Find().SetProjection(bson.M{"password": 0})
	cursor, err := userCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUser updates a user in the database
func UpdateUser(id primitive.ObjectID, update bson.M) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	return userCollection.UpdateOne(ctx, filter, bson.M{"$set": update})
}

// DeleteUser deletes a user by ID
func DeleteUser(id primitive.ObjectID) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	return userCollection.DeleteOne(ctx, filter)
}

// HashPassword hashes a plain text password with configurable cost
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	return string(bytes), err
}

// CheckPassword compares a hashed password with a plain text password
func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GetUserByID returns user without password projection
func GetUserByID(id primitive.ObjectID) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.FindOne().SetProjection(bson.M{"password": 0})
	err := userCollection.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments
		}
		return user, err
	}
	return user, nil
}

func GetUserByUsername(username string) (models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.FindOne().SetProjection(bson.M{"password": 0})
	err := userCollection.FindOne(ctx, bson.M{"username": username}, opts).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments
		}
		return user, err
	}

	return user, nil
}

func GetUserByEmail(email string) (models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parola hariç dönen versiyon (diğer api'lar için güvenli)
	opts := options.FindOne().SetProjection(bson.M{"password": 0})
	err := userCollection.FindOne(ctx, bson.M{"email": email}, opts).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments
		}
		return user, err
	}

	return user, nil
}

// Yeni: login/auth için parola dahil dönen fonksiyonlar
func GetUserByEmailWithPassword(email string) (models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parola dahil döndür (sadece auth akışlarında kullanılmalı)
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments
		}
		return user, err
	}

	return user, nil
}

func GetUserByUsernameWithPassword(username string) (models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments
		}
		return user, err
	}

	return user, nil
}

func GetUserByPhoneWithPassword(phone string) (models.User, error) {
	var user models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := userCollection.FindOne(ctx, bson.M{"phone_number": phone}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments
		}
		return user, err
	}

	return user, nil
}

func GetUserByPhone(phone string) (models.User, error) {
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.FindOne().SetProjection(bson.M{"password": 0})
	err := userCollection.FindOne(ctx, bson.M{"phone_number": phone}, opts).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return user, mongo.ErrNoDocuments
		}
		return user, err
	}
	return user, nil
}

// VerifyUserAccount sets the is_verified field to true for a specific user
func VerifyUserAccount(ctx context.Context, userID primitive.ObjectID) error {
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"is_verified": true}}

	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUserEmailByID retrieves the email address of a user by their ID
func GetUserEmailByID(ctx context.Context, userID primitive.ObjectID) (string, error) {
	var user struct {
		Email string `bson:"email"`
	}

	// Kullanıcıyı veritabanında bul
	err := userCollection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	return user.Email, nil
}

// UpdateUserPassword updates the password of a user by their ID
func UpdateUserPassword(ctx context.Context, userID primitive.ObjectID, newPassword string) error {
	// Şifreyi hashle
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcryptCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// Kullanıcının şifresini güncelle
	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"password": string(hashedPassword)}}

	result, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

// GetUserIDByEmail retrieves the user ID for a given email address
func GetUserIDByEmail(ctx context.Context, email string) (primitive.ObjectID, error) {
	var user struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	// Kullanıcıyı email adresine göre bul
	err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return primitive.NilObjectID, errors.New("user not found")
		}
		return primitive.NilObjectID, err
	}

	return user.ID, nil
}

func UpdateUserPreferredLanguage(userID string, languageCode string) error {
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	_, err = userCollection.UpdateOne(
		context.Background(),
		bson.M{"_id": objectID},
		bson.M{"$set": bson.M{"preferred_language": languageCode, "updated_at": primitive.NewDateTimeFromTime(time.Now())}},
	)
	return err
}

func IsLanguageEnabled(languageCode string) (bool, error) {
	var lang models.Language
	err := languageCollection.FindOne(context.Background(), bson.M{"code": languageCode, "enabled": true}).Decode(&lang)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
