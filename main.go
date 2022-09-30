package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/handler"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type API struct {
	logger *log.Logger
	db     *gorm.DB
}

var (
	mongourl = "localhost:27017"
	dbname   = "platform"
	pgaddr   = "host=ultrawaf.cluster-civmdcc4k4yd.us-east-1.rds.amazonaws.com dbname=ultrawaf sslmode=require user=ultraadmin password=jEJ38dHJanfijpJQjl8321JFQQApju"
)

func main() {
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() //cancel context
	logger := log.New(os.Stdout, "go-seed", log.Ldate|log.Ltime|log.Llongfile)

	a := &API{logger: logger}
	dbcloser, err := a.ConnectDB(pgaddr)
	defer dbcloser()
	r := mux.NewRouter()
	h := r.PathPrefix("/v1").Subrouter()
	g := r.PathPrefix("/v2").Subrouter()

	//h.Use(auth.Auth0Middleware) check in go-rest-seed
	h.HandleFunc("/trainer", a.ListTrainers).Methods("GET")
	h.HandleFunc("/trainer", a.CreateTrainer).Methods("POST")
	h.HandleFunc("/trainer/{id:[a-z0-9]+}", a.DeleteTrainer).Methods("DELETE")
	h.HandleFunc("/trainer/{id:[a-z0-9]+}", a.GetTrainer).Methods("GET")
	h.HandleFunc("/cheese", a.SayCheese).Methods("GET")

	//Playground second parameter must be the uri path to the graphql APIs
	g.HandleFunc("/playground", handler.Playground("GraphQL playground", "/v2/gql"))
	//g.HandleFunc("/gql", handler.GraphQL(NewExecutableSchema(Config{Resolvers: &Resolver{db: db}})))

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
