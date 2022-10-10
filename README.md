###GO Rest DB Seed Project

## Steps:
brew install mongodb-community@4.2
brew services start mongodb-community
mongo
> use platform
> load("$GOPATH/src/git.nexgen.neustar.biz/foundation/user-service/scripts/mongodb/seeds/platform/platform.auditlog.seed.js")

### Mongo
1. 	Start mongo server: brew services start mongodb-community
2. Install Studio 3T for mongo client to view documents
3. Create DB : use platform
4. Create a document in the DB (for some reason the DB vanishes after server restart if it doesnt have any documents) db.products.insert( { item: "card", qty: 15 }).  insert documents into the products collection. If the collection does not exist, the insert() method creates the collection.
5. If you want to load js you can use : load("~/path_to_the_folder/mongo-seed.js") see a sample file attached here
6. go build 
7. ./go-seed-db (Do not do go run main.go)
8. MyAPIs collection in Postman has all the APIs (sameeksha.shetty@neustar.biz)

### Rest
http://localhost:8080/v1/trainer
http://localhost:8080/v1/trainer/b6fb8b4a-fbc1-498b-9aa7-b6d8770daf7d

### GQL
Generate GQL models, generated.go, gqlgenyml:
go run github.com/99designs/gqlgen
go mod download - If running gqlgen wipes out the go.sum entry of gqlgen.
go build


GQL playground
GetTrainerByID

{
	TrainerByID(id:"b6fb8b4a-fbc1-498b-9aa7-b6d8770daf7d"){
    name
  }
}

GetTrainers
{
	Trainers{
    name
  }
}


### Run
1. go build
2. ./go-seed-db