package webpttparser

import (
	"strings"
	"testing"
)

const sampleHotboardsHTML = `

<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		

<meta name="viewport" content="width=device-width, initial-scale=1">

<title>熱門看板 - 批踢踢實業坊</title>

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
		<a class="right small" href="/about.html">關於我們</a>
		<a class="right small" href="/contact.html">聯絡資訊</a>
	</div>
</div>

<div id="main-container">
	<div id="action-bar-container">
		<div class="action-bar">
			<div class="btn-group btn-group-cls">
				<a class="btn selected" href="/bbs/hotboards.html">熱門看板</a>
                <a class="btn" href="/cls/1">分類看板</a>
			</div>
		</div>
	</div>

	<div class="b-list-container action-bar-margin bbs-screen">
        
        <div class="b-ent">
            <a class="board" href="/bbs/Gossiping/index.html">
                <div class="board-name">Gossiping</div>
                <div class="board-nuser"><span class="hl f4">6606</span></div>
                <div class="board-class">綜合</div>
                <div class="board-title">&#9678;[八翻] 10/24~26臺灣光復節連假♬</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Baseball/index.html">
                <div class="board-name">Baseball</div>
                <div class="board-nuser"><span class="hl f4">5056</span></div>
                <div class="board-class">棒球</div>
                <div class="board-title">&#9678;[棒球] TS猿爪巨蛋熱戰游刃搏終耀</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/LoL/index.html">
                <div class="board-name">LoL</div>
                <div class="board-nuser"><span class="hl f1">3686</span></div>
                <div class="board-class">遊戲</div>
                <div class="board-title">&#9678;[LoL] CFO獎勵一把HLE / PSG淘汰</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/C_Chat/index.html">
                <div class="board-name">C_Chat</div>
                <div class="board-nuser"><span class="hl f1">3038</span></div>
                <div class="board-class">閒談</div>
                <div class="board-title">&#9678;[希洽] 這裡是希洽閒聊板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Stock/index.html">
                <div class="board-name">Stock</div>
                <div class="board-nuser"><span class="hl f1">2592</span></div>
                <div class="board-class">學術</div>
                <div class="board-title">&#9678;[股票] 股票板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/NBA/index.html">
                <div class="board-name">NBA</div>
                <div class="board-nuser"><span class="hl">918</span></div>
                <div class="board-class">NBA.</div>
                <div class="board-title">&#9678;[NBA] 熱身賽</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Elephants/index.html">
                <div class="board-name">Elephants</div>
                <div class="board-nuser"><span class="hl">917</span></div>
                <div class="board-class">CPBL</div>
                <div class="board-title">&#9678;[兄弟] 【進化‧決勝】</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Lifeismoney/index.html">
                <div class="board-name">Lifeismoney</div>
                <div class="board-nuser"><span class="hl">840</span></div>
                <div class="board-class">省錢</div>
                <div class="board-title">&#9678;[省錢] 省錢板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/HatePolitics/index.html">
                <div class="board-name">HatePolitics</div>
                <div class="board-nuser"><span class="hl">576</span></div>
                <div class="board-class">Hate</div>
                <div class="board-title">&#9678;[政黑] 板主投票徵選中</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/KoreaStar/index.html">
                <div class="board-name">KoreaStar</div>
                <div class="board-nuser"><span class="hl">497</span></div>
                <div class="board-class">韓國</div>
                <div class="board-title">&#9678;[韓星] 1017 BM Rami㊣ HBD!! </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/car/index.html">
                <div class="board-name">car</div>
                <div class="board-nuser"><span class="hl">464</span></div>
                <div class="board-class">車車</div>
                <div class="board-title">&#9678;[汽車] 恭喜lll156k1529當選</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/home-sale/index.html">
                <div class="board-name">home-sale</div>
                <div class="board-nuser"><span class="hl">362</span></div>
                <div class="board-class">房屋</div>
                <div class="board-title">&#9678;[房版] 請詳閱新板規避免違規</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Beauty/index.html">
                <div class="board-name">Beauty</div>
                <div class="board-nuser"><span class="hl">322</span></div>
                <div class="board-class">聊天</div>
                <div class="board-title">&#9678;[表特] 貼AI圖 一律水桶+退文</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/SportLottery/index.html">
                <div class="board-name">SportLottery</div>
                <div class="board-nuser"><span class="hl">321</span></div>
                <div class="board-class">博弈</div>
                <div class="board-title">&#9678;[運彩] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/basketballTW/index.html">
                <div class="board-name">basketballTW</div>
                <div class="board-nuser"><span class="hl">315</span></div>
                <div class="board-class">籃球</div>
                <div class="board-title">&#9678;[台籃] 例行賽</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/MobileComm/index.html">
                <div class="board-name">MobileComm</div>
                <div class="board-nuser"><span class="hl">304</span></div>
                <div class="board-class">資訊</div>
                <div class="board-title">&#9678;[通訊] 手機o平板o資費 請進!</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Japan_Travel/index.html">
                <div class="board-name">Japan_Travel</div>
                <div class="board-nuser"><span class="hl">285</span></div>
                <div class="board-class">旅遊</div>
                <div class="board-title">&#9678;前往山區小心野生動物</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/movie/index.html">
                <div class="board-name">movie</div>
                <div class="board-nuser"><span class="hl">266</span></div>
                <div class="board-class">綜合</div>
                <div class="board-title">&#9678;[電影] 取消票房文規範</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/PC_Shopping/index.html">
                <div class="board-name">PC_Shopping</div>
                <div class="board-nuser"><span class="hl">261</span></div>
                <div class="board-class">硬體</div>
                <div class="board-title">&#9678;[電蝦]全台最大客服中心</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/sex/index.html">
                <div class="board-name">sex</div>
                <div class="board-nuser"><span class="hl">254</span></div>
                <div class="board-class">男女</div>
                <div class="board-title">&#9678;[西斯] 愛對了人 情人節每天都過</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Steam/index.html">
                <div class="board-name">Steam</div>
                <div class="board-nuser"><span class="hl">244</span></div>
                <div class="board-class">平台</div>
                <div class="board-title">&#9678;本板禁止買賣合購代購代刷</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/iOS/index.html">
                <div class="board-name">iOS</div>
                <div class="board-nuser"><span class="hl">236</span></div>
                <div class="board-class">系統</div>
                <div class="board-title">&#9678;[iOS] Let Loose 超...燃</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Monkeys/index.html">
                <div class="board-name">Monkeys</div>
                <div class="board-nuser"><span class="hl">212</span></div>
                <div class="board-class">CPBL</div>
                <div class="board-title">&#9678;[桃猿] 桃猿制霸!!!! TS 0:2領先</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/japanavgirls/index.html">
                <div class="board-name">japanavgirls</div>
                <div class="board-nuser"><span class="hl">200</span></div>
                <div class="board-class">綜合</div>
                <div class="board-title">&#9678;AV女優板 清新、優質,神人照格式</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Tech_Job/index.html">
                <div class="board-name">Tech_Job</div>
                <div class="board-nuser"><span class="hl">196</span></div>
                <div class="board-class">工作</div>
                <div class="board-title">&#9678;[科技] B0988698088水桶十年</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/creditcard/index.html">
                <div class="board-name">creditcard</div>
                <div class="board-nuser"><span class="hl">178</span></div>
                <div class="board-class">理財</div>
                <div class="board-title">&#9678;[卡板] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Military/index.html">
                <div class="board-name">Military</div>
                <div class="board-nuser"><span class="hl">174</span></div>
                <div class="board-class">軍事</div>
                <div class="board-title">&#9678;[軍事]理性討論軍事 請詳閱板規</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/FORMULA1/index.html">
                <div class="board-name">FORMULA1</div>
                <div class="board-nuser"><span class="hl">173</span></div>
                <div class="board-class">賽車</div>
                <div class="board-title">&#9678;F1板 </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/KoreaDrama/index.html">
                <div class="board-name">KoreaDrama</div>
                <div class="board-nuser"><span class="hl">159</span></div>
                <div class="board-class">韓劇</div>
                <div class="board-title">&#9678;[韓劇] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Marginalman/index.html">
                <div class="board-name">Marginalman</div>
                <div class="board-nuser"><span class="hl">157</span></div>
                <div class="board-class">心情</div>
                <div class="board-title">&#9678;[邊緣] CFO vs HLE</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/HardwareSale/index.html">
                <div class="board-name">HardwareSale</div>
                <div class="board-nuser"><span class="hl">146</span></div>
                <div class="board-class">買賣</div>
                <div class="board-title">&#9678;[硬交] 詐騙會從內文+LINE找上門</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/China-Drama/index.html">
                <div class="board-name">China-Drama</div>
                <div class="board-nuser"><span class="hl">144</span></div>
                <div class="board-class">中劇</div>
                <div class="board-title">&#9678;[中劇] 讓空降來得更猛烈些吧！</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/DIABLO/index.html">
                <div class="board-name">DIABLO</div>
                <div class="board-nuser"><span class="hl">140</span></div>
                <div class="board-class">線上</div>
                <div class="board-title">&#9678;[暗黑] D4 PTR 10/21~10/28</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Tainan/index.html">
                <div class="board-name">Tainan</div>
                <div class="board-nuser"><span class="hl">137</span></div>
                <div class="board-class">台南</div>
                <div class="board-title">&#9678;[台南] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/TY_Research/index.html">
                <div class="board-name">TY_Research</div>
                <div class="board-nuser"><span class="hl">134</span></div>
                <div class="board-class">大氣</div>
                <div class="board-title">&#9678;大氣科學板23週年板慶</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Kaohsiung/index.html">
                <div class="board-name">Kaohsiung</div>
                <div class="board-nuser"><span class="hl">131</span></div>
                <div class="board-class">高雄</div>
                <div class="board-title">&#9678;[高雄] 翅膀硬惹</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Japandrama/index.html">
                <div class="board-name">Japandrama</div>
                <div class="board-nuser"><span class="hl">122</span></div>
                <div class="board-class">日劇</div>
                <div class="board-title">&#9678;[日劇] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/PlayStation/index.html">
                <div class="board-name">PlayStation</div>
                <div class="board-nuser"><span class="hl">117</span></div>
                <div class="board-class">主機</div>
                <div class="board-title">&#9678;[PS5] 10/21 忍者外傳4</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/MacShop/index.html">
                <div class="board-name">MacShop</div>
                <div class="board-nuser"><span class="hl">108</span></div>
                <div class="board-class">資訊</div>
                <div class="board-title">&#9678;[麥蝦] 禁止販售非現貨商品</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/BabyMother/index.html">
                <div class="board-name">BabyMother</div>
                <div class="board-nuser"><span class="hl">108</span></div>
                <div class="board-class">家庭</div>
                <div class="board-title">&#9678;[寶寶] 記得打流感疫苗</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/miHoYo/index.html">
                <div class="board-name">miHoYo</div>
                <div class="board-nuser"><span class="hl">101</span></div>
                <div class="board-class">米哈</div>
                <div class="board-title">&#9678;[米哈遊] 鐵道要用行動網路開</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/TaichungBun/index.html">
                <div class="board-name">TaichungBun</div>
                <div class="board-nuser"><span class="hl">101</span></div>
                <div class="board-class">台中</div>
                <div class="board-title">&#9678;[台中] 新板規施行過度期至10/24</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/AC_In/index.html">
                <div class="board-name">AC_In</div>
                <div class="board-nuser"><span class="hl">100</span></div>
                <div class="board-class">閒談</div>
                <div class="board-title">&#9678;[裏洽]</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/mobilesales/index.html">
                <div class="board-name">mobilesales</div>
                <div class="board-nuser"><span class="hl">96</span></div>
                <div class="board-class">資訊</div>
                <div class="board-title">&#9678;[Mobile] 遵守版規，進桶看通知信</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/CFantasy/index.html">
                <div class="board-name">CFantasy</div>
                <div class="board-nuser"><span class="hl">95</span></div>
                <div class="board-class">玄幻</div>
                <div class="board-title">&#9678;[玄幻] 玄幻小說板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/biker/index.html">
                <div class="board-name">biker</div>
                <div class="board-nuser"><span class="hl">93</span></div>
                <div class="board-class">車車</div>
                <div class="board-title">&#9678;[機車]光陽和他的網軍</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Aviation/index.html">
                <div class="board-name">Aviation</div>
                <div class="board-nuser"><span class="hl">91</span></div>
                <div class="board-class">交通</div>
                <div class="board-title">&#9678;[航空] 注意行動電源上機新規定</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/NSwitch/index.html">
                <div class="board-name">NSwitch</div>
                <div class="board-nuser"><span class="hl">90</span></div>
                <div class="board-class">主機</div>
                <div class="board-title">&#9678;[NS] 寶可夢傳說 Z-A 10/16</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/WomenTalk/index.html">
                <div class="board-name">WomenTalk</div>
                <div class="board-nuser"><span class="hl">88</span></div>
                <div class="board-class">聊天</div>
                <div class="board-title">&#9678;[女孩] 天氣熱小心中暑&lt;(@o@)/</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Gamesale/index.html">
                <div class="board-name">Gamesale</div>
                <div class="board-nuser"><span class="hl">87</span></div>
                <div class="board-class">綜合</div>
                <div class="board-title">&#9678;[GS] 遊戲交易板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/TaiwanDrama/index.html">
                <div class="board-name">TaiwanDrama</div>
                <div class="board-nuser"><span class="hl">85</span></div>
                <div class="board-class">臺劇</div>
                <div class="board-title">&#9678;[臺劇] 二代咖啡館的化外接班人</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Golden-Award/index.html">
                <div class="board-name">Golden-Award</div>
                <div class="board-nuser"><span class="hl">82</span></div>
                <div class="board-class">三金</div>
                <div class="board-title">&#9678;[三金] 金鐘60 戲劇節目獎 影后</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/TW_Entertain/index.html">
                <div class="board-name">TW_Entertain</div>
                <div class="board-nuser"><span class="hl">82</span></div>
                <div class="board-class">綜藝</div>
                <div class="board-title">&#9678;[台綜] 恭喜LULU，漢典</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/CarShop/index.html">
                <div class="board-name">CarShop</div>
                <div class="board-nuser"><span class="hl">78</span></div>
                <div class="board-class">買賣</div>
                <div class="board-title">&#9678;標題沒分類會直接被刪文</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/SakaTalk/index.html">
                <div class="board-name">SakaTalk</div>
                <div class="board-nuser"><span class="hl">76</span></div>
                <div class="board-class">日本</div>
                <div class="board-title">&#9678;坂道閒聊板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/AllTogether/index.html">
                <div class="board-name">AllTogether</div>
                <div class="board-nuser"><span class="hl">74</span></div>
                <div class="board-class">聯誼</div>
                <div class="board-title">&#9678;[歐兔] 祝 交得球友</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Badminton/index.html">
                <div class="board-name">Badminton</div>
                <div class="board-nuser"><span class="hl">74</span></div>
                <div class="board-class">羽球</div>
                <div class="board-title">&#9678;[羽球版]</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/CVS/index.html">
                <div class="board-name">CVS</div>
                <div class="board-nuser"><span class="hl">71</span></div>
                <div class="board-class">資訊</div>
                <div class="board-title">&#9678;[CVS] 超商店員a.k.a.superman</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/watch/index.html">
                <div class="board-name">watch</div>
                <div class="board-nuser"><span class="hl">70</span></div>
                <div class="board-class">購一</div>
                <div class="board-title">&#9678;[錶板] Show me your Watch!</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/e-coupon/index.html">
                <div class="board-name">e-coupon</div>
                <div class="board-nuser"><span class="hl">64</span></div>
                <div class="board-class">理財</div>
                <div class="board-title">&#9678;新板主2025年4月25日上任</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/EAseries/index.html">
                <div class="board-name">EAseries</div>
                <div class="board-nuser"><span class="hl">63</span></div>
                <div class="board-class">歐美</div>
                <div class="board-title">&#9678;[EA] 歐美影集版 </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Boy-Girl/index.html">
                <div class="board-name">Boy-Girl</div>
                <div class="board-nuser"><span class="hl">59</span></div>
                <div class="board-class">心情</div>
                <div class="board-title">&#9678;[男女] 歡迎光臨男女板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/MLB/index.html">
                <div class="board-name">MLB</div>
                <div class="board-nuser"><span class="hl">58</span></div>
                <div class="board-class">#MLB</div>
                <div class="board-title">&#9678;[MLB] 水手3-2藍鳥 道奇晉級</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Tennis/index.html">
                <div class="board-name">Tennis</div>
                <div class="board-nuser"><span class="hl">56</span></div>
                <div class="board-class">網球</div>
                <div class="board-title">&#9678;[網球] 辛卡辛卡 鑰CoGa芭</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/KR_Entertain/index.html">
                <div class="board-name">KR_Entertain</div>
                <div class="board-nuser"><span class="hl">55</span></div>
                <div class="board-class">綜藝</div>
                <div class="board-title">&#9678;[韓綜] Show Me The 新板主!</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Headphone/index.html">
                <div class="board-name">Headphone</div>
                <div class="board-nuser"><span class="hl">52</span></div>
                <div class="board-class">資訊</div>
                <div class="board-title">&#9678;[耳機] 歡迎大家多發心得~</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/TWICE/index.html">
                <div class="board-name">TWICE</div>
                <div class="board-nuser"><span class="hl f3">50</span></div>
                <div class="board-class">韓國</div>
                <div class="board-title">&#9678;[TWICE]TWICE : ONE IN A MILL10N</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/DigiCurrency/index.html">
                <div class="board-name">DigiCurrency</div>
                <div class="board-nuser"><span class="hl f3">50</span></div>
                <div class="board-class">資訊</div>
                <div class="board-title">&#9678;[數位貨幣] real SoV will stay</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/YuanChuang/index.html">
                <div class="board-name">YuanChuang</div>
                <div class="board-nuser"><span class="hl f3">49</span></div>
                <div class="board-class">原創</div>
                <div class="board-title">&#9678;[原創] 連假就要看小說歐耶～</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/FAPL/index.html">
                <div class="board-name">FAPL</div>
                <div class="board-nuser"><span class="hl f3">48</span></div>
                <div class="board-class">聯賽</div>
                <div class="board-title">&#9678;歡樂．三喵．足球</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/marriage/index.html">
                <div class="board-name">marriage</div>
                <div class="board-nuser"><span class="hl f3">47</span></div>
                <div class="board-class">婚姻</div>
                <div class="board-title">&#9678;[婚姻] 歡迎光臨婚姻板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/CheerGirlsTW/index.html">
                <div class="board-name">CheerGirlsTW</div>
                <div class="board-nuser"><span class="hl f3">46</span></div>
                <div class="board-class">應援</div>
                <div class="board-title">&#9678;[真香] 阿金請勿發表不當言論</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Drama-Ticket/index.html">
                <div class="board-name">Drama-Ticket</div>
                <div class="board-nuser"><span class="hl f3">46</span></div>
                <div class="board-class">舞臺</div>
                <div class="board-title">&#9678;藝文票券轉售板 禁止徵票小心詐騙</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/DC_SALE/index.html">
                <div class="board-name">DC_SALE</div>
                <div class="board-nuser"><span class="hl f3">45</span></div>
                <div class="board-class">攝影</div>
                <div class="board-title">&#9678;詐騙請去至底文回覆</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Hsinchu/index.html">
                <div class="board-name">Hsinchu</div>
                <div class="board-nuser"><span class="hl f3">45</span></div>
                <div class="board-class">新竹</div>
                <div class="board-title">&#9678;[新竹] 新版規研議中~~\⊙▽⊙/</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/e-shopping/index.html">
                <div class="board-name">e-shopping</div>
                <div class="board-nuser"><span class="hl f3">43</span></div>
                <div class="board-class">網購</div>
                <div class="board-title">&#9678;[ＥＳ] 注意板規 遠離水桶</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Tobacco/index.html">
                <div class="board-name">Tobacco</div>
                <div class="board-nuser"><span class="hl f3">42</span></div>
                <div class="board-class">生一</div>
                <div class="board-title">&#9678;不是過了就可以上架 不是這樣喔</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/PokemonGO/index.html">
                <div class="board-name">PokemonGO</div>
                <div class="board-nuser"><span class="hl f3">42</span></div>
                <div class="board-class">抓寶</div>
                <div class="board-title">&#9678;[PMGO] 2025萬聖節活動＆GO集章趣</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/give/index.html">
                <div class="board-name">give</div>
                <div class="board-nuser"><span class="hl f3">41</span></div>
                <div class="board-class">贈送</div>
                <div class="board-title">&#9678;[贈送] 標題不合置底文規定  刪</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/H-GAME/index.html">
                <div class="board-name">H-GAME</div>
                <div class="board-nuser"><span class="hl f3">41</span></div>
                <div class="board-class">綜合</div>
                <div class="board-title">&#9678;H-GAME板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/MuscleBeach/index.html">
                <div class="board-name">MuscleBeach</div>
                <div class="board-nuser"><span class="hl f3">41</span></div>
                <div class="board-class">美體</div>
                <div class="board-title">&#9678;[健身] 發文前請務必詳閱板規</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/BabyProducts/index.html">
                <div class="board-name">BabyProducts</div>
                <div class="board-nuser"><span class="hl f3">40</span></div>
                <div class="board-class">買賣</div>
                <div class="board-title">&#9678;防詐騙，交易異常時，自行找客服</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/joke/index.html">
                <div class="board-name">joke</div>
                <div class="board-nuser"><span class="hl f3">39</span></div>
                <div class="board-class">娛樂</div>
                <div class="board-title">&#9678;[就可] 為花蓮祈禱</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/FATE_GO/index.html">
                <div class="board-name">FATE_GO</div>
                <div class="board-nuser"><span class="hl f3">36</span></div>
                <div class="board-class">手遊</div>
                <div class="board-title">&#9678;[FGO] 日 殺職冠位戴冠戰 開幕</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/marvel/index.html">
                <div class="board-name">marvel</div>
                <div class="board-nuser"><span class="hl f3">36</span></div>
                <div class="board-class">生二</div>
                <div class="board-title">&#9678;[媽佛] 你吃了誰煮的麵</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/DSLR/index.html">
                <div class="board-name">DSLR</div>
                <div class="board-nuser"><span class="hl f3">34</span></div>
                <div class="board-class">攝影</div>
                <div class="board-title">&#9678;[攝影] 2025</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Wanted/index.html">
                <div class="board-name">Wanted</div>
                <div class="board-nuser"><span class="hl f3">34</span></div>
                <div class="board-class">徵求</div>
                <div class="board-title">&#9678;[汪踢] 10/24~26光復節連假</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/WOW/index.html">
                <div class="board-name">WOW</div>
                <div class="board-nuser"><span class="hl f3">34</span></div>
                <div class="board-class">線上</div>
                <div class="board-title">&#9678;[WoW] 軍團REMIX/11.2.5 10/9上線</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Lakers/index.html">
                <div class="board-name">Lakers</div>
                <div class="board-nuser"><span class="hl f3">34</span></div>
                <div class="board-class">W-PA</div>
                <div class="board-title">&#9678;[Lakers] Offseason</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Bank_Service/index.html">
                <div class="board-name">Bank_Service</div>
                <div class="board-nuser"><span class="hl f3">34</span></div>
                <div class="board-class">銀行</div>
                <div class="board-title">&#9678;[銀行服務板] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/E-appliance/index.html">
                <div class="board-name">E-appliance</div>
                <div class="board-nuser"><span class="hl f3">33</span></div>
                <div class="board-class">資訊</div>
                <div class="board-title">&#9678;[家電] 發文請照規定發文</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/kartrider/index.html">
                <div class="board-name">kartrider</div>
                <div class="board-nuser"><span class="hl f3">32</span></div>
                <div class="board-class">線上</div>
                <div class="board-title">&#9678;跑跑卡丁車：飄移 10/16結束營運</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/PublicServan/index.html">
                <div class="board-name">PublicServan</div>
                <div class="board-nuser"><span class="hl f3">32</span></div>
                <div class="board-class">公職</div>
                <div class="board-title">&#9678;[公職] 記得錄影 好好保護自己</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Lions/index.html">
                <div class="board-name">Lions</div>
                <div class="board-nuser"><span class="hl f3">32</span></div>
                <div class="board-class">CPBL</div>
                <div class="board-title">&#9678;[獅隊] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Soft_Job/index.html">
                <div class="board-name">Soft_Job</div>
                <div class="board-nuser"><span class="hl f3">31</span></div>
                <div class="board-class">工作</div>
                <div class="board-title">&#9678;[軟工] 成為全宇宙のNo.1吧</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/nb-shopping/index.html">
                <div class="board-name">nb-shopping</div>
                <div class="board-nuser"><span class="hl f3">31</span></div>
                <div class="board-class">硬體</div>
                <div class="board-title">&#9678;[筆電蝦] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Key_Mou_Pad/index.html">
                <div class="board-name">Key_Mou_Pad</div>
                <div class="board-nuser"><span class="hl f3">31</span></div>
                <div class="board-class">生活</div>
                <div class="board-title">&#9678;[鍵鼠] MxMaster4滑鼠不附線</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/fastfood/index.html">
                <div class="board-name">fastfood</div>
                <div class="board-nuser"><span class="hl f3">31</span></div>
                <div class="board-class">美食</div>
                <div class="board-title">&#9678;[速食] 慟！麥香魚起士只放半片</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/CPBL_ticket/index.html">
                <div class="board-name">CPBL_ticket</div>
                <div class="board-nuser"><span class="hl f3">30</span></div>
                <div class="board-class">球聚</div>
                <div class="board-title">&#9678;票券版-進版請先看置底版規及公告</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Road_Running/index.html">
                <div class="board-name">Road_Running</div>
                <div class="board-nuser"><span class="hl f3">30</span></div>
                <div class="board-class">跑步</div>
                <div class="board-title">&#9678;跑步版-天氣熱多喝水</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/ONE_PIECE/index.html">
                <div class="board-name">ONE_PIECE</div>
                <div class="board-nuser"><span class="hl f3">30</span></div>
                <div class="board-class">日本</div>
                <div class="board-title">&#9678;[海賊王] 1160 《神之谷事件》</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/PokeMon/index.html">
                <div class="board-name">PokeMon</div>
                <div class="board-nuser"><span class="hl f3">30</span></div>
                <div class="board-class">日本</div>
                <div class="board-title">&#9678;[PM] 寶可夢Z-A 1016發售！</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/BeautySalon/index.html">
                <div class="board-name">BeautySalon</div>
                <div class="board-nuser"><span class="hl f3">30</span></div>
                <div class="board-class">美容</div>
                <div class="board-title">&#9678;[美保] 禁止洗文章數</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Salary/index.html">
                <div class="board-name">Salary</div>
                <div class="board-nuser"><span class="hl f3">29</span></div>
                <div class="board-class">職場</div>
                <div class="board-title">&#9678;[職場] 工作職場板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Taoyuan/index.html">
                <div class="board-name">Taoyuan</div>
                <div class="board-nuser"><span class="hl f3">29</span></div>
                <div class="board-class">桃園</div>
                <div class="board-title">&#9678;[桃園] </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Audiophile/index.html">
                <div class="board-name">Audiophile</div>
                <div class="board-nuser"><span class="hl f3">26</span></div>
                <div class="board-class">資訊</div>
                <div class="board-title">&#9678;[喇音] 注意 Line 加好友詐騙</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/ToS/index.html">
                <div class="board-name">ToS</div>
                <div class="board-nuser"><span class="hl f3">26</span></div>
                <div class="board-class">轉珠</div>
                <div class="board-title">&#9678;[神魔] 七王子</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/MobilePay/index.html">
                <div class="board-name">MobilePay</div>
                <div class="board-nuser"><span class="hl f3">26</span></div>
                <div class="board-class">理財</div>
                <div class="board-title">&#9678;刷卡回饋問題請至信用卡板</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Hearthstone/index.html">
                <div class="board-name">Hearthstone</div>
                <div class="board-nuser"><span class="hl f3">24</span></div>
                <div class="board-class">線上</div>
                <div class="board-title">&#9678;[爐石] 賀 西陵珩 獲得夏季賽亞軍</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/JP_Entertain/index.html">
                <div class="board-name">JP_Entertain</div>
                <div class="board-nuser"><span class="hl f3">24</span></div>
                <div class="board-class">綜藝</div>
                <div class="board-title">&#9678;[日綜] 歡迎來日綜 (づ′▽)づ</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/KanColle/index.html">
                <div class="board-name">KanColle</div>
                <div class="board-nuser"><span class="hl f3">23</span></div>
                <div class="board-class">艦娘</div>
                <div class="board-title">&#9678;[艦娘] 瑞鶴居然又有新立繪了</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Broad_Band/index.html">
                <div class="board-name">Broad_Band</div>
                <div class="board-nuser"><span class="hl f3">23</span></div>
                <div class="board-class">網路</div>
                <div class="board-title">&#9678;[寬頻板] 禁發廣告文/轉讓文</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/MakeUp/index.html">
                <div class="board-name">MakeUp</div>
                <div class="board-nuser"><span class="hl f3">23</span></div>
                <div class="board-class">美容</div>
                <div class="board-title">&#9678;[美妝] 10月怎麼還可以那麼熱</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/PuzzleDragon/index.html">
                <div class="board-name">PuzzleDragon</div>
                <div class="board-nuser"><span class="hl f3">22</span></div>
                <div class="board-class">轉珠</div>
                <div class="board-title">&#9678;[P＆D] 萬聖節檔期 </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/KoreanPop/index.html">
                <div class="board-name">KoreanPop</div>
                <div class="board-nuser"><span class="hl f3">22</span></div>
                <div class="board-class">音樂</div>
                <div class="board-title">&#9678;抵制買榜 抵制作弊，韓樂需要正義</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Railway/index.html">
                <div class="board-name">Railway</div>
                <div class="board-nuser"><span class="hl f3">22</span></div>
                <div class="board-class">交通</div>
                <div class="board-title">&#9678;[鐵道] 丟錢 說去台中 票丟給你 </div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Digitalhome/index.html">
                <div class="board-name">Digitalhome</div>
                <div class="board-nuser"><span class="hl f3">22</span></div>
                <div class="board-class">數位</div>
                <div class="board-title">&#9678;[電視] 交易請貨到付款防詐騙</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/cookclub/index.html">
                <div class="board-name">cookclub</div>
                <div class="board-nuser"><span class="hl f3">22</span></div>
                <div class="board-class">烹飪</div>
                <div class="board-title">&#9678;歡迎分享手做美食</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/DMM_GAMES/index.html">
                <div class="board-name">DMM_GAMES</div>
                <div class="board-nuser"><span class="hl f3">22</span></div>
                <div class="board-class">遊戲</div>
                <div class="board-title">&#9678;[DMMG] 日幣匯率低 課金好時機</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/MRT/index.html">
                <div class="board-name">MRT</div>
                <div class="board-nuser"><span class="hl f3">22</span></div>
                <div class="board-class">交通</div>
                <div class="board-title">&#9678;[捷運] 北捷多元支付還要等明年</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/bicycle/index.html">
                <div class="board-name">bicycle</div>
                <div class="board-nuser"><span class="hl f3">21</span></div>
                <div class="board-class">車車</div>
                <div class="board-title">&#9678;[單車] 騎車注意交通安全</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/SuperBike/index.html">
                <div class="board-name">SuperBike</div>
                <div class="board-nuser"><span class="hl f3">21</span></div>
                <div class="board-class">車車</div>
                <div class="board-title">&#9678;PTT重機板 ~ 發買賣文請詳閱置底</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/HelpBuy/index.html">
                <div class="board-name">HelpBuy</div>
                <div class="board-nuser"><span class="hl f3">20</span></div>
                <div class="board-class">買賣</div>
                <div class="board-title">&#9678;[代買] 貼文請先看板規及置底公告</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/LCD/index.html">
                <div class="board-name">LCD</div>
                <div class="board-nuser"><span class="hl f3">20</span></div>
                <div class="board-class">硬體</div>
                <div class="board-title">&#9678;[螢幕] LCD板請先看置頂文</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/TypeMoon/index.html">
                <div class="board-name">TypeMoon</div>
                <div class="board-nuser"><span class="hl f3">19</span></div>
                <div class="board-class">型月</div>
                <div class="board-title">&#9678;[F/GO] 浪劍+死神+芙莉蓮聯動(X</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/facelift/index.html">
                <div class="board-name">facelift</div>
                <div class="board-nuser"><span class="hl f3">19</span></div>
                <div class="board-class">美容</div>
                <div class="board-title">&#9678;[醫學美容版]　嚴禁揪團買賣廣告</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/EatToDie/index.html">
                <div class="board-name">EatToDie</div>
                <div class="board-nuser"><span class="hl f3">18</span></div>
                <div class="board-class">美食</div>
                <div class="board-title">&#9678;[飽板] 能吃就是福</div>
            </a>
        </div>
        
        <div class="b-ent">
            <a class="board" href="/bbs/Palmar_Drama/index.html">
                <div class="board-name">Palmar_Drama</div>
                <div class="board-nuser"><span class="hl f3">18</span></div>
                <div class="board-class">布袋</div>
                <div class="board-title">&#9678;4/18 道劫龍戰</div>
            </a>
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

func TestParseHotboardsPage(t *testing.T) {
	p, err := ParseHotboardsPage([]byte(sampleHotboardsHTML))
	if err != nil {
		t.Fatalf("ParseHotboardsPage error: %v", err)
	}
	if len(p.Boards) < 10 { // sanity: large sample should have many boards
		t.Fatalf("expected at least 10 boards, got %d", len(p.Boards))
	}

	// Find Gossiping
	var gossip *HotboardEntry
	for i := range p.Boards {
		if p.Boards[i].BrdName == "Gossiping" {
			gossip = &p.Boards[i]
			break
		}
	}
	if gossip == nil {
		t.Fatalf("Gossiping not found in parsed boards")
	}
	if !gossip.IsBoard || gossip.URL != "/bbs/Gossiping/index.html" || gossip.Class != "綜合" || gossip.Nuser != 6606 {
		t.Fatalf("unexpected Gossiping: %+v", *gossip)
	}
	if gossip.Title == "" || strings.HasPrefix(gossip.Title, "◎") {
		t.Fatalf("unexpected Gossiping title (should not start with ◎): %q", gossip.Title)
	}

	// Spot check another entry exists, e.g., C_Chat
	var cchatFound bool
	for _, b := range p.Boards {
		if b.BrdName == "C_Chat" {
			cchatFound = true
			if !b.IsBoard || b.URL != "/bbs/C_Chat/index.html" || b.Class != "閒談" || b.Nuser <= 0 {
				t.Fatalf("unexpected C_Chat: %+v", b)
			}
			break
		}
	}
	if !cchatFound {
		t.Fatalf("C_Chat not found in parsed boards")
	}
}

func TestParseHotboardsPage_SmallWithGroup(t *testing.T) {
	const sampleSmall = `
    <div class="b-list-container action-bar-margin bbs-screen">
      <div class="b-ent">
        <a class="board" href="/bbs/Gossiping/index.html">
          <div class="board-name">Gossiping</div>
          <div class="board-nuser">8714</div>
          <div class="board-class">綜合</div>
          <div class="board-title">&#9678;八卦</div>
        </a>
      </div>
      <div class="b-ent">
        <a class="board" href="/cls/802">
          <div class="board-name">H_Group</div>
          <div class="board-nuser"></div>
          <div class="board-class">一一</div>
          <div class="board-title">&#931;戰略高手</div>
        </a>
      </div>
    </div>`

	p, err := ParseHotboardsPage([]byte(sampleSmall))
	if err != nil {
		t.Fatalf("ParseHotboardsPage error: %v", err)
	}
	if len(p.Boards) != 2 {
		t.Fatalf("expected 2 boards, got %d", len(p.Boards))
	}
	b0 := p.Boards[0]
	if b0.BrdName != "Gossiping" || b0.Nuser != 8714 || !b0.IsBoard || b0.Class != "綜合" || b0.Title == "" || b0.URL != "/bbs/Gossiping/index.html" {
		t.Fatalf("unexpected b0: %+v", b0)
	}
	if strings.HasPrefix(b0.Title, "◎") {
		t.Fatalf("title should have leading symbol stripped: %q", b0.Title)
	}
	b1 := p.Boards[1]
	if b1.BrdName != "H_Group" || b1.IsBoard || b1.URL != "/cls/802" || b1.Class != "一一" {
		t.Fatalf("unexpected b1: %+v", b1)
	}
	if strings.HasPrefix(b1.Title, "Σ") {
		t.Fatalf("group title should have leading symbol stripped: %q", b1.Title)
	}
}
