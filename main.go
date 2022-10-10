package main

import (
	"context"
	"log"
	"net/http"

	"net/http/pprof"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"
	"github.com/sshetty10/go-seed-db/generated"
	"gorm.io/gorm"
)

type API struct {
	logger *log.Logger
	db     *gorm.DB
}

var (
	mongourl = "localhost:27017"
	dbname   = "platform"
	pgaddr   = "host=localhost dbname=mydb user=scott password=tiger"
)

func main() {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() //cancel context
	logger := log.New(os.Stdout, "go-seed", log.Ldate|log.Ltime|log.Llongfile)

	a := &API{logger: logger}
	dbcloser, err := a.ConnectDB(pgaddr)
	if err != nil {
		logger.Fatalf("Connection to DB failed %v", err)
	}
	defer dbcloser()
	r := mux.NewRouter()
	h := r.PathPrefix("/v1").Subrouter()
	g := r.PathPrefix("/v2").Subrouter()

	// Handlers for profiling
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.Handle("/debug/pprof/block", pprof.Handler("block"))

	//h.Use(auth.Auth0Middleware) check in go-rest-seed
	h.HandleFunc("/trainer", a.ListTrainers).Methods("GET")
	h.HandleFunc("/trainer", a.CreateTrainer).Methods("POST")
	h.HandleFunc("/trainer/{id:[a-z0-9]+}", a.DeleteTrainer).Methods("DELETE")
	h.HandleFunc("/trainer/{id:[a-z0-9-]+}", a.GetTrainer).Methods("GET")
	h.HandleFunc("/cheese", a.SayCheese).Methods("GET")

	//Playground second parameter must be the uri path to the graphql APIs
	g.HandleFunc("/playground", handler.Playground("GraphQL playground", "/v2/gql"))
	g.HandleFunc("/gql", handler.GraphQL(generated.NewExecutableSchema(generated.Config{Resolvers: &Resolver{api: a}})))

	log.Println("connect to http://localhost:8080/v2/playground for GraphQL playground")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln(err)
	}
}

//ConnectMongo connects to the mongodb
/*func ConnectMongo(ctx context.Context, url string, dbname string) (*Db, error) {
	// Rest of the code will go here
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", url))
	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &Db{MongoCli: client, MongoDb: client.Database(dbname)}, nil
}*/
