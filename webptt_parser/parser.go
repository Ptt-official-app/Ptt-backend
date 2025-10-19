package webpttparser

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// GetPttPage fetches the HTML content of a PTT page given its URL path.
func GetPttPage(url string) ([]byte, error) {
	requestURL := fmt.Sprintf("https://www.ptt.cc/bbs/%v", url)
	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get page: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// HandlePage processes the HTML content of a PTT article page and extracts the main content.
func HandlePage(input []byte) []byte {
	buf := []byte{}

	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		slog.Error("parse error:", "err", err)
		return nil
	}
	// acceptMetaLines := true
	var mainContentNode *html.Node
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == "main-content" {
					mainContentNode = n
				}
			}
		}
	}

	if mainContentNode == nil {
		slog.Error("find main-content fail")
		return nil
	}
	state := "find-作者-line"
	for n := range mainContentNode.ChildNodes() {
		// slog.Info("scan Descendants.Child", "node", n.Data, "attr", n.Attr)
		if n.Type == html.TextNode {
			buf = append(buf, []byte(n.Data)...)
			continue
		}

		if n.Type != html.ElementNode {
			slog.Error("unknown node", "node", n.Data, "type", n.Type)
			continue
		}

		if n.Data == "a" {
			buf = handleATag(buf, n)
			continue
		} else if n.Data == "span" {
			buf = handleSpanTag(buf, n)
			continue
		} else if n.Data == "div" {
			isRichcontent := false
			for _, attr := range n.Attr {
				if attr.Key == "class" {
					if attr.Val == "richcontent" {
						// skip richcontent
						isRichcontent = true
						break
					}
				}
			}
			if isRichcontent {
				// skip richcontent
				continue
			}
		}

		for _, attr := range n.Attr {

			if attr.Key == "class" && attr.Val == "push" {
				buf = handleDivPushTag(buf, n)
				break
			}

			if attr.Key == "class" && attr.Val == "article-metaline" {
				// slog.Info("found", "cn", cn)
				// buf = append(buf, []byte("state:" + state + "\n")...)
				if state == "find-作者-line" {

					// buf = append(buf, []byte("\n\n\n\n")...)
					state2 := "find-作者-tag"
					for cn2 := range n.ChildNodes() {

						for _, attr := range cn2.Attr {
							if attr.Key == "class" && attr.Val == "article-meta-tag" {
								if state2 == "find-作者-tag" {
									// 作者
									for cn3 := range cn2.ChildNodes() {
										// buf = append(buf, []byte("    " + cn3.Data + "[")...)
										// for _,attr := range cn3.Attr {
										// 	buf = append(buf, []byte(attr.Key + "=" + attr.Val)...)
										// }
										// buf = append(buf, []byte("]\n")...)
										if cn3.Data == "作者" {
											buf = append(buf, []byte("作者: ")...)
										}
									}
									state2 = "find-作者-value"
								}
							} else if attr.Key == "class" && attr.Val == "article-meta-value" {
								if state2 == "find-作者-value" {
									// 作者值
									for cn3 := range cn2.ChildNodes() {
										// buf = append(buf, []byte("    " + cn3.Data + "[")...)
										// for _,attr := range cn3.Attr {
										// 	buf = append(buf, []byte(attr.Key + "=" + attr.Val)...)
										// }

										buf = append(buf, []byte(cn3.Data+" ")...)
									}
									state2 = "find-看板-tag"
								}
							}
						}
					}

					state = "find-看板-line"
				} else if state == "find-Meta-line" {
					// 標題
					state2 := "find-tag"
					for cn2 := range n.ChildNodes() {
						for _, attr := range cn2.Attr {
							if attr.Key == "class" && attr.Val == "article-meta-tag" {
								if state2 == "find-tag" {
									// 標題
									for cn3 := range cn2.ChildNodes() {
										// if cn3.Data == "標題"{
										buf = append(buf, []byte(cn3.Data+": ")...)
										// }
									}
									state2 = "find-value"
								}
							} else if attr.Key == "class" && attr.Val == "article-meta-value" {
								if state2 == "find-value" {
									// 標題值
									for cn3 := range cn2.ChildNodes() {
										buf = append(buf, []byte(cn3.Data+"\n")...)
									}
								}
							}
						}
					}
					state = "find-Meta-line"
				}
			} else if attr.Key == "class" && attr.Val == "article-metaline-right" {
				if state == "find-看板-line" {
					// 看板
					state2 := "find-tag"
					for cn2 := range n.ChildNodes() {
						for _, attr := range cn2.Attr {
							if attr.Key == "class" && attr.Val == "article-meta-tag" {
								if state2 == "find-tag" {
									// 看板
									for cn3 := range cn2.ChildNodes() {
										if cn3.Data == "看板" {
											buf = append(buf, []byte("看板: ")...)
										}
									}
									state2 = "find-value"
								}
							} else if attr.Key == "class" && attr.Val == "article-meta-value" {
								if state2 == "find-value" {
									// 看板值
									for cn3 := range cn2.ChildNodes() {
										buf = append(buf, []byte(cn3.Data+"\n")...)
									}
								}
							}
						}
					}
					state = "find-Meta-line"
				}
			} else {
				slog.Info("unknown node", "node", n.Data)
				if state == "find-Meta-line" {
					// dont parse header any more
					// buf = append(buf, []byte(n.Data)...)
					// acceptMetaLines = false
				}
			}
		} // scan node attr
	}
	slog.Info("handleFin")

	return buf
}

func handleATag(buf []byte, n *html.Node) []byte {
	// should have more a in children
	for cn := range n.ChildNodes() {
		if cn.Type == html.TextNode {
			buf = append(buf, []byte(cn.Data)...)
		}
	}
	return buf

}

const (
	ClassBgPrefix  string = "b"
	ClassfgPrefix  string = "f"
	ClassHighlight string = "hl"

	ClassPushTag = `push-tag`

	ClassPushUserId     = `push-userid`
	ClassPushContent    = `push-content`
	ClassPushIpDatetime = `push-ipdatetime`

	ClassArticleMetaLine      = `article-metaline`
	ClassArticleMetaLineRight = `article-metaline-right`
	ClassArticleMetaTag       = `article-meta-tag`
	ClassArticleMetaValue     = `article-meta-value`
)

// Flags
const (
	NoFlags = 1 << iota
	Highlighted
)

const (
	DefaultFg = 7
	DefaultBg = 0
)

type EscapeSequence struct {
	IsCSI        bool
	PrivateModes []rune // not supported, will not be parsed
	Nums         []int
	Trailings    []rune
	Mode         rune
}

func (e *EscapeSequence) Reset() {
	e.IsCSI = false

	if e.PrivateModes == nil {
		e.PrivateModes = make([]rune, 0, 4)
	} else {
		e.PrivateModes = e.PrivateModes[0:0]
	}

	if e.Nums == nil {
		e.Nums = make([]int, 0, 4)
	} else {
		e.Nums = e.Nums[0:0]
	}

	if e.Trailings == nil {
		e.Trailings = make([]rune, 0, 4)
	} else {
		e.Trailings = e.Trailings[0:0]
	}

	e.Mode = 0
}

func (e *EscapeSequence) ParseNumbers(buf []rune) {
	part := make([]rune, 0, 4)
	for i, r := range buf {
		if r != ';' {
			part = append(part, r)
		}
		if r == ';' || i == len(buf)-1 {
			switch len(part) {
			case 0:
				// Treat empty parameter as 0 (eg. "\e[;34m" as "\e[0;34m").
				// It is not stated in the spec whether it's valid. But most
				// terminal does this and are being relied by ascii art
				// creators. Let's allow this.
				e.Nums = append(e.Nums, 0)
			default:
				num, err := strconv.Atoi(string(part))
				if err != nil {
					continue // be nice
				}
				e.Nums = append(e.Nums, num)
				part = part[0:0]
			}
		}
	}
}

type TerminalState struct {
	fg, bg, flags int
}

func (t *TerminalState) Reset() {
	t.fg = DefaultFg
	t.bg = DefaultBg
	t.flags = NoFlags
}

func (t *TerminalState) IsDefaultState() bool {
	return t.fg == DefaultFg && t.bg == DefaultBg && t.flags == NoFlags
}

func (t *TerminalState) SetColor(fg, bg, flags int) {
	t.fg = fg
	t.bg = bg
	t.flags = flags
}

func (t *TerminalState) ApplyEscapeSequence(esc EscapeSequence) {
	switch esc.Mode {
	case 'm':
		if len(esc.Nums) == 0 {
			t.Reset()
			return
		}
		fg, bg, flags := t.fg, t.bg, t.flags
		for _, ctl := range esc.Nums {
			switch {
			case ctl == 0:
				fg = DefaultFg
				bg = DefaultBg
				flags = NoFlags
			case ctl == 1:
				flags |= Highlighted
			case ctl == 22:
				flags &= ^Highlighted
			case ctl >= 30 && ctl <= 37:
				fg = ctl % 10
			case ctl >= 40 && ctl <= 47:
				bg = ctl % 10
			default:
				// be nice
			}
		}
		t.SetColor(fg, bg, flags)
	}
}

func (t *TerminalState) Equal(u *TerminalState) bool {
	return t.fg == u.fg && t.bg == u.bg && t.flags == u.flags
}

func (t *TerminalState) Fg() int {
	return t.fg
}

func (t *TerminalState) Bg() int {
	return t.bg
}

func (t *TerminalState) Flags() int {
	return t.flags
}

func (t *TerminalState) HasFlags(f int) bool {
	return t.flags&f == f
}

func handleSpanTag(buf []byte, n *html.Node) []byte {
	// should have more a in children
	class := ""
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			class = attr.Val
		}
	}
	classes := strings.Split(class, " ")
	t := TerminalState{}
	t.Reset()
	for _, c := range classes {
		fg, bg, flags := t.Fg(), t.Bg(), t.Flags()
		if strings.HasPrefix(c, "f") { // `f`
			// foreground color
			v, err := strconv.Atoi(c[len("b"):])
			if err != nil {
				slog.Warn("strconv.Atoi error", "err", err)
				continue
			}
			fg = 30 + v
		} else if strings.HasPrefix(c, ClassBgPrefix) { // `b`
			// background color
			v, err := strconv.Atoi(c[len(ClassBgPrefix):])
			if err != nil {
				slog.Warn("strconv.Atoi error", "err", err)
				continue
			}
			bg = 40 + v
		} else if c == ClassHighlight { // `hl`
			// highlighted
			flags |= Highlighted
		} else if c == ClassPushTag ||
			c == ClassPushUserId ||
			c == ClassPushContent ||
			c == ClassPushIpDatetime {
			continue
		} else { // `p`
			slog.Warn("unknown class", "class", c)
		}
		t.SetColor(fg, bg, flags)
	}
	for cn := range n.ChildNodes() {
		// slog.Info("handleDivPushTag", "node", cn.Data, "attr", cn.Attr, "t", t)
		// escape := "\033"
		// if !t.IsDefaultState() {
		// slog.Info("handleDivPushTag not dfs", "node", cn.Data, "attr", cn.Attr, "t", t)
		escape := "\033["
		seg := []string{}
		if t.HasFlags(Highlighted) {
			seg = append(seg, "1")
		} else {
			seg = append(seg, "0")
		}
		if t.Fg() != DefaultFg {
			seg = append(seg, strconv.Itoa(t.Fg()))
		}
		if t.Bg() != DefaultBg {
			seg = append(seg, strconv.Itoa(t.Bg()))
		}
		escape += strings.Join(seg, ";") + "m"
		buf = append(buf, []byte(escape)...)
		// }
		if cn.Type == html.TextNode {
			buf = append(buf, []byte(cn.Data)...)
		} else {
			if cn.Data == "a" {
				buf = append(buf, []byte("\033[m")...)
				buf = handleATag(buf, cn)
				// buf = append(buf, []byte("\n")...)
			} else {
				slog.Warn("handleDivPushTag unknown tag in span", "node", cn.Data)
			}
		}
		if !t.IsDefaultState() {
			buf = append(buf, []byte("\033[m")...)
		}
	}
	return buf

}

func handleDivPushTag(buf []byte, n *html.Node) []byte {
	for cn := range n.ChildNodes() {
		if cn.Data == "span" {
			buf = handleSpanTag(buf, cn)
			continue
		}
		if cn.Type == html.ElementNode {
			// slog.Warn("handleDivPushTag unknown node", "node", cn.Data, "type", cn.Type)
			for cn2 := range cn.ChildNodes() {

				// slog.Info("handleDivPushTag", "node", cn2.Data, "attr", cn2.Attr)
				buf = append(buf, []byte(cn2.Data)...)
			}
		}
	}
	return buf
}

// ==========================
// Board article list parsing
// ==========================

// BoardEntry represents one row in a board index list (r-ent)
type BoardEntry struct {
	// Raw/derived values
	RecommendStr string // original string in nrec (e.g., "爆", "X1", "10", "")
	Recommend    int    // parsed number (爆=100, Xn=-n, empty=0)

	Title     string // article title or "(本文已被刪除) [user]"
	Owner     string // author id; may be "-" for deleted; try best-effort from title if deleted
	Date      string // e.g. " 4/02"
	Mark      string // sticky/mark column text; empty if none
	FileName  string // article file name, e.g., "M.1742152964.A.660"; for deleted set to "M.0.A.0"
	URL       string // href to the article if present; empty if deleted
	IsDeleted bool   // true if article is deleted (no link)
}

// BoardIndexPage captures a parsed board index page including pagination links
type BoardIndexPage struct {
	FirstPage string
	PrevPage  string
	NextPage  string
	LastPage  string

	Articles []BoardEntry
	Bottoms  []BoardEntry
}

// ParseBoardIndexPage parses PTT board index HTML and returns entries and pagination info.
func ParseBoardIndexPage(input []byte) (*BoardIndexPage, error) {
	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	out := &BoardIndexPage{}

	// Find pagination first (btn-group btn-group-paging)
	if node := findFirst(doc, func(n *html.Node) bool { return hasClass(n, "btn-group") && hasClass(n, "btn-group-paging") }); node != nil {
		for child := range node.ChildNodes() {
			if child.Type == html.ElementNode && child.Data == "a" {
				txt := strings.TrimSpace(textContent(child))
				href, _ := getAttr(child, "href")
				isDisabled := hasClass(child, "disabled") || href == ""
				switch {
				case strings.Contains(txt, "最舊"):
					if !isDisabled {
						out.FirstPage = href
					}
				case strings.Contains(txt, "上頁"):
					if !isDisabled {
						out.PrevPage = href
					}
				case strings.Contains(txt, "下頁"):
					if !isDisabled {
						out.NextPage = href
					}
				case strings.Contains(txt, "最新"):
					if !isDisabled {
						out.LastPage = href
					}
				}
			}
		}
	}

	// Find container of entries
	container := findFirst(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "r-list-container")
	})
	if container == nil {
		// Some pages may still be parsable by scanning all r-ent
		slog.Warn("r-list-container not found; scanning entire document for r-ent")
	}

	var entries []BoardEntry
	var bottoms []BoardEntry
	walker := container
	if walker == nil {
		walker = doc
	}

	if container != nil {
		afterSep := false
		for c := range container.ChildNodes() {
			if c.Type != html.ElementNode || c.Data != "div" {
				continue
			}
			if hasClass(c, "r-list-sep") {
				afterSep = true
				continue
			}
			if hasClass(c, "r-ent") {
				e := parseREnt(c)
				if afterSep {
					bottoms = append(bottoms, e)
				} else {
					entries = append(entries, e)
				}
			}
		}
	} else {
		// Fallback: just collect all r-ent as Articles
		for n := range walker.Descendants() {
			if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "r-ent") {
				entry := parseREnt(n)
				entries = append(entries, entry)
			}
		}
	}

	out.Articles = entries
	out.Bottoms = bottoms
	return out, nil
}

// =====================
// Hotboards page parser
// =====================

// HotboardEntry represents one board item in hotboards/classlist
type HotboardEntry struct {
	BrdName string // board-name
	Title   string // board-title (without Σ/◎ prefixes)
	Class   string // board-class
	Nuser   int    // board-nuser (online users)
	IsBoard bool   // true if it's a board (◎), false if group (Σ)
	URL     string // href of the entry
}

// HotboardsPage holds a list of hotboard entries
type HotboardsPage struct {
	Boards []HotboardEntry
}

// ParseHotboardsPage parses a PTT hotboards/classlist HTML page
// Expected structure similar to classlist.html with .b-list-container and .b-ent nodes.
func ParseHotboardsPage(input []byte) (*HotboardsPage, error) {
	doc, err := html.Parse(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}
	out := &HotboardsPage{}

	container := findFirst(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "b-list-container")
	})
	if container == nil {
		// fallback: scan entire doc
		container = doc
	}

	for n := range container.Descendants() {
		if n.Type == html.ElementNode && n.Data == "div" && hasClass(n, "b-ent") {
			// inside it there should be a link .board
			link := findFirstChild(n, func(c *html.Node) bool { return c.Type == html.ElementNode && c.Data == "a" && hasClass(c, "board") })
			if link == nil {
				continue
			}
			entry := HotboardEntry{}
			if href, ok := getAttr(link, "href"); ok {
				entry.URL = href
			}
			// children divs: board-name, board-nuser, board-class, board-title
			if name := findFirstChild(link, func(c *html.Node) bool {
				return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "board-name")
			}); name != nil {
				entry.BrdName = strings.TrimSpace(textContent(name))
			}
			if nuser := findFirstChild(link, func(c *html.Node) bool {
				return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "board-nuser")
			}); nuser != nil {
				s := strings.TrimSpace(textContent(nuser))
				if v, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
					entry.Nuser = v
				} else {
					// Some skins use formatted nuser like "1,234"; strip non-digits
					digits := regexp.MustCompile(`[^0-9]`).ReplaceAllString(s, "")
					if v, err := strconv.Atoi(digits); err == nil {
						entry.Nuser = v
					}
				}
			}
			if class := findFirstChild(link, func(c *html.Node) bool {
				return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "board-class")
			}); class != nil {
				entry.Class = strings.TrimSpace(textContent(class))
			}
			if title := findFirstChild(link, func(c *html.Node) bool {
				return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "board-title")
			}); title != nil {
				t := strings.TrimSpace(textContent(title))
				// Determine IsBoard based on prefix in template (◎ board, Σ group)
				// However, the visible symbol may not be present in raw text depending on renderer.
				// Fallback: detect from href pattern: /bbs/<name>/index.html => board; /cls/<bid> => group
				if strings.HasPrefix(t, "◎") {
					entry.IsBoard = true
					t = strings.TrimPrefix(t, "◎")
				} else if strings.HasPrefix(t, "Σ") || strings.HasPrefix(t, "&#931;") {
					entry.IsBoard = false
					t = strings.TrimPrefix(t, "Σ")
				}
				entry.Title = strings.TrimSpace(t)
			}
			if entry.URL != "" {
				if strings.HasPrefix(entry.URL, "/bbs/") {
					entry.IsBoard = true
				} else if strings.HasPrefix(entry.URL, "/cls/") {
					entry.IsBoard = false
				}
			}
			out.Boards = append(out.Boards, entry)
		}
	}

	return out, nil
}

func parseREnt(n *html.Node) BoardEntry {
	var e BoardEntry

	// nrec
	if nn := findFirstChild(n, func(c *html.Node) bool { return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "nrec") }); nn != nil {
		e.RecommendStr = strings.TrimSpace(textContent(nn))
		e.Recommend = parseRecommend(e.RecommendStr)
	}

	// title (a?)
	if tn := findFirstChild(n, func(c *html.Node) bool { return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "title") }); tn != nil {
		link := findFirstChild(tn, func(c *html.Node) bool { return c.Type == html.ElementNode && c.Data == "a" })
		if link != nil {
			e.Title = strings.TrimSpace(textContent(link))
			href, _ := getAttr(link, "href")
			e.URL = href
			e.FileName = extractFilenameFromHref(href)
			e.IsDeleted = false
		} else {
			// Deleted post, take raw text
			titleText := strings.TrimSpace(textContent(tn))
			e.Title = titleText
			e.IsDeleted = true
			e.FileName = "M.0.A.0"
			// Try to extract owner from [owner]
			if m := regexp.MustCompile(`\[(.*?)\]`).FindStringSubmatch(titleText); len(m) == 2 {
				e.Owner = m[1]
			}
		}
	}

	// meta
	if mn := findFirstChild(n, func(c *html.Node) bool { return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "meta") }); mn != nil {
		if an := findFirstChild(mn, func(c *html.Node) bool { return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "author") }); an != nil {
			author := strings.TrimSpace(textContent(an))
			if author != "" && author != "-" { // '-' is placeholder for deleted
				e.Owner = author
			}
		}
		if dn := findFirstChild(mn, func(c *html.Node) bool { return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "date") }); dn != nil {
			e.Date = strings.TrimSpace(textContent(dn))
		}
		if mk := findFirstChild(mn, func(c *html.Node) bool { return c.Type == html.ElementNode && c.Data == "div" && hasClass(c, "mark") }); mk != nil {
			e.Mark = strings.TrimSpace(textContent(mk))
		}
	}

	return e
}

func parseRecommend(s string) int {
	if s == "" {
		return 0
	}
	if s == "爆" {
		return 100
	}
	if strings.HasPrefix(s, "X") || strings.HasPrefix(s, "x") {
		if v, err := strconv.Atoi(strings.TrimPrefix(strings.ToUpper(s), "X")); err == nil {
			return -v
		}
		return 0
	}
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return 0
}

func extractFilenameFromHref(href string) string {
	// Expecting /bbs/<Board>/<Filename>.html
	base := path.Base(href)
	if strings.HasSuffix(base, ".html") {
		return strings.TrimSuffix(base, ".html")
	}
	return base
}

// Helpers
func getAttr(n *html.Node, key string) (string, bool) {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val, true
		}
	}
	return "", false
}

func hasClass(n *html.Node, class string) bool {
	if n == nil {
		return false
	}
	for _, a := range n.Attr {
		if a.Key == "class" {
			for _, c := range strings.Fields(a.Val) {
				if c == class {
					return true
				}
			}
		}
	}
	return false
}

func findFirst(root *html.Node, pred func(*html.Node) bool) *html.Node {
	if root == nil {
		return nil
	}
	for n := range root.Descendants() {
		if pred(n) {
			return n
		}
	}
	return nil
}

func findFirstChild(n *html.Node, pred func(*html.Node) bool) *html.Node {
	for c := range n.ChildNodes() {
		if pred(c) {
			return c
		}
	}
	return nil
}

func textContent(n *html.Node) string {
	var b strings.Builder
	var rec func(*html.Node)
	rec = func(x *html.Node) {
		switch x.Type {
		case html.TextNode:
			b.WriteString(x.Data)
		case html.ElementNode:
			for c := range x.ChildNodes() {
				rec(c)
			}
		}
	}
	rec(n)
	return b.String()
}
