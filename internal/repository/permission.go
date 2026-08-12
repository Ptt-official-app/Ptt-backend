package repository

import (
	"github.com/Ptt-official-app/go-bbs"
	"github.com/Ptt-official-app/go-bbs/pttbbs"
)

const pttOpenBoardPermissionMask uint32 = 0o37

type permissionLevelProvider interface {
	PermissionLevel() uint32
}

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

func BoardIsPublic(record bbs.BoardRecord) bool {
	if record == nil || record.BoardID() == "" || record.IsClass() {
		return false
	}
	settings, ok := record.(bbs.BoardRecordSettings)
	if !ok {
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
