package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Ptt-official-app/Ptt-backend/internal/usecase"
	"github.com/Ptt-official-app/go-bbs"
)

func (delivery *Delivery) getBoardList(w http.ResponseWriter, r *http.Request) {
	delivery.logger.Debugf("getBoardList: %v", r)

	token := delivery.getTokenFromRequest(r)
	userID, err := delivery.usecase.GetUserIDFromToken(token)
	if err != nil {
		userID = "guest" // TODO: use const variable
		// user permission error
		// Support Guest?
		if !supportGuest() {
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte(`{"error":"token_invalid"}`))
			if err != nil {
				delivery.logger.Errorf("getBoardList write token invalid response err: %w", err)
			}
			return
		}
	}

	boards := delivery.usecase.GetBoards(context.Background(), userID)

	dataList := make([]interface{}, 0, len(boards))
	for _, board := range boards {
		dataList = append(dataList, marshalBoardHeader(board))
	}

	responseMap := map[string]interface{}{
		"data": dataList,
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getBoardList write response err: %w", err)
	}
}

func (delivery *Delivery) getPopularBoardList(w http.ResponseWriter, r *http.Request) {
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
		_, err := w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getPopularBoardList write error response err: %w", err)
		}
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
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getPopularBoardList write success response err: %w", err)
	}
}

func (delivery *Delivery) getBoardInformation(w http.ResponseWriter, r *http.Request, boardID string) {
	delivery.logger.Debugf("getBoardInformation: %v", r)
	token := delivery.getTokenFromRequest(r)
	ctx := context.Background()

	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadBoardInformation}, map[string]string{
		"board_id": boardID,
	})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	brd, err := delivery.usecase.GetBoardByID(ctx, boardID)
	if err != nil {
		// TODO: record error
		delivery.logger.Warningf("find board %s failed: %v", boardID, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_board_error",
			"error_description": "get board for " + boardID + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getBoardInformation write error response err: %w", err)
		}
		return
	}

	limitation, err := delivery.usecase.GetBoardPostsLimitation(ctx, boardID)
	if err != nil {
		delivery.logger.Warningf("get board %s post_limitation failed: %v", boardID, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "get_board_post_limitation_error",
			"error_description": "get board post_limitation for " + boardID + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getBoardInformation write success response err: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": marshalBoardHeaderWithInfo(brd, limitation),
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getBoardInformation write success response err: %w", err)
	}
}

func (delivery *Delivery) getBoardSettings(w http.ResponseWriter, r *http.Request, boardID string) {
	delivery.logger.Debugf("getBoardSettings: %v", r)
	token := delivery.getTokenFromRequest(r)
	ctx := context.Background()

	err := delivery.usecase.CheckPermission(token, []usecase.Permission{usecase.PermissionReadBoardInformation}, map[string]string{
		"board_id": boardID,
	})

	if err != nil {
		// TODO: record unauthorized access
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	brd, err := delivery.usecase.GetBoardByID(ctx, boardID)
	if err != nil {
		// TODO: record error
		delivery.logger.Warningf("find board %s failed: %v", boardID, err)
		w.WriteHeader(http.StatusInternalServerError)
		m := map[string]string{
			"error":             "find_board_error",
			"error_description": "get board for " + boardID + " failed",
		}
		b, _ := json.MarshalIndent(m, "", "  ")
		_, err = w.Write(b)
		if err != nil {
			delivery.logger.Errorf("getBoardSettings write error response err: %w", err)
		}
		return
	}

	responseMap := map[string]interface{}{
		"data": marshalBoardSettings(brd),
	}

	b, _ := json.MarshalIndent(responseMap, "", "  ")
	_, err = w.Write(b)
	if err != nil {
		delivery.logger.Errorf("getBoardSettings write success response err: %w", err)
	}
}

// marshal generate board or class metadata object,
// b is input header
func marshalBoardHeader(b bbs.BoardRecord) map[string]interface{} {
	return marshalBoardHeaderWithInfo(b, nil)
}

func marshalBoardHeaderWithInfo(b bbs.BoardRecord, l *usecase.BoardPostLimitation) map[string]interface{} {
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
		ret["id"] = b.BoardID()
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

func marshalBoardSettings(b bbs.BoardRecord) map[string]interface{} {
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
