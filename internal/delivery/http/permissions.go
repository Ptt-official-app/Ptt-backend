package http

type permission string

const (
	PermissionReadUserInformation     permission = "READ_USER_INFORMATION"
	PermissionReadBoardInformation    permission = "READ_BOARD_INFORMATION"
	PermissionReadTreasureInformation permission = "READ_TREASURE_INFORMATION"
	PermissionReadFavorite            permission = "READ_FAVORITE"
)

func checkTokenPermission(token string, permissionId []permission, userInfo map[string]string) error {
	return nil
}
