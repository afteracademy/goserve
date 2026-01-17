package user

import (
	"context"

	"github.com/afteracademy/goserve/api/user/dto"
	"github.com/afteracademy/goserve/api/user/model"
	"github.com/afteracademy/goserve/arch/network"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	FetchUserPrivateProfile(user *model.User) (*dto.UserPrivate, error)
	FetchUserPublicProfile(userId uuid.UUID) (*dto.UserPublic, error)
	FindRoleByCode(ctx context.Context, code model.RoleCode) (*model.Role, error)
	FindRoles(ctx context.Context, roleIds []uuid.UUID) ([]*model.Role, error)
	FindUserById(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindUserRoles(ctx context.Context, user model.User) ([]*model.Role, error)
	FindUserPrivateProfile(ctx context.Context, user *model.User) (*model.User, error)
	FindUserPublicProfile(ctx context.Context, userId uuid.UUID) (*model.User, error)
	DeleteUserByEmail(ctx context.Context, email string) (bool, error)
}

type service struct {
	network.BaseService
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) Service {
	return &service{
		BaseService: network.NewBaseService(),
		db:          db,
	}
}

func (s *service) FetchUserPrivateProfile(user *model.User) (*dto.UserPrivate, error) {
	// leverage the role from the auth context
	return dto.NewUserPrivate(user), nil
}

func (s *service) FetchUserPublicProfile(userId uuid.UUID) (*dto.UserPublic, error) {
	user, err := s.FindUserPublicProfile(context.Background(), userId)
	if err != nil {
		return nil, network.NewNotFoundError("user does not exists", err)
	}
	return dto.NewUserPublic(user), nil
}

func (s *service) FindRoleByCode(
	ctx context.Context,
	code model.RoleCode,
) (*model.Role, error) {

	query := `
		SELECT
			id,
			code,
			status,
			created_at,
			updated_at
		FROM roles
		WHERE code = $1
		  AND status = TRUE
	`

	var role model.Role

	err := s.db.QueryRow(ctx, query, code).
		Scan(
			&role.ID,
			&role.Code,
			&role.Status,
			&role.CreatedAt,
			&role.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *service) FindRoles(
	ctx context.Context,
	roleIDs []uuid.UUID,
) ([]*model.Role, error) {

	if len(roleIDs) == 0 {
		return []*model.Role{}, nil
	}

	query := `
		SELECT
			id,
			code,
			status,
			created_at,
			updated_at
		FROM roles
		WHERE id = ANY($1)
		  AND status = TRUE
	`

	rows, err := s.db.Query(ctx, query, roleIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*model.Role

	for rows.Next() {
		var role model.Role
		if err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.Status,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *service) FindUserById(
	ctx context.Context,
	id uuid.UUID,
) (*model.User, error) {

	userQuery := `
		SELECT
			id,
			email,
			name,
			profile_pic_url,
			verified,
			status,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		  AND status = TRUE
	`

	var user model.User

	err := s.db.QueryRow(ctx, userQuery, id).
		Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.ProfilePicURL,
			&user.Verified,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	roles, err := s.FindUserRoles(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

func (s *service) FindUserByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {

	userQuery := `
		SELECT
			id,
			email,
			password,
			name,
			profile_pic_url,
			verified,
			status,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
		  AND status = TRUE
	`

	var user model.User

	err := s.db.QueryRow(ctx, userQuery, email).
		Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.ProfilePicURL,
			&user.Verified,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	roles, err := s.FindUserRoles(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Roles = roles

	return &user, nil
}

func (s *service) FindUserRoles(ctx context.Context, user model.User) ([]*model.Role, error) {
	roleQuery := `
		SELECT
			r.id,
			r.code,
			r.status,
			r.created_at,
			r.updated_at
		FROM roles r
		INNER JOIN user_roles ur
			ON ur.role_id = r.id
		WHERE ur.user_id = $1
		  AND r.status = TRUE
	`
	rows, err := s.db.Query(ctx, roleQuery, user.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(
			&role.ID,
			&role.Code,
			&role.Status,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *service) CreateUser(
	ctx context.Context,
	user *model.User,
) (*model.User, error) {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Insert user
	query := `
		INSERT INTO users (
			email,
			password,
			name,
			profile_pic_url,
			verified,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err = tx.QueryRow(
		ctx,
		query,
		user.Email,
		user.Password,
		user.Name,
		user.ProfilePicURL,
		user.Verified,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return nil, err
	}

	// Insert user roles
	const roleInsert = `
		INSERT INTO user_roles (user_id, role_id)
		VALUES ($1, $2)
	`

	for _, role := range user.Roles {
		_, err := tx.Exec(ctx, roleInsert, user.ID, role.ID)
		if err != nil {
			return nil, err
		}
	}

	// Commit
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) FindUserPrivateProfile(
	ctx context.Context,
	user *model.User,
) (*model.User, error) {

	query := `
		SELECT
			id,
			email,
			name,
			profile_pic_url,
			verified,
			status,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
		  AND status = TRUE
	`

	var result model.User

	err := s.db.QueryRow(ctx, query, user.ID).
		Scan(
			&result.ID,
			&result.Email,
			&result.Name,
			&result.ProfilePicURL,
			&result.Verified,
			&result.Status,
			&result.CreatedAt,
			&result.UpdatedAt,
		)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *service) FindUserPublicProfile(
	ctx context.Context,
	userID uuid.UUID,
) (*model.User, error) {

	query := `
		SELECT
			id,
			name,
			profile_pic_url
		FROM users
		WHERE id = $1
		  AND status = TRUE
	`

	var user model.User

	err := s.db.QueryRow(ctx, query, userID).
		Scan(&user.ID, &user.Name, &user.ProfilePicURL)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) DeleteUserByEmail(ctx context.Context, email string) (bool, error) {
	query := `
		DELETE FROM users
		WHERE email = $1
	`

	tag, err := s.db.Exec(ctx, query, email)
	if err != nil {
		return false, err
	}

	return tag.RowsAffected() > 0, nil
}
