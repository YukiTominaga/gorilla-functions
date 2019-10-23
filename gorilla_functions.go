package gorilla_functions

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

var myMux = NewRouter()

var projectID = os.Getenv("GCP_PROJECT")
var redisHost = os.Getenv("REDIS_HOST")
var nodePort  = os.Getenv("NODE_PORT")

var redisClient *redis.Client

type Value struct {
	Value string `json:"value"`
}

func init() {
	var err error

	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, nodePort),
		Password: "",
		DB: 0,
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
}

func RedisGetHandler(w http.ResponseWriter, r *http.Request) {
	myMux.ServeHTTP(w, r)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	key := params["key"]
	val, err := redisClient.Get(key).Result()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w)
		return
	}

	value := Value{
		Value: val,
	}

	response, _ := json.Marshal(&value)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(response))
}

func indexHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Printf("%s", "indexHandler")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w)
}

// Route ルーター構造体
type Route struct {
	// Name 名前
	Name string
	// Method HTTPメソッド
	Method string
	// Pattern URLパターン
	Pattern string
	// HandlerFunc 実行される関数
	HandlerFunc http.HandlerFunc
}

// Routes ルーター
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		indexHandler,
	},

	Route{
		"GetHandler",
		"GET",
		"/key/{key}",
		getHandler,
	},
}

// NewRouter コンストラクタ　routesで定義したroute構造体配列を用いて、muxのルーターを作成
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		handler = http.Handler(route.HandlerFunc)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
