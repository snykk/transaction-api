package records

import (
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type Users struct {
	Id        string     `db:"user_id"`
	Username  string     `db:"username"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Active    bool       `db:"active"`
	RoleId    int        `db:"role_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

// Mapper
func (u *Users) ToV1Domain() V1Domains.UserDomain {
	return V1Domains.UserDomain{
		ID:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Active:    u.Active,
		RoleID:    u.RoleId,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func FromUsersV1Domain(u *V1Domains.UserDomain) Users {
	return Users{
		Id:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Active:    u.Active,
		RoleId:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func ToArrayOfUsersV1Domain(u *[]Users) []V1Domains.UserDomain {
	var result []V1Domains.UserDomain

	for _, val := range *u {
		result = append(result, val.ToV1Domain())
	}

	return result
}
