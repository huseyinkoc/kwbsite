package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"admin-panel/configs"
	"admin-panel/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v5"
)

var refreshCollection *mongo.Collection

// InitAuthService initializes collections used by auth service
func InitAuthService(client *mongo.Client) {
	refreshCollection = client.Database("admin_panel").Collection("refresh_tokens")
}

// getJWTSecret returns JWT signing key from env
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: JWT_SECRET not set, using insecure fallback (development only)")
		secret = "dev_fallback_secret"
	}
	return secret
}

// GenerateAccessToken creates a short-lived JWT access token
func GenerateAccessToken(user models.User) (string, time.Time, error) {
	exp := time.Now().Add(15 * time.Minute)
	claims := jwt.MapClaims{
		"sub":  user.ID.Hex(),
		"iat":  time.Now().Unix(),
		"exp":  exp.Unix(),
		"role": user.Roles,
		"aud":  "admin-api",
		"iss":  "aystek",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(configs.GetJWTSecret())
	signed, err := token.SignedString(secret)
	return signed, exp, err
}

// GenerateAndStoreRefreshToken creates a refresh token, stores hashed version in Mongo and returns plaintext token (id:value)
func GenerateAndStoreRefreshToken(userID primitive.ObjectID) (string, time.Time, error) {
	if refreshCollection == nil {
		return "", time.Time{}, errors.New("refresh collection not initialized")
	}

	// token id
	tokenID := primitive.NewObjectID().Hex()

	// token value
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", time.Time{}, err
	}
	tokenValue := base64.RawURLEncoding.EncodeToString(b)

	// hash token value (bcrypt)
	hash, err := bcrypt.GenerateFromPassword([]byte(tokenValue), bcryptCost)
	if err != nil {
		return "", time.Time{}, err
	}

	exp := time.Now().Add(7 * 24 * time.Hour) // 7 days

	doc := bson.M{
		"_id":        tokenID,
		"user_id":    userID,
		"token_hash": string(hash),
		"created_at": time.Now(),
		"expires_at": exp,
		"revoked":    false,
	}

	if _, err := refreshCollection.InsertOne(context.Background(), doc); err != nil {
		return "", time.Time{}, err
	}

	plain := fmt.Sprintf("%s:%s", tokenID, tokenValue)
	return plain, exp, nil
}

// VerifyAndRotateRefreshToken verifies provided refresh token, rotates it (single-use) and returns new plaintext token + userID
func VerifyAndRotateRefreshToken(plain string) (string, primitive.ObjectID, error) {
	if refreshCollection == nil {
		return "", primitive.NilObjectID, errors.New("refresh collection not initialized")
	}
	parts := []byte(plain)
	idx := 0
	for i, b := range parts {
		if b == ':' {
			idx = i
			break
		}
	}
	if idx == 0 {
		return "", primitive.NilObjectID, errors.New("invalid token format")
	}
	tokenID := string(parts[:idx])
	tokenValue := string(parts[idx+1:])

	// find token document
	var doc struct {
		ID        string             `bson:"_id"`
		UserID    primitive.ObjectID `bson:"user_id"`
		TokenHash string             `bson:"token_hash"`
		ExpiresAt time.Time          `bson:"expires_at"`
		Revoked   bool               `bson:"revoked"`
	}
	err := refreshCollection.FindOne(context.Background(), bson.M{"_id": tokenID}).Decode(&doc)
	if err != nil {
		return "", primitive.NilObjectID, errors.New("invalid token")
	}
	if doc.Revoked {
		return "", primitive.NilObjectID, errors.New("token revoked")
	}
	if time.Now().After(doc.ExpiresAt) {
		return "", primitive.NilObjectID, errors.New("token expired")
	}

	// compare hash
	if err := bcrypt.CompareHashAndPassword([]byte(doc.TokenHash), []byte(tokenValue)); err != nil {
		return "", primitive.NilObjectID, errors.New("invalid token")
	}

	// rotate: create new token and mark old revoked
	newPlain, _, err := GenerateAndStoreRefreshToken(doc.UserID)
	if err != nil {
		return "", primitive.NilObjectID, err
	}

	// revoke old token, set replaced_by
	_, _ = refreshCollection.UpdateOne(context.Background(), bson.M{"_id": tokenID}, bson.M{
		"$set": bson.M{
			"revoked":     true,
			"revoked_at":  time.Now(),
			"replaced_by": stringsOrEmpty(newPlain), // small helper to persist replacement id/value if needed
		},
	})

	return newPlain, doc.UserID, nil
}

func stringsOrEmpty(s string) string {
	if s == "" {
		return ""
	}
	return s
}

// RevokeRefreshTokenByID revokes a refresh token by its tokenID
func RevokeRefreshTokenByID(tokenID string) error {
	if refreshCollection == nil {
		return errors.New("refresh collection not initialized")
	}
	_, err := refreshCollection.UpdateOne(context.Background(), bson.M{"_id": tokenID}, bson.M{
		"$set": bson.M{"revoked": true, "revoked_at": time.Now()},
	})
	return err
}

// RevokeAllRefreshTokensForUser revokes all tokens for given user
func RevokeAllRefreshTokensForUser(userID primitive.ObjectID) error {
	if refreshCollection == nil {
		return errors.New("refresh collection not initialized")
	}
	_, err := refreshCollection.UpdateMany(context.Background(), bson.M{"user_id": userID}, bson.M{
		"$set": bson.M{"revoked": true, "revoked_at": time.Now()},
	})
	return err
}

// Failed login / lockout helpers

// IncrementFailedLoginByEmail increments failed_attempts and sets lockout if threshold reached
func IncrementFailedLoginByEmail(email string) (bool, error) {
	// Threshold and duration
	const threshold = 5
	lockDuration := 15 * time.Minute

	// increment and get updated doc
	filter := bson.M{"email": email}
	update := bson.M{
		"$inc": bson.M{"failed_attempts": 1},
	}
	opts := optionsAfter()
	var res struct {
		FailedAttempts int                `bson:"failed_attempts"`
		ID             primitive.ObjectID `bson:"_id"`
	}
	err := userCollection.FindOneAndUpdate(context.Background(), filter, update, opts).Decode(&res)
	if err != nil {
		return false, err
	}

	if res.FailedAttempts >= threshold {
		_, err := userCollection.UpdateOne(context.Background(), bson.M{"_id": res.ID}, bson.M{
			"$set": bson.M{"locked_until": time.Now().Add(lockDuration)},
		})
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

// ResetFailedAttempts resets failed_attempts and locked_until for user
func ResetFailedAttempts(userID primitive.ObjectID) error {
	_, err := userCollection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{
		"$set":   bson.M{"failed_attempts": 0},
		"$unset": bson.M{"locked_until": ""},
	})
	return err
}

// IsAccountLockedByEmail checks if account is locked
func IsAccountLockedByEmail(email string) (bool, time.Time, error) {
	var res struct {
		LockedUntil *time.Time `bson:"locked_until"`
	}
	err := userCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&res)
	if err != nil {
		return false, time.Time{}, err
	}
	if res.LockedUntil != nil && time.Now().Before(*res.LockedUntil) {
		return true, *res.LockedUntil, nil
	}
	return false, time.Time{}, nil
}

// IsRefreshTokenValid checks the provided plain refresh token (id:value) against DB without rotating.
// Returns userID, valid, error.
func IsRefreshTokenValid(plain string) (primitive.ObjectID, bool, error) {
	if refreshCollection == nil {
		return primitive.NilObjectID, false, errors.New("refresh collection not initialized")
	}

	parts := strings.SplitN(plain, ":", 2)
	if len(parts) != 2 {
		return primitive.NilObjectID, false, errors.New("invalid token format")
	}
	tokenID := parts[0]
	tokenValue := parts[1]

	var doc struct {
		ID        string             `bson:"_id"`
		UserID    primitive.ObjectID `bson:"user_id"`
		TokenHash string             `bson:"token_hash"`
		ExpiresAt time.Time          `bson:"expires_at"`
		Revoked   bool               `bson:"revoked"`
	}
	ctx := context.Background()
	err := refreshCollection.FindOne(ctx, bson.M{"_id": tokenID}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return primitive.NilObjectID, false, nil
		}
		return primitive.NilObjectID, false, err
	}

	if doc.Revoked || time.Now().After(doc.ExpiresAt) {
		return doc.UserID, false, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(doc.TokenHash), []byte(tokenValue)); err != nil {
		return doc.UserID, false, nil
	}

	return doc.UserID, true, nil
}

// helper: options to return document after update (to read failed_attempts)
func optionsAfter() *options.FindOneAndUpdateOptions {
	opt := options.FindOneAndUpdate()
	opt.SetReturnDocument(options.After)
	return opt
}
