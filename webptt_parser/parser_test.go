package webpttparser

import (
	"testing"
)

const sampleIndexHTML = `

<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		

<meta name="viewport" content="width=device-width, initial-scale=1">

<title>看板 SYSOP 文章列表 - 批踢踢實業坊</title>

<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-common.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-base.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-custom.css">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/pushstream.css" media="screen">
<link rel="stylesheet" type="text/css" href="//images.ptt.cc/bbs/v2.27/bbs-print.css" media="print">





	</head>
    <body>
		
<div id="topbar-container">
	<div id="topbar" class="bbs-content">
		<a id="logo" href="/bbs/">批踢踢實業坊</a>
		<span>&rsaquo;</span>
		<a class="board" href="/bbs/SYSOP/index.html"><span class="board-label">看板 </span>SYSOP</a>
		<a class="right small" href="/about.html">關於我們</a>
		<a class="right small" href="/contact.html">聯絡資訊</a>
	</div>
</div>

<div id="main-container">
	<div id="action-bar-container">
		<div class="action-bar">
			<div class="btn-group btn-group-dir">
				<a class="btn selected" href="/bbs/SYSOP/index.html">看板</a>
				<a class="btn" href="/man/SYSOP/index.html">精華區</a>
			</div>
			<div class="btn-group btn-group-paging">
				<a class="btn wide" href="/bbs/SYSOP/index1.html">最舊</a>
				<a class="btn wide" href="/bbs/SYSOP/index968.html">&lsaquo; 上頁</a>
				<a class="btn wide disabled">下頁 &rsaquo;</a>
				<a class="btn wide" href="/bbs/SYSOP/index.html">最新</a>
			</div>
		</div>
	</div>

	<div class="r-list-container action-bar-margin bbs-screen">
		<div class="search-bar">
			<form type="get" action="search" id="search-bar">
				<input class="query" type="text" name="q" value="" placeholder="搜尋文章&#x22ef;">
			</form>
		</div>

		
		
            
        
        
		<div class="r-ent">
			<div class="nrec"><span class="hl f2">1</span></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1757260299.A.E94.html">[情報] 42屆小天使招考（9/8～9/17）</a>
			
			</div>
			<div class="meta">
				<div class="author">hp6304</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E6%83%85%E5%A0%B1%5D&#43;42%E5%B1%86%E5%B0%8F%E5%A4%A9%E4%BD%BF%E6%8B%9B%E8%80%83%EF%BC%889%2F8%EF%BD%9E9%2F17%EF%BC%89">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3Ahp6304">搜尋看板內 hp6304 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date"> 9/07</div>
				<div class="mark">!</div>
			</div>
		</div>

		
            
        
        
		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1759419781.A.3A1.html">[問題] 手機已被使用請版主協助認證</a>
			
			</div>
			<div class="meta">
				<div class="author">pangmaomi</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E5%95%8F%E9%A1%8C%5D&#43;%E6%89%8B%E6%A9%9F%E5%B7%B2%E8%A2%AB%E4%BD%BF%E7%94%A8%E8%AB%8B%E7%89%88%E4%B8%BB%E5%8D%94%E5%8A%A9%E8%AA%8D%E8%AD%89">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3Apangmaomi">搜尋看板內 pangmaomi 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date">10/02</div>
				<div class="mark"></div>
			</div>
		</div>

		
            
        
        
		<div class="r-ent">
			<div class="nrec"><span class="hl f2">1</span></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1760274859.A.41A.html">[問題] 無法通過驗證該怎麼辦</a>
			
			</div>
			<div class="meta">
				<div class="author">cchyl</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E5%95%8F%E9%A1%8C%5D&#43;%E7%84%A1%E6%B3%95%E9%80%9A%E9%81%8E%E9%A9%97%E8%AD%89%E8%A9%B2%E6%80%8E%E9%BA%BC%E8%BE%A6">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3Acchyl">搜尋看板內 cchyl 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date">10/12</div>
				<div class="mark"></div>
			</div>
		</div>

		
            
        
        
		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1760425266.A.3B8.html">[問題] 無法收到驗證碼</a>
			
			</div>
			<div class="meta">
				<div class="author">honestrbb</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E5%95%8F%E9%A1%8C%5D&#43;%E7%84%A1%E6%B3%95%E6%94%B6%E5%88%B0%E9%A9%97%E8%AD%89%E7%A2%BC">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3Ahonestrbb">搜尋看板內 honestrbb 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date">10/14</div>
				<div class="mark"></div>
			</div>
		</div>

		
        
        <div class="r-list-sep"></div>
            
                
        
        
		<div class="r-ent">
			<div class="nrec"><span class="hl f2">6</span></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1126081374.A.D45.html">[公告] 網友需為言論負責，切莫誤觸法網</a>
			
			</div>
			<div class="meta">
				<div class="author">in2</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E5%85%AC%E5%91%8A%5D&#43;%E7%B6%B2%E5%8F%8B%E9%9C%80%E7%82%BA%E8%A8%80%E8%AB%96%E8%B2%A0%E8%B2%AC%EF%BC%8C%E5%88%87%E8%8E%AB%E8%AA%A4%E8%A7%B8%E6%B3%95%E7%B6%B2">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3Ain2">搜尋看板內 in2 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date"> 9/07</div>
				<div class="mark">!</div>
			</div>
		</div>

            
                
        
        
		<div class="r-ent">
			<div class="nrec"><span class="hl f1">爆</span></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1290228854.A.352.html">[回報] 斷線情形回報推文專用</a>
			
			</div>
			<div class="meta">
				<div class="author">wens</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E5%9B%9E%E5%A0%B1%5D&#43;%E6%96%B7%E7%B7%9A%E6%83%85%E5%BD%A2%E5%9B%9E%E5%A0%B1%E6%8E%A8%E6%96%87%E5%B0%88%E7%94%A8">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3Awens">搜尋看板內 wens 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date">11/20</div>
				<div class="mark">S</div>
			</div>
		</div>

            
                
        
        
		<div class="r-ent">
			<div class="nrec"></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1563944454.A.B6D.html">[公告] 請使用者多加注意我國保護兒少的法令</a>
			
			</div>
			<div class="meta">
				<div class="author">longbow2</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E5%85%AC%E5%91%8A%5D&#43;%E8%AB%8B%E4%BD%BF%E7%94%A8%E8%80%85%E5%A4%9A%E5%8A%A0%E6%B3%A8%E6%84%8F%E6%88%91%E5%9C%8B%E4%BF%9D%E8%AD%B7%E5%85%92%E5%B0%91%E7%9A%84%E6%B3%95%E4%BB%A4">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3Alongbow2">搜尋看板內 longbow2 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date"> 7/24</div>
				<div class="mark">!</div>
			</div>
		</div>

            
                
        
        
		<div class="r-ent">
			<div class="nrec"><span class="hl f1">爆</span></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1627279813.A.1C8.html">[公告] 如何以 AOTP 手機驗證註冊本站帳號教學</a>
			
			</div>
			<div class="meta">
				<div class="author">s5048218</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E5%85%AC%E5%91%8A%5D&#43;%E5%A6%82%E4%BD%95%E4%BB%A5&#43;AOTP&#43;%E6%89%8B%E6%A9%9F%E9%A9%97%E8%AD%89%E8%A8%BB%E5%86%8A%E6%9C%AC%E7%AB%99%E5%B8%B3%E8%99%9F%E6%95%99%E5%AD%B8">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3As5048218">搜尋看板內 s5048218 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date"> 7/26</div>
				<div class="mark">M</div>
			</div>
		</div>

            
                
        
        
		<div class="r-ent">
			<div class="nrec"><span class="hl f1">爆</span></div>
			<div class="title">
			
				<a href="/bbs/SYSOP/M.1650509751.A.2F9.html">Fw: [公告] 帳號被盜洽ID_GetBack；復權洽 ID_Problem</a>
			
			</div>
			<div class="meta">
				<div class="author">Makotoyen</div>
				<div class="article-menu">
					
					<div class="trigger">&#x22ef;</div>
					<div class="dropdown">
						<div class="item"><a href="/bbs/SYSOP/search?q=thread%3A%5B%E5%85%AC%E5%91%8A%5D&#43;%E5%B8%B3%E8%99%9F%E8%A2%AB%E7%9B%9C%E6%B4%BDID_GetBack%EF%BC%9B%E5%BE%A9%E6%AC%8A%E6%B4%BD&#43;ID_Problem">搜尋同標題文章</a></div>
						
						<div class="item"><a href="/bbs/SYSOP/search?q=author%3AMakotoyen">搜尋看板內 Makotoyen 的文章</a></div>
						
					</div>
					
				</div>
				<div class="date"> 4/21</div>
				<div class="mark">M</div>
			</div>
		</div>

            
        
	</div>

    
</div>

		



<script async src="https://www.googletagmanager.com/gtag/js?id=G-DZ6Y3BY9GW"></script>
<script>
      window.dataLayer = window.dataLayer || [];
      function gtag(){dataLayer.push(arguments);}
      gtag('js', new Date());

      gtag('config', 'G-DZ6Y3BY9GW');
</script>
<script>
  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
  })(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

  ga('create', 'UA-32365737-1', {
    cookieDomain: 'ptt.cc',
    legacyCookieDomain: 'ptt.cc'
  });
  ga('send', 'pageview');
</script>


		
<script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
<script src="//images.ptt.cc/bbs/v2.27/bbs.js"></script>

    <script defer src="https://static.cloudflareinsights.com/beacon.min.js/vcd15cbe7772f49c399c6a5babf22c1241717689176015" integrity="sha512-ZpsOmlRQV6y907TI0dKBHq9Md29nnaEIPlkf84rnaERnq6zvWvPUqr2ft8M1aS28oN72PdrCzSjY4U6VaAw1EQ==" data-cf-beacon='{"version":"2024.11.0","token":"515615eb5fab4c9b91a11e9bf529e6cf","r":1,"server_timing":{"name":{"cfCacheStatus":true,"cfEdge":true,"cfExtPri":true,"cfL4":true,"cfOrigin":true,"cfSpeedBrain":true},"location_startswith":null}}' crossorigin="anonymous"></script>
</body>
</html>
`

func TestParseBoardIndexPage(t *testing.T) {
	page, err := ParseBoardIndexPage([]byte(sampleIndexHTML))
	if err != nil {
		t.Fatalf("ParseBoardIndexPage returned error: %v", err)
	}
	if page.FirstPage != "/bbs/SYSOP/index1.html" {
		t.Fatalf("FirstPage = %q", page.FirstPage)
	}
	if page.PrevPage != "/bbs/SYSOP/index968.html" {
		t.Fatalf("PrevPage = %q", page.PrevPage)
	}
	if page.NextPage != "" {
		t.Fatalf("NextPage should be empty, got %q", page.NextPage)
	}
	if page.LastPage != "/bbs/SYSOP/index.html" {
		t.Fatalf("LastPage = %q", page.LastPage)
	}

	if len(page.Articles) != 4 {
		t.Fatalf("expected 4 normal articles before separator, got %d", len(page.Articles))
	}
	if len(page.Bottoms) != 5 {
		t.Fatalf("expected 5 bottom (sticky) articles after separator, got %d", len(page.Bottoms))
	}

	a0 := page.Articles[0]
	if a0.Owner != "hp6304" || a0.FileName != "M.1757260299.A.E94" || a0.URL == "" || a0.Recommend != 1 || a0.Mark != "!" {
		t.Fatalf("unexpected first entry: %+v", a0)
	}
	if a0.Title == "" || a0.Date == "" {
		t.Fatalf("first entry should have title and date, got: %+v", a0)
	}

	a1 := page.Articles[1]
	if a1.Owner != "pangmaomi" || a1.FileName != "M.1759419781.A.3A1" || a1.Recommend != 0 {
		t.Fatalf("unexpected second entry: %+v", a1)
	}

	b0 := page.Bottoms[0]
	if b0.Owner != "in2" || b0.FileName != "M.1126081374.A.D45" || b0.Recommend != 6 || b0.Mark != "!" {
		t.Fatalf("unexpected first bottom entry: %+v", b0)
	}

	b1 := page.Bottoms[1]
	if b1.Owner != "wens" || b1.Recommend != 100 || b1.Mark != "S" || b1.FileName != "M.1290228854.A.352" {
		t.Fatalf("unexpected second bottom entry (爆推 S): %+v", b1)
	}

	b3 := page.Bottoms[3]
	if b3.Owner != "s5048218" || b3.Recommend != 100 || b3.Mark != "M" || b3.FileName != "M.1627279813.A.1C8" {
		t.Fatalf("unexpected fourth bottom entry (爆推 M): %+v", b3)
	}
}
