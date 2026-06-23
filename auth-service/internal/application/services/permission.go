package services

import "auth-service/internal/domain"

func PermissionsForUser(userType domain.UserType, role domain.UserRole) []string {
	if userType == domain.UserTypeCustomer {
		return []string{
			"profile:read",
			"profile:write",
			"cart:read",
			"cart:write",
			"orders:read_own",
			"orders:create",
		}
	}

	switch role {
	case domain.UserRoleAdminViewer:
		return []string{
			"clothes:read",
			"images:read",
			"change_logs:read",
		}

	case domain.UserRoleAdminEditor:
		return []string{
			"clothes:read",
			"clothes:write",
			"images:read",
			"images:write",
			"change_logs:read",
		}

	case domain.UserRoleAdminOwner:
		return []string{
			"clothes:read",
			"clothes:write",
			"images:read",
			"images:write",
			"change_logs:read",
			"admin_users:read",
			"admin_users:write",
			"customers:read",
			"customers:write",
		}

	default:
		return []string{}
	}
}
