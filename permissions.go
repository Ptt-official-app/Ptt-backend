package main

type permission string

const (
	PermissionReadUserInformation permission = "READ_USER_INFORMATION"
	PermissionReadFavorite        permission = "READ_FAVORITE"
)

func checkTokenPermission(token string, permissionId []permission, userInfo map[string]string) error {
	return nil
}
