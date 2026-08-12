package repository

import (
	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

// permissionLevelProvider lets other BBS implementations expose permission
// levels without making the usecase layer depend on a concrete BBS driver.
type permissionLevelProvider interface {
	PermissionLevel() uint32
}

// BoardReadPermissionLevel returns the permission bits required to read a
// board. On PTT, a board with BRD_POSTMASK uses Level for posting restrictions,
// so it must not be treated as a read restriction.
func BoardReadPermissionLevel(record bbs.BoardRecord) (uint32, bool) {
	if settings, ok := record.(bbs.BoardRecordSettings); ok && settings.IsPostMask() {
		return 0, true
	}

	if levelRecord, ok := record.(permissionLevelProvider); ok {
		return levelRecord.PermissionLevel(), true
	}

	rawRecord, ok := record.(*pttbbs.BoardHeader)
	if !ok {
		return 0, false
	}

	return rawRecord.Level, true
}

// UserPermissionLevel returns the permission bits assigned to a user.
func UserPermissionLevel(record bbs.UserRecord) (uint32, bool) {
	if levelRecord, ok := record.(permissionLevelProvider); ok {
		return levelRecord.PermissionLevel(), true
	}

	if wrappedRecord, ok := record.(*bbsUserRecord); ok {
		return UserPermissionLevel(wrappedRecord.UserRecord)
	}

	rawRecord, ok := record.(*pttbbs.Userec)
	if !ok {
		return 0, false
	}

	return rawRecord.UserLevel, true
}
