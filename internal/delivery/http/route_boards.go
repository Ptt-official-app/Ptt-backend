package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"github.com/Ptt-official-app/go-bbs"
)

func (delivery *httpDelivery) getBoardList(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getBoardList: %v", r)

	token := delivery.getTokenFromRequest(r)
	userId, err := delivery.usecase.GetUserIdFromToken(token)
	if err != nil {
		// user permission error
		// Support Guest?
		if !supportGuest() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"token_invalid"}`))
			return
		} else {
			userId = "guest" // TODO: use const variable
		}
	}

	boards := delivery.usecase.GetBoards(context.Background(), userId)

	dataList := make([]interface{}, 0, len(boards))
	for _, board := range boards {
		dataList = append(dataList, marshalBoardHeader(board))
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
}

func (delivery *httpDelivery) getPopularBoardList(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getPopularBoardList: %v", r)

	boards, err := delivery.usecase.GetPopularBoards(context.Background())
	if err != nil {
		// TODO: record error
		delivery.logger.Errorf("find popular board failed: %v", err)
		m := map[string]string{
			"error":             "find_popular_board_error",
			"error_description": "get popular board failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(b)
		return
	}

	dataList := make([]interface{}, 0, len(boards))
	for _, board := range boards {
		dataList = append(dataList, marshalBoardHeader(board))
	}

	responseMap := map[string]interface{}{
		"data": struct {
			Items []interface{} `json:"items"`
		}{dataList},
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
}

func (delivery *httpDelivery) getBoardInformation(w http.ResponseWriter, r *http.Request, boardId string) {
	delivery.logger.Debugf("getBoardInformation: %v", r)
	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	brd, err := delivery.usecase.GetBoardByID(context.Background(), boardId)
	if err != nil {
		// TODO: record error
		delivery.logger.Warningf("find board %s failed: %v", boardId, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_board_error",
			"error_description": "get board for " + boardId + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.Write(b)
		return
	}

	limitation, err := delivery.usecase.GetBoardPostsLimitation(context.Background(), boardId)
	if err != nil {
		delivery.logger.Warningf("get board %s post_limitation failed: %v", boardId, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "get_board_post_limitation_error",
			"error_description": "get board post_limitation for " + boardId + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.Write(b)
		return
	}

	responseMap := map[string]interface{}{
		"data": marshalBoardHeaderInfo(brd, limitation),
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
}

func (delivery *httpDelivery) getBoardSettings(w http.ResponseWriter, r *http.Request, boardId string) {
	delivery.logger.Debugf("getBoardSettings: %v", r)
	token := delivery.getTokenFromRequest(r)
	err := delivery.usecase.CheckPermission(token,
		[]usecase.Permission{usecase.PermissionReadBoardInformation},
		map[string]string{
			"board_id": boardId,
		})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	brd, err := delivery.usecase.GetBoardByID(context.Background(), boardId)
	if err != nil {
		// TODO: record error
		delivery.logger.Warningf("find board %s failed: %v", boardId, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_board_error",
			"error_description": "get board for " + boardId + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		w.Write(b)
		return
	}

	responseMap := map[string]interface{}{
		"data": marshalBoardHeaderSettings(brd),
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	w.Write(b)
}

// marshal generate board or class metadata object,
// b is input header
func marshalBoardHeader(b bbs.BoardRecord) map[string]interface{} {
	return marshalBoardHeaderInfo(b, nil)
}

func marshalBoardHeaderInfo(b bbs.BoardRecord, l *usecase.BoardPostLimitation) map[string]interface{} {
	ret := map[string]interface{}{
		"title":          b.Title(),
		"number_of_user": "0",
		"moderators":     b.BM(),
	}
	if b.IsClass() {
		// class
		// Assign ID from foreach loop
		ret["type"] = "class"
	} else {
		// board
		ret["id"] = b.BoardId()
		ret["type"] = "board"
	}
	if l != nil {
		ret["post_limitation"] = map[string]interface{}{
			"posts":   strconv.Itoa(int(l.PostsLimit)),
			"logins":  strconv.Itoa(int(l.LoginsLimit)),
			"badpost": strconv.Itoa(int(l.BadPostLimit)),
		}
	}
	return ret
}

func marshalBoardHeaderSettings(b bbs.BoardRecord) map[string]interface{} {
	ret := map[string]interface{}{
		// "hide":                               b.IsHide(),
		// "restricted_post_or_read_permission": b.IsPostMask(),
		// "anonymous":                          b.IsAnonymous(),
		// "default_anonymous":                  b.IsDefaultAnonymous(),
		// "no_money":                           b.IsNoCredit(),
		// "vote_board":                         b.IsVoteBoard(),
		// "warnel":                             b.IsWarnEL(),
		// "top":                                b.IsTop(),
		// "no_comment":                         b.IsNoRecommend(),
		// "angel_anonymous":                    b.IsAngelAnonymous(),
		// "bm_count":                           b.IsBMCount(),
		// "no_boo":                             b.IsNoBoo(),
		// "allow_list_post_only":               b.IsRestrictedPost(),
		// "guest_post_only":                    b.IsGuestPost(),
		// "cooldown":                           b.IsCooldown(),
		// "cross_post_log":                     b.IsCPLog(),
		// "no_fast_comment":                    b.IsNoFastRecommend(),
		// "log_ip_when_comment":                b.IsIPLogRecommend(),
		// "over18":                             b.IsOver18(),
		// "no_reply":                           b.IsNoReply(),
		// "aligned_comment":                    b.IsAlignedComment(),
		// "no_self_delete_post":                b.IsNoSelfDeletePost(),
		// "bm_mask_content":                    b.IsBMMaskContent(),
		"hide":                               false,
		"restricted_post_or_read_permission": false,
		"anonymous":                          false,
		"default_anonymous":                  false,
		"no_money":                           false,
		"vote_board":                         false,
		"warnel":                             false,
		"top":                                false,
		"no_comment":                         false,
		"angel_anonymous":                    false,
		"bm_count":                           false,
		"no_boo":                             false,
		"allow_list_post_only":               false,
		"guest_post_only":                    false,
		"cooldown":                           false,
		"cross_post_log":                     false,
		"no_fast_comment":                    false,
		"log_ip_when_comment":                false,
		"over18":                             false,
		"no_reply":                           false,
		"aligned_comment":                    false,
		"no_self_delete_post":                false,
		"bm_mask_content":                    false,
	}
	return ret
}
