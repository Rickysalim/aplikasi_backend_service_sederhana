package domain

import (
	"errors"
	"strings"
)

type RolePermissions struct {
	rolePermissions map[string][]string
}

func (p RolePermissions) IsAuthorizedFor(role string, routeName string) error {
	perms := p.rolePermissions[role]
	for _, r := range perms {
		if r == strings.TrimSpace(routeName) {
			return nil
		}
	}
	return errors.New("role is unathorized")
}

func GetRolePermissions() RolePermissions {
	return RolePermissions{map[string][]string{
		"admin": {"/v1/products","/v1/products/:product_id"},
		// "user": {"/v1/products","/v1/products/:product_id"},
		"user": {"/v1/products"},
	}}
}