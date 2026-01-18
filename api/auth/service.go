package auth

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/afteracademy/goserve/api/auth/dto"
	"github.com/afteracademy/goserve/api/auth/model"
	"github.com/afteracademy/goserve/api/user"
	userModel "github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/afteracademy/goserve/config"
	"github.com/afteracademy/goserve/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUpBasic(signUpDto *dto.SignUpBasic) (*dto.UserAuth, error)
	SignInBasic(signInDto *dto.SignInBasic) (*dto.UserAuth, error)
	RenewToken(tokenRefreshDto *dto.TokenRefresh, accessToken string) (*dto.Tokens, error)
	SignOut(keystore *model.Keystore) error
	IsEmailRegisted(email string) bool
	GenerateToken(user *userModel.User) (string, string, error)
	FetchKeystore(client *userModel.User, primaryKey string) (*model.Keystore, error)
	VerifyToken(tokenStr string) (*jwt.RegisteredClaims, error)
	DecodeToken(tokenStr string) (*jwt.RegisteredClaims, error)
	SignToken(claims jwt.RegisteredClaims) (string, error)
	ValidateClaims(claims *jwt.RegisteredClaims) bool
	FetchApiKey(key string) (*model.ApiKey, error)

	/*--------only for tests----------*/
	CreateApiKey(key string, version int, permissions []model.Permission, comments []string) (*model.ApiKey, error)
	DeleteApiKey(apikey *model.ApiKey) (bool, error)
	/*--------------------------------*/
}

type service struct {
	network.BaseService
	db          *pgxpool.Pool
	userService user.Service
	// token
	rsaPrivateKey        *rsa.PrivateKey
	rsaPublicKey         *rsa.PublicKey
	accessTokenValidity  time.Duration
	refreshTokenValidity time.Duration
	tokenIssuer          string
	tokenAudience        string
}

func NewService(
	db *pgxpool.Pool,
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
		BaseService: network.NewBaseService(),
		userService: userService,
		db:          db,
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

	role, err := s.userService.FetchRoleByCode(userModel.RoleCodeLearner)
	if err != nil {
		return nil, err
	}
	roles := make([]*userModel.Role, 1)
	roles[0] = role

	hashed, err := bcrypt.GenerateFromPassword([]byte(signUpDto.Password), 5)
	if err != nil {
		return nil, err
	}

	user, err := s.userService.CreateUser(signUpDto.Email, string(hashed), signUpDto.Name, signUpDto.ProfilePicUrl, roles)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	tokens := dto.NewTokens(accessToken, refreshToken)
	return dto.NewUserAuth(user, tokens), nil
}

func (s *service) SignInBasic(signInDto *dto.SignInBasic) (*dto.UserAuth, error) {
	user, err := s.userService.FetchUserByEmail(signInDto.Email)
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

	tokens := dto.NewTokens(accessToken, refreshToken)
	return dto.NewUserAuth(user, tokens), nil
}

func (s *service) SignOut(keystore *model.Keystore) error {
	ctx := context.Background()

	query := `
		DELETE FROM keystore
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, keystore.ID)
	return err
}

func (s *service) IsEmailRegisted(email string) bool {
	exists, _ := s.userService.IsEmailExists(email)
	return exists
}

func (s *service) RenewToken(tokenRefreshDto *dto.TokenRefresh, accessToken string) (*dto.Tokens, error) {
	ctx := context.Background()

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

	userId, _ := uuid.Parse(refreshClaims.Subject)
	user, err := s.userService.FetchUserById(userId)
	if err != nil {
		return nil, network.NewUnauthorizedError("permission denied: invalid refresh claims subject", nil)
	}

	keystore, err := s.FindRefreshKeystore(ctx, user, accessClaims.ID, refreshClaims.ID)
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

	return dto.NewTokens(accessToken, refreshToken), nil
}

func (s *service) GenerateToken(user *userModel.User) (string, string, error) {
	ctx := context.Background()
	primaryKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}
	secondaryKey, err := utils.GenerateRandomString(32)
	if err != nil {
		return "", "", err
	}

	_, err = s.CreateKeystore(ctx, user, primaryKey, secondaryKey)
	if err != nil {
		return "", "", err
	}

	now := jwt.NewNumericDate(time.Now())

	accessTokenClaims := jwt.RegisteredClaims{
		Issuer:    s.tokenIssuer,
		Subject:   user.ID.String(),
		Audience:  []string{s.tokenAudience},
		IssuedAt:  now,
		NotBefore: now,
		ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTokenValidity * time.Second)),
		ID:        primaryKey,
	}

	refreshTokenClaims := jwt.RegisteredClaims{
		Issuer:    s.tokenIssuer,
		Subject:   user.ID.String(),
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

func (s *service) GenerateKeystore(
	client *userModel.User,
	primaryKey string,
	secondaryKey string,
) (*model.Keystore, error) {
	return s.CreateKeystore(context.Background(), client, primaryKey, secondaryKey)
}

func (s *service) CreateKeystore(
	ctx context.Context,
	client *userModel.User,
	primaryKey string,
	secondaryKey string,
) (*model.Keystore, error) {

	var ks = model.Keystore{}

	query := `
		INSERT INTO keystore (
			user_id,
			p_key,
			s_key
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			created_at,
			updated_at
	`

	err := s.db.QueryRow(
		ctx,
		query,
		client.ID,
		primaryKey,
		secondaryKey,
	).Scan(
		&ks.ID,
		&ks.CreatedAt,
		&ks.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ks, nil
}

func (s *service) FetchKeystore(
	client *userModel.User,
	primaryKey string,
) (*model.Keystore, error) {
	ctx := context.Background()
	query := `
		SELECT
			id,
			user_id,
			p_key,
			s_key,
			status,
			created_at,
			updated_at
		FROM keystore
		WHERE user_id = $1
		  AND p_key = $2
		  AND status = TRUE
	`

	var ks model.Keystore

	err := s.db.QueryRow(ctx, query, client.ID, primaryKey).
		Scan(
			&ks.ID,
			&ks.UserID,
			&ks.PrimaryKey,
			&ks.SecondaryKey,
			&ks.Status,
			&ks.CreatedAt,
			&ks.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &ks, nil
}

func (s *service) FindRefreshKeystore(
	ctx context.Context,
	client *userModel.User,
	primaryKey string,
	secondaryKey string,
) (*model.Keystore, error) {

	query := `
		SELECT
			id,
			user_id,
			p_key,
			s_key,
			status,
			created_at,
			updated_at
		FROM keystore
		WHERE user_id = $1
		  AND p_key = $2
		  AND s_key = $3
		  AND status = TRUE
	`

	var ks model.Keystore

	err := s.db.QueryRow(
		ctx,
		query,
		client.ID,
		primaryKey,
		secondaryKey,
	).Scan(
		&ks.ID,
		&ks.UserID,
		&ks.PrimaryKey,
		&ks.SecondaryKey,
		&ks.Status,
		&ks.CreatedAt,
		&ks.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &ks, nil
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

func (s *service) FetchApiKey(
	key string,
) (*model.ApiKey, error) {
	ctx := context.Background()
	query := `
		SELECT
			id,
			key,
			permissions,
			comments,
			version,
			status,
			created_at,
			updated_at
		FROM api_keys
		WHERE key = $1
		  AND status = TRUE
	`

	var apiKey model.ApiKey

	err := s.db.QueryRow(ctx, query, key).
		Scan(
			&apiKey.ID,
			&apiKey.Key,
			&apiKey.Permissions,
			&apiKey.Comments,
			&apiKey.Version,
			&apiKey.Status,
			&apiKey.CreatedAt,
			&apiKey.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (s *service) CreateApiKey(
	key string,
	version int,
	permissions []model.Permission,
	comments []string,
) (*model.ApiKey, error) {
	ctx := context.Background()
	apiKey := model.ApiKey{}

	query := `
		INSERT INTO api_keys (
			key,
			permissions,
			comments,
			version
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			status,
			created_at,
			updated_at
	`

	err := s.db.QueryRow(
		ctx,
		query,
		key,
		permissions,
		comments,
		version,
	).Scan(
		&apiKey.ID,
		&apiKey.Status,
		&apiKey.CreatedAt,
		&apiKey.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &apiKey, nil
}

func (s *service) DeleteApiKey(
	apiKey *model.ApiKey,
) (bool, error) {
	ctx := context.Background()
	query := `
		DELETE FROM api_keys
		WHERE id = $1
	`

	tag, err := s.db.Exec(ctx, query, apiKey.ID)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}
