package repositories

import (
	"fmt"

	"flutter-bengkel/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(id uuid.UUID, user *models.User) error
	SoftDelete(id uuid.UUID, deletedBy uuid.UUID) error
	Restore(id uuid.UUID) error
	PermanentDelete(id uuid.UUID) error
	List(offset, limit int, includeDeleted bool) ([]models.User, int64, error)
	UpdateLastLogin(id uuid.UUID) error
	ChangePassword(id uuid.UUID, hashedPassword string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, full_name, phone, role_id, outlet_id, is_active)
		VALUES (:username, :email, :password_hash, :full_name, :phone, :role_id, :outlet_id, :is_active)
	`
	
	result, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get user ID: %w", err)
	}
	
	user.ID = id
	return nil
}

func (r *userRepository) GetByID(id int64) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.phone, 
			   u.role_id, u.outlet_id, u.is_active, u.last_login_at, u.created_at, u.updated_at,
			   r.id as "role.id", r.name as "role.name", r.description as "role.description",
			   o.id as "outlet.id", o.name as "outlet.name", o.address as "outlet.address"
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN outlets o ON u.outlet_id = o.id
		WHERE u.id = ?
	`
	
	var user models.User
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.phone, 
			   u.role_id, u.outlet_id, u.is_active, u.last_login_at, u.created_at, u.updated_at,
			   r.id as "role.id", r.name as "role.name", r.description as "role.description",
			   o.id as "outlet.id", o.name as "outlet.name", o.address as "outlet.address"
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN outlets o ON u.outlet_id = o.id
		WHERE u.username = ?
	`
	
	var user models.User
	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.full_name, u.phone, 
			   u.role_id, u.outlet_id, u.is_active, u.last_login_at, u.created_at, u.updated_at,
			   r.id as "role.id", r.name as "role.name", r.description as "role.description",
			   o.id as "outlet.id", o.name as "outlet.name", o.address as "outlet.address"
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN outlets o ON u.outlet_id = o.id
		WHERE u.email = ?
	`
	
	var user models.User
	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	
	return &user, nil
}

func (r *userRepository) Update(id int64, user *models.User) error {
	query := `
		UPDATE users 
		SET username = :username, email = :email, full_name = :full_name, 
			phone = :phone, role_id = :role_id, outlet_id = :outlet_id, 
			is_active = :is_active
		WHERE id = :id
	`
	
	user.ID = id
	_, err := r.db.NamedExec(query, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	
	return nil
}

func (r *userRepository) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	return nil
}

func (r *userRepository) List(offset, limit int) ([]models.User, int64, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM users`
	var total int64
	err := r.db.Get(&total, countQuery)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}
	
	// Get users
	query := `
		SELECT u.id, u.username, u.email, u.full_name, u.phone, 
			   u.role_id, u.outlet_id, u.is_active, u.last_login_at, u.created_at, u.updated_at,
			   r.id as "role.id", r.name as "role.name", r.description as "role.description",
			   o.id as "outlet.id", o.name as "outlet.name", o.address as "outlet.address"
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		LEFT JOIN outlets o ON u.outlet_id = o.id
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?
	`
	
	var users []models.User
	err = r.db.Select(&users, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	
	return users, total, nil
}

func (r *userRepository) UpdateLastLogin(id int64) error {
	query := `UPDATE users SET last_login_at = NOW() WHERE id = ?`
	
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	
	return nil
}

func (r *userRepository) ChangePassword(id int64, hashedPassword string) error {
	query := `UPDATE users SET password_hash = ? WHERE id = ?`
	
	_, err := r.db.Exec(query, hashedPassword, id)
	if err != nil {
		return fmt.Errorf("failed to change password: %w", err)
	}
	
	return nil
}

// Role Repository
type RoleRepository interface {
	GetByID(id uuid.UUID) (*models.Role, error)
	GetPermissionsByRoleID(roleID uuid.UUID) ([]models.Permission, error)
	List() ([]models.Role, error)
}

type roleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetByID(id int64) (*models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = ?`
	
	var role models.Role
	err := r.db.Get(&role, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	
	return &role, nil
}

func (r *roleRepository) GetPermissionsByRoleID(roleID int64) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at
		FROM permissions p
		INNER JOIN role_has_permissions rhp ON p.id = rhp.permission_id
		WHERE rhp.role_id = ?
	`
	
	var permissions []models.Permission
	err := r.db.Select(&permissions, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get permissions: %w", err)
	}
	
	return permissions, nil
}

func (r *roleRepository) List() ([]models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles ORDER BY name`
	
	var roles []models.Role
	err := r.db.Select(&roles, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list roles: %w", err)
	}
	
	return roles, nil
}

// Permission Repository
type PermissionRepository interface {
	List() ([]models.Permission, error)
}

type permissionRepository struct {
	db *sqlx.DB
}

func NewPermissionRepository(db *sqlx.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) List() ([]models.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at, updated_at FROM permissions ORDER BY resource, action`
	
	var permissions []models.Permission
	err := r.db.Select(&permissions, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list permissions: %w", err)
	}
	
	return permissions, nil
}

// Outlet Repository
type OutletRepository interface {
	GetByID(id uuid.UUID) (*models.Outlet, error)
	List() ([]models.Outlet, error)
}

type outletRepository struct {
	db *sqlx.DB
}

func NewOutletRepository(db *sqlx.DB) OutletRepository {
	return &outletRepository{db: db}
}

func (r *outletRepository) GetByID(id int64) (*models.Outlet, error) {
	query := `SELECT id, name, address, phone, email, is_active, created_at, updated_at FROM outlets WHERE id = ?`
	
	var outlet models.Outlet
	err := r.db.Get(&outlet, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get outlet: %w", err)
	}
	
	return &outlet, nil
}

func (r *outletRepository) List() ([]models.Outlet, error) {
	query := `SELECT id, name, address, phone, email, is_active, created_at, updated_at FROM outlets WHERE is_active = true ORDER BY name`
	
	var outlets []models.Outlet
	err := r.db.Select(&outlets, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list outlets: %w", err)
	}
	
	return outlets, nil
}