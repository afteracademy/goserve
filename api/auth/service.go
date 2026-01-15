package auth

import (
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/afteracademy/goserve/api/auth/dto"
	"github.com/afteracademy/goserve/api/auth/model"
	"github.com/afteracademy/goserve/api/user"
	userModel "github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/mongo"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/config"
	"github.com/afteracademy/goserve/utils"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUpBasic(signUpDto *dto.SignUpBasic) (*dto.UserAuth, error)
	SignInBasic(signInDto *dto.SignInBasic) (*dto.UserAuth, error)
	RenewToken(tokenRefreshDto *dto.TokenRefresh, accessToken string) (*dto.UserTokens, error)
	SignOut(keystore *model.Keystore) error
	IsEmailRegisted(email string) bool
	GenerateToken(user *userModel.User) (string, string, error)
	CreateKeystore(client *userModel.User, primaryKey string, secondaryKey string) (*model.Keystore, error)
	FindKeystore(client *userModel.User, primaryKey string) (*model.Keystore, error)
	FindRefreshKeystore(client *userModel.User, pKey string, sKey string) (*model.Keystore, error)
	VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error)
	SignToken(claims jwt.RegisteredClaims) (string, error)
	ValidateClaims(claims *jwt.RegisteredClaims) bool
	FindApiKey(key string) (*model.ApiKey, error)
	CreateApiKey(key string, version int, permissions []model.Permission, comments []string) (*model.ApiKey, error)
	DeleteApiKey(apikey *model.ApiKey) (bool, error)
}

type service struct {
	network.BaseService
	keystoreQueryBuilder mongo.QueryBuilder[model.Keystore]
	apikeyQueryBuilder   mongo.QueryBuilder[model.ApiKey]
	userService          user.Service
	// token
	rsaPrivateKey        *rsa.PrivateKey
	rsaPublicKey         *rsa.PublicKey
	accessTokenValidity  time.Duration
	refreshTokenValidity time.Duration
	tokenIssuer          string
	tokenAudience        string
}

func NewService(
	db mongo.Database,
	env *config.Env,
	userService user.Service,
) Service {
	privatePem, err := utils.LoadPEMFileInto(env.RSAPrivateKeyPath)
	if err != nil {
		panic(err)
	}
	rsaPrivateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
	if err != nil {
		panic(err)
	}

	publicPem, err := utils.LoadPEMFileInto(env.RSAPublicKeyPath)
	if err != nil {
		panic(err)
	}

	rsaPublicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPem)
	if err != nil {
		panic(err)
	}

	return &service{
		BaseService:          network.NewBaseService(),
		userService:          userService,
		keystoreQueryBuilder: mongo.NewQueryBuilder[model.Keystore](db, model.KeystoreCollectionName),
		apikeyQueryBuilder:   mongo.NewQueryBuilder[model.ApiKey](db, model.ApiKeyCollectionName),
		// token key
		rsaPrivateKey: rsaPrivateKey,
		rsaPublicKey:  rsaPublicKey,
		// token claim
		accessTokenValidity:  time.Duration(env.AccessTokenValiditySec),
		refreshTokenValidity: time.Duration(env.RefreshTokenValiditySec),
		tokenIssuer:          env.TokenIssuer,
		tokenAudience:        env.TokenAudience,
	}
}

func (s *service) SignUpBasic(signUpDto *dto.SignUpBasic) (*dto.UserAuth, error) {
	exists := s.IsEmailRegisted(signUpDto.Email)
	if exists {
		return nil, network.NewBadRequestError("user already registered", nil)
	}

	role, err := s.userService.FindRoleByCode(userModel.RoleCodeLearner)
	if err != nil {
		return nil, err
	}
	roles := make([]*userModel.Role, 1)
	roles[0] = role

	hashed, err := bcrypt.GenerateFromPassword([]byte(signUpDto.Password), 5)
	if err != nil {
		return nil, err
	}

	user, err := userModel.NewUser(signUpDto.Email, string(hashed), signUpDto.Name, signUpDto.ProfilePicUrl, roles)
	if err != nil {
		return nil, err
	}

	user, err = s.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	tokens := dto.NewUserTokens(accessToken, refreshToken)
	return dto.NewUserAuth(user, tokens), nil
}

func (s *service) SignInBasic(signInDto *dto.SignInBasic) (*dto.UserAuth, error) {
	user, err := s.userService.FindUserByEmail(signInDto.Email)
	if err != nil {
		return nil, network.NewNotFoundError("user not registerd", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(signInDto.Password))
	if err != nil {
		return nil, network.NewUnauthorizedError("wrong password", err)
	}

	accessToken, refreshToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	tokens := dto.NewUserTokens(accessToken, refreshToken)
	return dto.NewUserAuth(user, tokens), nil
}

func (s *service) SignOut(keystore *model.Keystore) error {
	filter := bson.M{"_id": keystore.ID}
	_, err := s.keystoreQueryBuilder.SingleQuery().DeleteOne(filter)
	return err
}

func (s *service) IsEmailRegisted(email string) bool {
	user, _ := s.userService.FindUserByEmail(email)
	return user != nil
}

func (s *service) RenewToken(tokenRefreshDto *dto.TokenRefresh, accessToken string) (*dto.UserTokens, error) {
	accessClaims, err := s.DecodeToken(accessToken)
	if err != nil {
		return nil, err
	}

	valid := s.ValidateClaims(accessClaims)
	if !valid {
		return nil, network.NewUnauthorizedError("permission denied: invalid access claims", nil)
	}

	refreshClaims, err := s.VerifyToken(tokenRefreshDto.RefreshToken)
	if err != nil {
		return nil, err
	}

	valid = s.ValidateClaims(refreshClaims)
	if !valid {
		return nil, network.NewUnauthorizedError("permission denied: invalid refresh claims", nil)
	}

	if accessClaims.Subject != refreshClaims.Subject {
		return nil, network.NewUnauthorizedError("permission denied: access and refresh claims mismatch", nil)
	}

	userId, _ := mongo.NewObjectID(refreshClaims.Subject)
	user, err := s.userService.FindUserById(userId)
	if err != nil {
		return nil, network.NewUnauthorizedError("permission denied: invalid refresh claims subject", nil)
	}

	keystore, err := s.FindRefreshKeystore(user, accessClaims.ID, refreshClaims.ID)
	if err != nil {
		return nil, network.NewUnauthorizedError("permission denied: claims ids", nil)
	}

	err = s.SignOut(keystore)
	if err != nil {
		return nil, nil
	}

	accessToken, refreshToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return dto.NewUserTokens(accessToken, refreshToken), nil
}

func (s *service) GenerateToken(user *userModel.User) (string, string, error) {
	primaryKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	secondaryKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}

	_, err = s.CreateKeystore(user, primaryKey, secondaryKey)
	if err != nil {
		return "", "", err
	}

	now := jwt.NewNumericDate(time.Now())

	accessTokenClaims := jwt.RegisteredClaims{
		Issuer:    s.tokenIssuer,
		Subject:   user.ID.Hex(),
		Audience:  []string{s.tokenAudience},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenValidity * time.Second)),
		ID:        primaryKey,
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		Issuer:    s.tokenIssuer,
		Subject:   user.ID.Hex(),
		Audience:  []string{s.tokenAudience},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshTokenValidity * time.Second)),
		ID:        secondaryKey,
	}

	accessToken, err := s.SignToken(accessTokenClaims)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.SignToken(refreshTokenClaims)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) CreateKeystore(client *userModel.User, primaryKey string, secondaryKey string) (*model.Keystore, error) {
	doc, err := model.NewKeystore(client.ID, primaryKey, secondaryKey)
	if err != nil {
		return nil, err
	}

	id, err := s.keystoreQueryBuilder.SingleQuery().InsertOne(doc)
	if err != nil {
		return nil, err
	}

	doc.ID = *id
	return doc, nil
}

func (s *service) FindKeystore(client *userModel.User, primaryKey string) (*model.Keystore, error) {
	filter := bson.M{"client": client.ID, "pKey": primaryKey, "status": true}
	return s.keystoreQueryBuilder.SingleQuery().FindOne(filter, nil)
}

func (s *service) FindRefreshKeystore(client *userModel.User, primaryKey string, secondaryKey string) (*model.Keystore, error) {
	filter := bson.M{"client": client.ID, "pKey": primaryKey, "sKey": secondaryKey, "status": true}
	return s.keystoreQueryBuilder.SingleQuery().FindOne(filter, nil)
}

func (s *service) SignToken(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(s.rsaPrivateKey)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (s *service) VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(tkn *jwt.Token) (any, error) {
		return s.rsaPublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
			return claims, nil
		}
	}

	return nil, jwt.ErrTokenMalformed
}

func (s *service) DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(tkn *jwt.Token) (any, error) {
		return s.rsaPublicKey, nil
	})
	if token == nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

func (s *service) ValidateClaims(claims *jwt.RegisteredClaims) bool {
	invalid := claims.Issuer != s.tokenIssuer ||
		claims.Subject == "" ||
		len(claims.Audience) == 0 ||
		claims.Audience[0] != s.tokenAudience ||
		claims.NotBefore == nil ||
		claims.ExpiresAt == nil ||
		claims.ID == ""

	if invalid {
		return false
	}

	return utils.IsValidObjectID(claims.Subject)
}

func (s *service) FindApiKey(key string) (*model.ApiKey, error) {
	filter := bson.M{"key": key, "status": true}

	apikey, err := s.apikeyQueryBuilder.SingleQuery().FindOne(filter, nil)
	if err != nil {
		return nil, err
	}

	return apikey, nil
}

func (s *service) CreateApiKey(key string, version int, permissions []model.Permission, comments []string) (*model.ApiKey, error) {
	doc := model.NewApiKey(key, version, permissions, comments)

	id, err := s.apikeyQueryBuilder.SingleQuery().InsertOne(doc)
	if err != nil {
		return nil, err
	}

	doc.ID = *id
	return doc, nil
}

func (s *service) DeleteApiKey(apikey *model.ApiKey) (bool, error) {
	filter := bson.M{"_id": apikey.ID}
	result, err := s.apikeyQueryBuilder.SingleQuery().DeleteOne(filter)
	if err != nil {
		return false, err
	}
	return result.DeletedCount > 0, nil
}
