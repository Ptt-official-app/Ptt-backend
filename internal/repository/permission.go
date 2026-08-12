package repository

import (
	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

// PTT considers these basic permission bits compatible with an open board.
// This mirrors IS_OPENBRD in pttbbs/include/perm.h.
const pttOpenBoardPermissionMask uint32 = 0o37

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

// BoardIsPublic reports whether a board may be included in public aggregate
// endpoints. The PTT-specific rules mirror IS_OPENBRD and additionally skip
// group boards because they do not contain normal article lists.
func BoardIsPublic(record bbs.BoardRecord) bool {
	if record == nil || record.BoardID() == "" || record.IsClass() {
		return false
	}

	settings, ok := record.(bbs.BoardRecordSettings)
	if !ok {
		// Aggregate endpoints must fail closed when visibility cannot be proven.
		return false
	}
	if settings.IsHide() || settings.IsTop() {
		return false
	}
	if settings.IsPostMask() {
		return true
	}

	level, ok := BoardReadPermissionLevel(record)
	if !ok {
		return false
	}

	return level&^pttOpenBoardPermissionMask == 0
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

// UserIsSYSOP reports whether the user has PTT's SYSOP permission bit.
func UserIsSYSOP(record bbs.UserRecord) bool {
	level, ok := UserPermissionLevel(record)
	return ok && level&pttbbs.PermSYSOP != 0
}
