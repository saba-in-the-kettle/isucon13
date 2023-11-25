package main

// ISUCON的な参考: https://github.com/isucon/isucon12-qualify/blob/main/webapp/go/isuports.go#L336
// sqlx的な参考: https://jmoiron.github.io/sqlx/

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/isucon/isucon13/webapp/go/isuutil"
	"github.com/kaz/pprotein/integration/echov4"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	echolog "github.com/labstack/gommon/log"
)

const (
	listenPort                     = 8080
	powerDNSSubdomainAddressEnvKey = "ISUCON13_POWERDNS_SUBDOMAIN_ADDRESS"
)

var (
	powerDNSSubdomainAddress string
	dbConn                   *sqlx.DB
	secret                   = []byte("isucon13_session_cookiestore_defaultsecret")
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if secretKey, ok := os.LookupEnv("ISUCON13_SESSION_SECRETKEY"); ok {
		secret = []byte(secretKey)
	}
}

type InitializeResponse struct {
	Language string `json:"language"`
}

func connectDB(logger echo.Logger) (*sqlx.DB, error) {
	const (
		networkTypeEnvKey = "ISUCON13_MYSQL_DIALCONFIG_NET"
		addrEnvKey        = "ISUCON13_MYSQL_DIALCONFIG_ADDRESS"
		portEnvKey        = "ISUCON13_MYSQL_DIALCONFIG_PORT"
		userEnvKey        = "ISUCON13_MYSQL_DIALCONFIG_USER"
		passwordEnvKey    = "ISUCON13_MYSQL_DIALCONFIG_PASSWORD"
		dbNameEnvKey      = "ISUCON13_MYSQL_DIALCONFIG_DATABASE"
		parseTimeEnvKey   = "ISUCON13_MYSQL_DIALCONFIG_PARSETIME"
	)

	conf := mysql.NewConfig()

	// 環境変数がセットされていなかった場合でも一旦動かせるように、デフォルト値を入れておく
	// この挙動を変更して、エラーを出すようにしてもいいかもしれない
	conf.Net = "tcp"
	conf.Addr = net.JoinHostPort("127.0.0.1", "3306")
	conf.User = "isucon"
	conf.Passwd = "isucon"
	conf.DBName = "isupipe"
	conf.ParseTime = true

	if v, ok := os.LookupEnv(networkTypeEnvKey); ok {
		conf.Net = v
	}
	if addr, ok := os.LookupEnv(addrEnvKey); ok {
		if port, ok2 := os.LookupEnv(portEnvKey); ok2 {
			conf.Addr = net.JoinHostPort(addr, port)
		} else {
			conf.Addr = net.JoinHostPort(addr, "3306")
		}
	}
	if v, ok := os.LookupEnv(userEnvKey); ok {
		conf.User = v
	}
	if v, ok := os.LookupEnv(passwordEnvKey); ok {
		conf.Passwd = v
	}
	if v, ok := os.LookupEnv(dbNameEnvKey); ok {
		conf.DBName = v
	}
	if v, ok := os.LookupEnv(parseTimeEnvKey); ok {
		parseTime, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("failed to parse environment variable '%s' as bool: %+v", parseTimeEnvKey, err)
		}
		conf.ParseTime = parseTime
	}

	db, err := isuutil.NewIsuconDB(conf)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func initializeHandler(c echo.Context) error {
	if out, err := exec.Command("../sql/init.sh").CombinedOutput(); err != nil {
		c.Logger().Errorf("init.sh failed with err=%s", string(out))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}

	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index livestream_tags_tag_id_index\n    on livestream_tags (tag_id);\n\n"); err != nil {
		c.Logger().Errorf("create index failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index livestream_tags_livestream_id_index\n    on livestream_tags (livestream_id);\n\n"); err != nil {
		c.Logger().Errorf("create index failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index icons_user_id_index\n    on icons (user_id);\n\n"); err != nil {
		c.Logger().Errorf("create index failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "alter table icons\n    add image_hash varchar(256) default '' not null;\n\n"); err != nil {
		c.Logger().Errorf("create image hash failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index themes_user_id_index\n    on themes (user_id);\n\n"); err != nil {
		c.Logger().Errorf("create image hash failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index reactions_livestream_id_created_at_index\n    on reactions (livestream_id asc, created_at desc);\n\n"); err != nil {
		c.Logger().Errorf("create image hash failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index ng_words_user_id_livestream_id_created_at_index\n    on ng_words (user_id asc, livestream_id asc, created_at desc);\n\n"); err != nil {
		c.Logger().Errorf("create image hash failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index livestream_viewers_history_user_id_livestream_id_index\n    on livestream_viewers_history (user_id, livestream_id);\n\n"); err != nil {
		c.Logger().Errorf("create image hash failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index livestream_viewers_history_livestream_id_index\n    on livestream_viewers_history (livestream_id);\n\n"); err != nil {
		c.Logger().Errorf("create image hash failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	if err := isuutil.CreateIndexIfNotExists(dbConn, "create index reaction_livestream_id_index\n    on reactions (livestream_id);\n\n"); err != nil {
		c.Logger().Errorf("create reaction_livestream_id_index failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}

	if out, err := exec.Command("../pdns/init_zone.sh").CombinedOutput(); err != nil {
		c.Logger().Warnf("init.sh failed with err=%s", string(out))
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}

	// iconsディレクトリの中身をすべて削除する
	if err := os.RemoveAll("../icons"); err != nil {
		c.Logger().Warnf("failed to remove icons directory with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}
	// iconsディレクトリを作成する
	if err := os.Mkdir("../icons", 0755); err != nil {
		c.Logger().Warnf("failed to create icons directory with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}

	if err := isuutil.KickPproteinCollect(); err != nil {
		c.Logger().Warnf("pprotein collect failed with err=%s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to initialize: "+err.Error())
	}

	c.Request().Header.Add("Content-Type", "application/json;charset=utf-8")
	return c.JSON(http.StatusOK, InitializeResponse{
		Language: "golang",
	})
}

func main() {
	_, err := isuutil.InitializeTracerProvider()
	if err != nil {
		panic(err)
	}

	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(echolog.DEBUG)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time=${time_rfc3339_nano} method=${method}, uri=${uri}, status=${status}, latency=${latency_human}, error=${error}\n",
	}))
	cookieStore := sessions.NewCookieStore(secret)
	cookieStore.Options.Domain = "*.u.isucon.dev"
	e.Use(session.Middleware(cookieStore))
	// e.Use(middleware.Recover())
	e.Use(otelecho.Middleware("webapp"))

	// 初期化
	e.POST("/api/initialize", initializeHandler)

	// top
	e.GET("/api/tag", getTagHandler)
	e.GET("/api/user/:username/theme", getStreamerThemeHandler)

	// livestream
	// reserve livestream
	e.POST("/api/livestream/reservation", reserveLivestreamHandler)
	// list livestream
	e.GET("/api/livestream/search", searchLivestreamsHandler)
	e.GET("/api/livestream", getMyLivestreamsHandler)
	e.GET("/api/user/:username/livestream", getUserLivestreamsHandler)
	// get livestream
	e.GET("/api/livestream/:livestream_id", getLivestreamHandler)
	// get polling livecomment timeline
	e.GET("/api/livestream/:livestream_id/livecomment", getLivecommentsHandler)
	// ライブコメント投稿
	e.POST("/api/livestream/:livestream_id/livecomment", postLivecommentHandler)
	e.POST("/api/livestream/:livestream_id/reaction", postReactionHandler)
	e.GET("/api/livestream/:livestream_id/reaction", getReactionsHandler)

	// (配信者向け)ライブコメントの報告一覧取得API
	e.GET("/api/livestream/:livestream_id/report", getLivecommentReportsHandler)
	e.GET("/api/livestream/:livestream_id/ngwords", getNgwords)
	// ライブコメント報告
	e.POST("/api/livestream/:livestream_id/livecomment/:livecomment_id/report", reportLivecommentHandler)
	// 配信者によるモデレーション (NGワード登録)
	e.POST("/api/livestream/:livestream_id/moderate", moderateHandler)

	// livestream_viewersにINSERTするため必要
	// ユーザ視聴開始 (viewer)
	e.POST("/api/livestream/:livestream_id/enter", enterLivestreamHandler)
	// ユーザ視聴終了 (viewer)
	e.DELETE("/api/livestream/:livestream_id/exit", exitLivestreamHandler)

	// user
	e.POST("/api/register", registerHandler)
	e.POST("/api/login", loginHandler)
	e.GET("/api/user/me", getMeHandler)
	// フロントエンドで、配信予約のコラボレーターを指定する際に必要
	e.GET("/api/user/:username", getUserHandler)
	e.GET("/api/user/:username/statistics", getUserStatisticsHandler)
	e.GET("/api/user/:username/icon", getIconHandler)
	e.POST("/api/icon", postIconHandler)

	// stats
	// ライブ配信統計情報
	e.GET("/api/livestream/:livestream_id/statistics", getLivestreamStatisticsHandler)

	// 課金情報
	e.GET("/api/payment", GetPaymentResult)

	e.HTTPErrorHandler = errorResponseHandler

	e.JSONSerializer = &isuutil.EchoJSONSerializer{}

	echov4.EnableDebugHandler(e)

	// DB接続
	conn, err := connectDB(e.Logger)
	if err != nil {
		e.Logger.Errorf("failed to connect db: %v", err)
		os.Exit(1)
	}
	defer conn.Close()
	dbConn = conn

	subdomainAddr, ok := os.LookupEnv(powerDNSSubdomainAddressEnvKey)
	if !ok {
		e.Logger.Errorf("environ %s must be provided", powerDNSSubdomainAddressEnvKey)
		os.Exit(1)
	}
	powerDNSSubdomainAddress = subdomainAddr

	// HTTPサーバ起動
	listenAddr := net.JoinHostPort("", strconv.Itoa(listenPort))
	if err := e.Start(listenAddr); err != nil {
		e.Logger.Errorf("failed to start HTTP server: %v", err)
		os.Exit(1)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func errorResponseHandler(err error, c echo.Context) {
	c.Logger().Errorf("error at %s: %+v", c.Path(), err)
	if he, ok := err.(*echo.HTTPError); ok {
		if e := c.JSON(he.Code, &ErrorResponse{Error: err.Error()}); e != nil {
			c.Logger().Errorf("%+v", e)
		}
		return
	}

	if e := c.JSON(http.StatusInternalServerError, &ErrorResponse{Error: err.Error()}); e != nil {
		c.Logger().Errorf("%+v", e)
	}
}
