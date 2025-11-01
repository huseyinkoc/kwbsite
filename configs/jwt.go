package configs

import (
    "log"
    "os"
)

// GetJWTSecret returns JWT signing key from env
func GetJWTSecret() string {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        log.Println("WARNING: JWT_SECRET not set, using insecure fallback (development only)")
        secret = "dev_fallback_secret"
    }
    return secret
}