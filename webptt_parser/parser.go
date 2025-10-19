package webpttparser

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

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
