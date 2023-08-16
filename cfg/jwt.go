package cfg

import (
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

var JwtKey = []byte(os.Getenv("JWT_PRIVATE_KEY"))

// Simpan token yang sudah digunakan dalam map
var UsedTokens = make(map[string]bool)

// payload untuk token
type Claims struct {
	Username string `json:"username"`
	UserID   int64  `json:"user_id"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// fungsi untuk membuat token
func CreateToken(username string, userID int64, role string) (string, error) {
	//tokenTTL, _ := strconv.Atoi(os.Getenv("TOKEN_TTL"))
	// mengatur waktu kadaluwarsa token
	expirationTime := time.Now().Add(time.Minute * 10)

	// membuat claims
	claims := &Claims{
		Username: username,
		UserID:   userID,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// membuat token dengan signing method HS256 dan secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)

	if err != nil {
		return "", err
	}

	//////simpan token yang di generate
	//UsedTokens[tokenString] = false

	return tokenString, nil
}
