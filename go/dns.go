package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/miekg/dns"
)

var domain = "u.isucon.dev"
var addr = ":53"
var ErrNotFound = fmt.Errorf("not found")

var initialSubdomains = []string{
	"ns1",
	"pipe",
	"test001",
	"www",
	"www1",
	"www2",
	"www3",
	"www4",
	"www5",
	"mail",
	"acr-nema",
	"afpovertcp",
	"afs3-bos",
	"afs3-callback",
	"afs3-fileserve",
	"afs3-kaserver",
	"afs3-prserver",
	"afs3-rmtsys",
	"afs3-update",
	"afs3-vlserver",
	"afs3-volser",
	"amanda",
	"amandaidx",
	"amidxtape",
	"amqp",
	"amqps",
	"asf-rmcp",
	"asp",
	"auth",
	"babel",
	"bacula-dir",
	"bacula-fd",
	"bacula-sd",
	"bbs",
	"bgp",
	"bgpd",
	"biff",
	"binkp",
	"bootpc",
	"bootps",
	"canna",
	"cfengine",
	"chargen",
	"cisco-sccp",
	"clc-build-daemon",
	"clearcase",
	"cmip-agent",
	"cmip-man",
	"codaauth2",
	"codasrv",
	"codasrv-se",
	"csync2",
	"cvspserver",
	"daap",
	"datametrics",
	"daytime",
	"db-lsp",
	"dcap",
	"dhcpv6-client",
	"dhcpv6-server",
	"dicom",
	"dict",
	"dircproxy",
	"discard",
	"distcc",
	"domain",
	"domain-s",
	"echo",
	"epmap",
	"epmd",
	"exec",
	"f5-globalsite",
	"f5-iquery",
	"fax",
	"fido",
	"finger",
	"font-service ",
	"freeciv",
	"fsp",
	"ftp",
	"ftp-data",
	"ftps",
	"ftps-data",
	"gdomap",
	"gds-db",
	"git",
	"gnunet",
	"gnutella-rtr ",
	"gnutella-svc ",
	"gopher",
	"gpsd",
	"gris",
	"groupwise",
	"gsidcap",
	"gsiftp",
	"gsigatekeeper",
	"hkp",
	"http",
	"http-alt",
	"https",
	"hylafax",
	"iax",
	"icpv2",
	"imap2",
	"imaps",
	"ingreslock",
	"ipp",
	"iprop",
	"ipsec-nat-t",
	"ipx",
	"ircd",
	"ircs-u",
	"isakmp",
	"iscsi-target ",
	"isisd",
	"isns",
	"iso-tsap",
	"kamandakerberos",
	"kerberos-adm",
	"kerberos-master",
	"kerberos4",
	"kermit",
	"klogin",
	"kpasswd",
	"krb-prop",
	"kshell",
	"l2f",
	"ldap",
	"ldaps",
	"ldp",
	"login",
	"lotusnote",
	"mailq",
	"mdns",
	"microsoft-ds",
	"moira-db",
	"moira-update",
	"moira-ureg",
	"mon",
	"ms-sql-m",
	"ms-sql-s",
	"ms-wbt-server",
	"mtn",
	"munin",
	"mysql",
	"mysql-proxy",
	"nbd",
	"nbp",
	"netbios-dgm",
	"netbios-ns",
	"netbios-ssn",
	"netstat",
	"nfs",
	"nntp",
	"nntps",
	"nqs",
	"nrpe",
	"nsca",
	"ntalk",
	"ntp",
	"ntske",
	"nut",
	"omniorb",
	"openvpn",
	"ospf6d",
	"ospfapi",
	"ospfd",
	"passwd-server",
	"pawserv",
	"pop3",
	"pop3s",
	"poppassd",
	"postgresql",
	"predict",
	"printer",
	"proofd",
	"ptp-event",
	"ptp-general",
	"puppet",
	"qmqp",
	"qmtp",
	"qotd",
	"radius",
	"radius-acct",
	"radmin-port",
	"redis",
	"remctl",
	"ripd",
	"ripngd",
	"rmiregistry",
	"rmtcfg",
	"rootd",
	"route",
	"rpc2portmap",
	"rplay",
	"rsync",
	"rtcm-sc104",
	"rtmp",
	"rtsp",
	"sa-msg-port",
	"saft",
	"sane-port",
	"sge-execd",
	"sge-qmaster",
	"sgi-cad",
	"sgi-cmsd",
	"sgi-crsd",
	"sgi-gcd",
	"shell",
	"sieve",
	"silc",
	"sip",
	"sip-tls",
	"skkserv",
	"smtp",
	"smux",
	"snmp",
	"snmp-trap",
	"snpp",
	"socks",
	"spamd",
	"ssh",
	"submission",
	"submissions",
	"sunrpc",
	"supfiledbg",
	"supfilesrv",
	"suucp",
	"svn",
	"svrloc",
	"syslog",
	"syslog-tls",
	"sysrqd",
	"systat",
	"tacacs",
	"talk",
	"tcpmux",
	"telnet",
	"telnets",
	"tfido",
	"tftp",
	"time",
	"tinc",
	"tproxy",
	"uucp",
	"venus",
	"venus-se",
	"webmin",
	"who",
	"whois",
	"wnn6",
	"x11",
	"x11-1",
	"x11-2",
	"x11-3",
	"x11-4",
	"x11-5",
	"x11-6",
	"x11-7",
	"xdmcp",
	"xinetd",
	"xmms2",
	"xmpp-client",
	"xmpp-server",
	"xtel",
	"xtelw",
	"z3950",
	"zabbix-agent",
	"zabbix-trapper",
	"zebra",
	"zebrasrv",
	"zephyr-clt",
	"zephyr-hm",
	"zephyr-srv",
	"zip",
	"zope",
	"zope-ftp",
	"zserv",
	"ayamazaki0",
	"yoshidamiki0",
	"hidekimurakami0",
	"akemikobayashi0",
	"eishikawa0",
	"tomoya100",
	"kobayashiminoru0",
	"yamadakaori0",
	"tanakahanako0",
	"taro660",
	"hashimotokenichi0",
	"yuta330",
	"kenichinakamura0",
	"taichikato0",
	"wfujita0",
	"saitotakuma0",
	"maayafujiwara0",
	"akira680",
	"vmaeda0",
	"jnakamura0",
	"suzukitsubasa0",
	"yoshidatomoya0",
	"qendo0",
	"haruka030",
	"saitotakuma1",
	"bsuzuki0",
	"shohei720",
	"naoko980",
	"suzukiryohei0",
	"kobayashisayuri0",
	"ykobayashi0",
	"asuzuki0",
	"sotaro880",
	"nyamaguchi0",
	"momokohashimoto0",
	"suzukinaoki0",
	"wgoto0",
	"tomoya540",
	"fujitayoichi0",
	"kaorikato0",
	"chiyo810",
	"yito0",
	"tomoyakato0",
	"ryosukeabe0",
	"smatsumoto0",
	"vwatanabe0",
	"sayuri650",
	"takahashinaoto0",
	"hiroshisuzuki0",
	"xyamamoto0",
	"osaito0",
	"esato0",
	"oshimizu0",
	"yamadamituru0",
	"wmori0",
	"saitorei0",
	"kimuramiki0",
	"sasakiyosuke0",
	"kumiko410",
	"qkobayashi0",
	"akondo0",
	"ywatanabe0",
	"otanaoki0",
	"naotosasaki0",
	"momokosuzuki0",
	"maayayoshida0",
	"hashimotoasuka0",
	"maayanakagawa0",
	"ltanaka0",
	"jyamada0",
	"yasuhiro650",
	"hiroshitanaka0",
	"yamashitaakemi0",
	"hidekiishikawa0",
	"kanatakahashi0",
	"akemiito0",
	"manabu800",
	"shota790",
	"naokoaoki0",
	"gyamamoto0",
	"xyamaguchi0",
	"katokumiko0",
	"xyamada0",
	"qyamashita0",
	"pwatanabe0",
	"ryoheiishikawa0",
	"nakamuraharuka0",
	"pfukuda0",
	"saitoyosuke0",
	"morijun0",
	"hiroshiokamoto0",
	"kazuyahayashi0",
	"anakajima0",
	"maaya110",
	"kazuya250",
	"kondoakira0",
	"atsushi920",
	"hanakokondo0",
	"naoki220",
	"hiroshi180",
	"jun820",
	"onakamura0",
	"takuma040",
	"kobayashiyasuhiro0",
	"reiwatanabe0",
	"hideki630",
	"yutayoshida0",
	"yutamori0",
	"ryosukeyamamoto0",
	"manabuota0",
	"nsakamoto0",
	"wtanaka0",
	"inakamura0",
	"ymori0",
	"nfujii0",
	"osamu950",
	"yuinakagawa0",
	"uhayashi0",
	"momoko070",
	"myamada0",
	"zyoshida0",
	"watanabesatomi0",
	"watanabemanabu0",
	"momokosuzuki1",
	"jokada0",
	"hanako640",
	"rmatsumoto0",
	"maayasasaki0",
	"ryohei450",
	"takuma570",
	"yokotanaka0",
	"mai270",
	"rikasakamoto0",
	"taro150",
	"suzukimomoko0",
	"watanabeyoichi0",
	"yuki620",
	"tsubasainoue0",
	"yokotakahashi0",
	"yumikohayashi0",
	"yamaguchiyumiko0",
	"kkato0",
	"minoruyoshida0",
	"tmurakami0",
	"hasegawanaoto0",
	"kumikowatanabe0",
	"lota0",
	"yoichi440",
	"sayurikondo0",
	"xogawa0",
	"naoki250",
	"eokada0",
	"satomiyamamoto0",
	"asukamaeda0",
	"momokowatanabe0",
	"nanamisuzuki0",
	"gotomanabu0",
	"maiyamazaki0",
	"junito0",
	"aokikazuya0",
	"shohei040",
	"momokotanaka0",
	"atsushimatsumoto0",
	"taichifukuda0",
	"haruka630",
	"tanakashohei0",
	"qnishimura0",
	"mkondo0",
	"kimurayui0",
	"tomoyanakajima0",
	"naoko310",
	"yukiokada0",
	"tsubasasuzuki0",
	"yutatakahashi0",
	"kyosuke140",
	"tomoya190",
	"ryohei610",
	"fukudahiroshi0",
	"tyamaguchi0",
	"rsasaki0",
	"satoyoichi0",
	"rikawatanabe0",
	"vokada0",
	"mituru070",
	"nakajimayui0",
	"nanamiota0",
	"smatsuda0",
	"suzukijun0",
	"saitosayuri0",
	"haruka730",
	"yamamotoakemi0",
	"satoyosuke0",
	"maayasato0",
	"nanami830",
	"kanaokamoto0",
	"usuzuki0",
	"manabusasaki0",
	"kenichi870",
	"tomoyasato0",
	"mituruhayashi0",
	"junnishimura0",
	"fujiwaramituru0",
	"yoshidamai0",
	"yosuke350",
	"xwatanabe0",
	"naokimatsumoto0",
	"kenichi170",
	"naoko340",
	"nakajimasotaro0",
	"eito0",
	"kenichi470",
	"qabe0",
	"junkondo0",
	"yosukewatanabe0",
	"ryosuke040",
	"itosotaro0",
	"kimuramikako0",
	"reiogawa0",
	"myamada1",
	"yamazakisayuri0",
	"gotoakemi0",
	"shohei660",
	"phasegawa0",
	"kana990",
	"yosuke710",
	"rika420",
	"ikedajun0",
	"takuma250",
	"yoichihayashi0",
	"yoichi441",
	"zmiura0",
	"ftakahashi0",
	"oinoue0",
	"osamu920",
	"takahashikumiko0",
	"tsubasamaeda0",
	"ryosuke220",
	"yoichimurakami0",
	"abeyumiko0",
	"tishii0",
	"vkato0",
	"mituru710",
	"asukamurakami0",
	"yuiwatanabe0",
	"jokamoto0",
	"akirafujiwara0",
	"lgoto0",
	"akiratakahashi0",
	"lhashimoto0",
	"ogawarika0",
	"manabu620",
	"osamuyamaguchi0",
	"fujiwararei0",
	"yosukewatanabe1",
	"asuka500",
	"otamaaya0",
	"okamotonaoko0",
	"tokamoto0",
	"saitomituru0",
	"suzukitakuma0",
	"suzukimaaya0",
	"nanamitanaka0",
	"suzukikana0",
	"sakamotoshohei0",
	"kanayamamoto0",
	"xtakahashi0",
	"wsuzuki0",
	"taichi990",
	"saitoyuki0",
	"rei050",
	"nmaeda0",
	"csuzuki0",
	"takahashiatsushi0",
	"yukifukuda0",
	"yumiko680",
	"yamadayoko0",
	"abeosamu0",
	"taro210",
	"katomai0",
	"hasegawakumiko0",
	"vhashimoto0",
	"tanakamai0",
	"hmori0",
	"naokoyamaguchi0",
	"yamaguchiyuta0",
	"iyamamoto0",
	"yutanakajima0",
	"kyosukesasaki0",
	"satosatomi0",
	"kkato1",
	"kumiko980",
	"jishii0",
	"yutafujiwara0",
	"fukudamomoko0",
	"kumikokobayashi0",
	"naoki350",
	"sayuri230",
	"morijun1",
	"mituru840",
	"eyoshida0",
	"taichi460",
	"esasaki0",
	"yumikonakamura0",
	"akemi540",
	"kazuya770",
	"suzukiharuka0",
	"harukagoto0",
	"vkato1",
	"bmatsumoto0",
	"yamazakikaori0",
	"yamashitasatomi0",
	"fujiimikako0",
	"ekato0",
	"asukaito0",
	"aokimiki0",
	"katotaichi0",
	"harukasasaki0",
	"manabuyamamoto0",
	"etanaka0",
	"nakagawamituru0",
	"hiroshikobayashi0",
	"satomi900",
	"akira180",
	"manabu710",
	"minoru160",
	"fujiwarataro0",
	"tomoya240",
	"taichinakajima0",
	"atsushitakahashi0",
	"yoichi150",
	"fito0",
	"satomitakahashi0",
	"akemi230",
	"akato0",
	"asuka250",
	"atsushiinoue0",
	"sasakirei0",
	"yamadamiki0",
	"taro890",
	"yamamotohiroshi0",
	"vtakahashi0",
	"yoichiyoshida0",
	"nishimurashohei0",
	"zsasaki0",
	"shimizumaaya0",
	"inoueosamu0",
	"hsato0",
	"kimuraharuka0",
	"takumayoshida0",
	"ikedahideki0",
	"kobayashitsubasa0",
	"jmurakami0",
	"kaori320",
	"suzukimiki0",
	"gtanaka0",
	"shota220",
	"fyamada0",
	"nishimuraminoru0",
	"ttanaka0",
	"satomi280",
	"maaya390",
	"snakamura0",
	"okamototsubasa0",
	"etakahashi0",
	"taro500",
	"yasuhirotakahashi0",
	"uhasegawa0",
	"miki470",
	"ryosukeito0",
	"kanasaito0",
	"dyamashita0",
	"taichi180",
	"iyamazaki0",
	"yui950",
	"kumikomaeda0",
	"yuiwatanabe1",
	"suzukishota0",
	"vsuzuki0",
	"tanakahiroshi0",
	"yosukekimura0",
	"saitoryosuke0",
	"ymori1",
	"xnakajima0",
	"zyamazaki0",
	"mikakoyamamoto0",
	"tanakahideki0",
	"ekato1",
	"momokomori0",
	"tsubasa260",
	"mikako290",
	"yoko160",
	"yuki480",
	"hasegawakumiko1",
	"itokana0",
	"fkobayashi0",
	"tomoya630",
	"yuki730",
	"vyamaguchi0",
	"tshimizu0",
	"pota0",
	"hiroshi960",
	"katoyui0",
	"mkondo1",
	"rika500",
	"einoue0",
	"yamashitayumiko0",
	"kazuya680",
	"asukakobayashi0",
	"miki320",
	"yumikosuzuki0",
	"naotoito0",
	"esuzuki0",
	"hiroshiaoki0",
	"vishikawa0",
	"akiraito0",
	"mikako020",
	"miurasatomi0",
	"etakahashi1",
	"tarosato0",
	"matsumotomanabu0",
	"naoko560",
	"kaori570",
	"minoru210",
	"jsuzuki0",
	"tanakashohei1",
	"tyamada0",
	"kenichi160",
	"shotayamada0",
	"hiroshi070",
	"hidekisakamoto0",
	"naokiikeda0",
	"nakagawamaaya0",
	"hiroshitanaka1",
	"yamashitamaaya0",
	"momokoabe0",
	"vito0",
	"sakamotomikako0",
	"asaito0",
	"ryosukewatanabe0",
	"nanami140",
	"tanakahanako1",
	"rei390",
	"rito0",
	"akemiinoue0",
	"jendo0",
	"cwatanabe0",
	"rikakobayashi0",
	"kaoriokamoto0",
	"momoko350",
	"sayuri130",
	"kaoritanaka0",
	"mikakokobayashi0",
	"naoto500",
	"suzukinaoto0",
	"myamamoto0",
	"yuta090",
	"asato0",
	"chiyo580",
	"lmatsumoto0",
	"kaori190",
	"mai860",
	"naokokimura0",
	"tomoya400",
	"dota0",
	"osamu590",
	"takahashitomoya0",
	"shimizuakemi0",
	"endokenichi0",
	"yoshidaryohei0",
	"ryosukemurakami0",
	"reitanaka0",
	"akemi910",
	"enakamura0",
	"ohayashi0",
	"takumasuzuki0",
	"suzukijun1",
	"bsaito0",
	"maifujiwara0",
	"ikedayoko0",
	"kimurakenichi0",
	"kyosuke210",
	"shota870",
	"fyamaguchi0",
	"miki730",
	"sayurihashimoto0",
	"naoki460",
	"satoyumiko0",
	"minoruhashimoto0",
	"nanamimurakami0",
	"mai600",
	"taro340",
	"yosuke950",
	"rokada0",
	"nanami710",
	"zishikawa0",
	"yutawatanabe0",
	"suzukiosamu0",
	"mituru130",
	"taichi000",
	"chiyoyamamoto0",
	"yoshidakyosuke0",
	"vyamazaki0",
	"wwatanabe0",
	"ykimura0",
	"tanakahideki1",
	"minoruito0",
	"shotamatsumoto0",
	"watanabehiroshi0",
	"taichi560",
	"nanami090",
	"akira800",
	"suzukiyasuhiro0",
	"momokofukuda0",
	"kana350",
	"tanakanaoto0",
	"akemi800",
	"tanakakenichi0",
	"oyoshida0",
	"yutamatsumoto0",
	"hiroshisuzuki1",
	"yamaguchikazuya0",
	"yoko800",
	"lnakamura0",
	"esato1",
	"atsushitakahashi1",
	"eshimizu0",
	"reitanaka1",
	"yuki500",
	"tanakashohei2",
	"akiranishimura0",
	"naokiyamamoto0",
	"shimizuatsushi0",
	"winoue0",
	"hiroshifujita0",
	"mikako380",
	"gotoshohei0",
	"yumikoyamamoto0",
	"gmatsuda0",
	"kumikokobayashi1",
	"tsubasatakahashi0",
	"yyamamoto0",
	"nakamurajun0",
	"yasuhirokato0",
	"vsuzuki1",
	"kobayashihideki0",
	"murakamiosamu0",
	"mikako340",
	"takumaokada0",
	"osamusato0",
	"aogawa0",
	"ikedayosuke0",
	"sakamotoyoichi0",
	"yamadakaori1",
	"ykobayashi1",
	"fukudananami0",
	"yasuhirohashimoto0",
	"akemi670",
	"rhasegawa0",
	"naoto740",
	"satoasuka0",
	"lyamashita0",
	"kumiko030",
	"qikeda0",
	"ptanaka0",
	"fmatsuda0",
	"akiramaeda0",
	"yoichiyamamoto0",
	"gsuzuki0",
	"naokiikeda1",
	"abekenichi0",
	"maiyamamoto0",
	"shimizukumiko0",
	"zsato0",
	"manabu050",
	"shoheinishimura0",
	"gyoshida0",
	"chiyo790",
	"hanakookada0",
	"xota0",
	"manabusakamoto0",
	"jsato0",
	"kenichi840",
	"kobayashishota0",
	"ainoue0",
	"takuma380",
	"dikeda0",
	"dtanaka0",
	"miturutakahashi0",
	"fmurakami0",
	"ishiiakira0",
	"maayaokamoto0",
	"asaito1",
	"atsushi470",
	"ekobayashi0",
	"ukobayashi0",
	"asaito2",
	"hayashirei0",
	"suzukisatomi0",
	"yumiko100",
	"okadahanako0",
	"mtanaka0",
	"hidekiyamamoto0",
	"yamamotomituru0",
	"kimuratomoya0",
	"mikako780",
	"pgoto0",
	"harukawatanabe0",
	"zkimura0",
	"taichikimura0",
	"moriminoru0",
	"thayashi0",
	"akira740",
	"suzukiharuka1",
	"yuta030",
	"fyamamoto0",
	"watanaberyosuke0",
	"maayamatsumoto0",
	"yoichiendo0",
	"satokaori0",
	"gmatsumoto0",
	"rika250",
	"jsasaki0",
	"yutasaito0",
	"kazuya810",
	"manabukobayashi0",
	"hashimotomikako0",
	"fujiwarahideki0",
	"yoichi730",
	"takahashiryosuke0",
	"tyamaguchi1",
	"otaminoru0",
	"momoko990",
	"kazuyahasegawa0",
	"mai770",
	"aishii0",
	"pokada0",
	"cyamaguchi0",
	"akemi500",
	"jtanaka0",
	"tnakagawa0",
	"kondoyoko0",
	"yoichiogawa0",
	"hayashikenichi0",
	"nanami990",
	"hiroshiendo0",
	"akirasuzuki0",
	"maaya140",
	"itomomoko0",
	"satorika0",
	"yoko110",
	"wyamazaki0",
	"yoko111",
	"osato0",
	"aokitaro0",
	"gnakamura0",
	"kobayashiyoko0",
	"suzukinanami0",
	"yosukeishikawa0",
	"manabu450",
	"naoki630",
	"kimuranaoki0",
	"kumiko770",
	"moriyuta0",
	"tsubasa300",
	"hidekinakamura0",
	"rei080",
	"shoheiyoshida0",
	"omatsumoto0",
	"haruka260",
	"takuma790",
	"satoryosuke0",
	"maedamaaya0",
	"mikakosasaki0",
	"vmiura0",
	"rika570",
	"myamaguchi0",
	"minoru110",
	"ishikawakana0",
	"yoko300",
	"tsuzuki0",
	"kenichiwatanabe0",
	"naokosakamoto0",
	"fujitataichi0",
	"naokigoto0",
	"yuta210",
	"sgoto0",
	"momokoishii0",
	"yoshidahideki0",
	"takumainoue0",
	"suzukiyuta0",
	"manabu370",
	"asuka120",
	"manabu621",
	"gogawa0",
	"yamashitaosamu0",
	"jnakamura1",
	"atakahashi0",
	"osamukobayashi0",
	"satomisato0",
	"satorika1",
	"tanakahiroshi1",
	"satoyuta0",
	"mai710",
	"nakamurakaori0",
	"naoko440",
	"jokamoto1",
	"kobayashirei0",
	"kyosukewatanabe0",
	"yamaguchijun0",
	"yoshidahiroshi0",
	"miurahanako0",
	"tsubasa790",
	"hiroshitakahashi0",
	"yukishimizu0",
	"inouenaoki0",
	"yuinakamura0",
	"hashimotoakemi0",
	"jwatanabe0",
	"yosuke890",
	"mikisato0",
	"naotoinoue0",
	"vinoue0",
	"takuma170",
	"yamashitatsubasa0",
	"mituru180",
	"ttanaka1",
	"ogawataichi0",
	"rsakamoto0",
	"satoakemi0",
	"lito0",
	"suzukimaaya1",
	"nakamurarei0",
	"momokofujiwara0",
	"yoko310",
	"fujiwaramanabu0",
	"takahashiosamu0",
	"nakamurakazuya0",
	"nnakamura0",
	"akirasuzuki1",
	"watanabenanami0",
	"rikaito0",
	"juntanaka0",
	"tanakamiki0",
	"ekato2",
	"xnishimura0",
	"ihayashi0",
	"atsushi921",
	"mikakoishikawa0",
	"junmatsuda0",
	"abejun0",
	"katoshota0",
	"shotaishikawa0",
	"yukikimura0",
	"xishii0",
	"rei240",
	"yutakobayashi0",
	"kazuya380",
	"nanamikobayashi0",
	"zyamada0",
	"inouehanako0",
	"yutaota0",
	"satomikato0",
	"yoshidayasuhiro0",
	"akemiikeda0",
	"yoichikobayashi0",
	"msuzuki0",
	"ssato0",
	"qnakamura0",
	"ryoheikobayashi0",
	"naoki880",
	"sakamotosotaro0",
	"haruka350",
	"satonaoki0",
	"jun430",
	"sasakimai0",
	"naoko880",
	"rikakobayashi1",
	"nanamikondo0",
	"yoshidamai1",
	"nakajimarei0",
	"morimanabu0",
	"tanakasayuri0",
	"momoko040",
	"kaori580",
	"haruka100",
	"akemi600",
	"miki170",
	"okamotokyosuke0",
	"kondonaoto0",
	"tsubasa610",
	"ryoheiyamada0",
	"naoki090",
	"yamazakimikako0",
	"minorusato0",
	"chiyofujii0",
	"takuma300",
	"taroyamamoto0",
	"suzukichiyo0",
	"wishikawa0",
	"reinakamura0",
	"ayamazaki1",
	"cogawa0",
	"hasegawamituru0",
	"cyamaguchi1",
	"hasegawamomoko0",
	"suzukishohei0",
	"qyamazaki0",
	"kyosuke370",
	"mgoto0",
	"yasuhiroogawa0",
	"maaya250",
	"kimurayumiko0",
	"yoko460",
	"atsushiyamaguchi0",
	"yamadaasuka0",
	"wyamaguchi0",
	"katoyuta0",
	"yoichi660",
	"lkato0",
	"tomoyatakahashi0",
	"rika370",
	"inouenaoto0",
	"kondoshota0",
	"taro820",
	"maisasaki0",
	"gyoshida1",
	"manabu770",
	"matsudayoichi0",
	"nakajimamiki0",
	"watanabeyoko0",
	"vwatanabe1",
	"yamaguchinaoto0",
	"saitojun0",
	"naotomori0",
	"lito1",
	"ryoheiikeda0",
	"yuisasaki0",
	"yishikawa0",
	"yamazakimituru0",
	"yumiko050",
	"okadasayuri0",
	"skobayashi0",
	"rika740",
	"yoko600",
	"wsato0",
	"taichi181",
	"hiroshi800",
	"yoko420",
	"naokokondo0",
	"sakamotorika0",
	"hiroshiyoshida0",
	"manabusato0",
	"rei340",
	"jsato1",
	"satomikako0",
	"murakamihiroshi0",
	"asukaishii0",
	"kana550",
	"fyamaguchi1",
	"miki190",
	"dkimura0",
	"ynakamura0",
	"endomomoko0",
	"itoyumiko0",
	"yamashitayosuke0",
	"yokokondo0",
	"tsubasa470",
	"matsudanaoko0",
	"miuraatsushi0",
	"miurataichi0",
	"manabusato1",
	"uhasegawa1",
	"nanamimaeda0",
	"katotsubasa0",
	"eyamazaki0",
	"wyamazaki1",
	"junyamada0",
	"akemigoto0",
	"uokada0",
	"nishii0",
	"takumatanaka0",
	"itomai0",
	"hidekimaeda0",
	"osuzuki0",
	"harukasuzuki0",
	"watanabeyuta0",
	"kobayashikaori0",
	"usuzuki1",
	"gotokazuya0",
	"kondotaichi0",
	"nakamurasatomi0",
	"ltakahashi0",
	"jun240",
	"sayuri860",
	"jmurakami1",
	"ftakahashi1",
	"nanami660",
	"otaasuka0",
	"oito0",
	"satomi720",
	"ryosuke680",
	"minoru370",
	"matsumotokumiko0",
	"sishii0",
	"maedanaoko0",
	"chiyo370",
	"sasakitakuma0",
	"reiikeda0",
	"chiyo540",
	"yamamotosayuri0",
	"jota0",
	"bsasaki0",
	"bmori0",
	"itojun0",
	"atsushikondo0",
	"jun320",
	"osamusasaki0",
	"suzukimaaya2",
	"shimizuhideki0",
	"endotomoya0",
	"yamaguchimanabu0",
	"rhashimoto0",
	"rei590",
	"yosuke490",
	"tanakaryohei0",
	"yukinakamura0",
	"akira580",
	"minoru400",
	"yamaguchimomoko0",
	"hanako410",
	"hyamamoto0",
	"maaya910",
	"sayurikobayashi0",
	"matsumotonaoko0",
	"rei250",
	"momoko980",
	"asuka370",
	"asuka100",
	"yukiito0",
	"mai500",
	"manabutanaka0",
	"taichinakamura0",
	"csato0",
	"mikiyoshida0",
	"fujiiatsushi0",
	"hayashinaoko0",
	"kaorikobayashi0",
	"lfujii0",
	"wyamaguchi1",
	"satomifujiwara0",
	"gyamashita0",
	"ryohei850",
	"dnakamura0",
	"rika040",
	"satoyuta1",
	"ynishimura0",
	"tsubasa240",
	"yukigoto0",
	"satomi200",
	"tsubasashimizu0",
	"suzukiyumiko0",
	"nakamurataichi0",
	"yumiko070",
	"shotafujiwara0",
	"maedataro0",
	"maifujita0",
	"fujiwarayuki0",
	"shohei240",
	"aokirika0",
	"hanako580",
	"rgoto0",
	"yuta100",
	"yosukeishikawa1",
	"dokada0",
	"kimurahanako0",
	"wnakagawa0",
	"yukiyoshida0",
	"naokitanaka0",
	"ysaito0",
	"miturusato0",
	"naoki870",
	"chiyo310",
	"naokitanaka1",
	"chiyonakamura0",
	"msaito0",
	"jsuzuki1",
	"suzukihanako0",
	"dokada1",
	"rei850",
	"kenichihashimoto0",
	"fnakamura0",
	"hasegawatakuma0",
	"lyamamoto0",
	"takuma600",
	"manabutakahashi0",
	"ryosukesakamoto0",
	"itoharuka0",
	"satomi130",
	"tomoya450",
}

var userNames = map[string]bool{}
var userNameLock sync.RWMutex

// panic時にプロセスを終了させずにログ出力するハンドラー
func handlePanic(w dns.ResponseWriter, r *dns.Msg) {
	if rcv := recover(); rcv != nil {
		log.Println("[ERR] panic", rcv, w, r)
	}
}

// クエリで指定されたサブドメインをレコードとして応答（エコー機能）するクエリハンドラー
func echoHandler(w dns.ResponseWriter, r *dns.Msg) {
	defer handlePanic(w, r)

	// 関数終了時にDNS応答
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false
	defer w.WriteMsg(m)

	for _, q := range r.Question {
		// サブドメイン部分のみ取り出して、レスポンスのリソースデータとして扱う
		dataLen := len(q.Name) - len(domain) - 2
		subDomain := ""
		if dataLen < -1 {
			continue
		}
		if dataLen >= 0 {
			subDomain = q.Name[:dataLen]
		}

		// この時点で応答できるかは未定だが、応答しようとしている内容をログ出力
		log.Printf("[INFO] query: name=%s class=%s type=%s\n",
			subDomain, dns.ClassToString[q.Qclass], dns.TypeToString[q.Qtype])

		// レスポンスの共通ヘッダー
		// クエリの内容をそのまま使用し、TTLは0固定でキャッシュさせない
		rr_header := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 300}
		rr_header_nx := dns.RR_Header{Name: q.Name, Rrtype: q.Qtype, Class: q.Qclass, Ttl: 0}

		// NewRRやNewZoneParserを使うとほとんどあらゆるレコードに対応可能になるが、
		// 同時にCNAMEによるSubdomain Takeoverなどリスクを負うことになるので必要なものだけ追加する
		switch q.Qclass {
		case dns.ClassINET:
			switch q.Qtype {
			case dns.TypeA:
				ip, err := getIp(subDomain)
				if err != nil {
					m.Answer = append(m.Answer, &dns.A{
						Hdr: rr_header_nx,
					})
					m.Rcode = dns.RcodeNameError
					err2 := w.WriteMsg(m)
					if !errors.Is(err, ErrNotFound) {
						log.Printf("[ERR] DNS Resolution %+v\n", err)
					}
					if err2 != nil {
						log.Printf("[ERR] %s\n", err2.Error())
					}
					continue
				}
				m.Answer = append(m.Answer, &dns.A{
					Hdr: rr_header,
					A:   net.ParseIP(ip),
				})
			case dns.TypeNS:
				if subDomain == "" {
					m.Answer = append(m.Answer, &dns.NS{
						Hdr: rr_header,
						Ns:  "ns1." + domain + ".",
					})
				} else {
					m.Answer = append(m.Answer, &dns.NS{
						Hdr: rr_header_nx,
					})
					m.Rcode = dns.RcodeNameError
				}
			case dns.TypeSOA:
				if subDomain == "" {
					m.Answer = append(m.Answer, &dns.SOA{
						Hdr:     rr_header,
						Ns:      "ns1." + domain + ".",
						Mbox:    "hostmaster." + domain + ".",
						Serial:  0,
						Refresh: 10800,
						Retry:   3600,
						Expire:  604800,
						Minttl:  3600,
					})
				} else {
					m.Answer = append(m.Answer, &dns.SOA{
						Hdr: rr_header_nx,
					})
					m.Rcode = dns.RcodeNameError
				}
			}
		}
	}
	err := w.WriteMsg(m)
	if err != nil {
		log.Printf("[ERR] %s", err.Error())
	}
}

// 指定ネットワークでDNSサーバー処理を実行
func serveDNS(server *dns.Server) {
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func getIp(subDomain string) (string, error) {
	if subDomain == "" {
		return powerDNSSubdomainAddress, nil
	}
	userNameLock.RLock()
	defer userNameLock.RUnlock()
	if _, ok := userNames[subDomain]; ok {
		return powerDNSSubdomainAddress, nil
	}

	//var i int
	//err := dbConn.Get(&i, "SELECT 1 FROM users WHERE name= ?", subDomain)
	//if err != nil {
	//	if errors.Is(err, sql.ErrNoRows) {
	//		return "", ErrNotFound
	//	}
	//	return "", fmt.Errorf("failed to select user: %w", err)
	//}

	//
	//return powerDNSSubdomainAddress, nil

	return "", ErrNotFound
}

func initializeDnsCache(dbOnly bool) error {
	// initialize は同時に呼ばれないので thread-safe ではなくて良い
	userNameLock.Lock()
	defer userNameLock.Unlock()
	userNames = map[string]bool{}
	if !dbOnly {
		for _, subdomain := range initialSubdomains {
			userNames[subdomain] = true
		}
	}

	var names []string
	err := dbConn.Select(&names, "SELECT name FROM users")
	if err != nil {
		return fmt.Errorf("failed to select users: %w", err)
	}

	for _, name := range names {
		userNames[name] = true
	}

	return nil
}
