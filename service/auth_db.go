package service

import (
	"context"
	"log"
	"time"

	"github.com/Mau-MR/authzService/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthDB struct {
	db             *DB
	userCollection *mongo.Collection
}

//This for creation of new accounts
type accountDB struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Surname  string             `bson:"surname"`
	User     string             `bson:"user"`
	Password string             `bson:"password"`
	Phone    string             `bson:"phone"`
	Bussines string             `bson:"bussines"`
	Role     string             `bson:"role"`
	DB       primitive.ObjectID `bson:"db"`
}

//TODO: Find a way to hide the database name with env and change it later
func NewAuthDB(db *DB) *AuthDB {
	userCollection := db.client.Database("logdb").Collection("user")
	return &AuthDB{db, userCollection}
}

//TODO: change for a trasaction to make a single call to check also the phone
func (auth *AuthDB) checkAccount(user string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//TODO:check if it is truly necesary to decode on this part
	account := accountDB{}
	if err := auth.userCollection.FindOne(ctx, bson.D{{Key: "user", Value: user}}).Decode(&account); err != nil {
		return logError(status.Errorf(codes.Internal, "user not found: :c"))
	}
	return nil
}

//TODO:find a way to use the check account method for this too
func (auth *AuthDB) login(req *pb.LoginRequest) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	account := User{}
	if err := auth.userCollection.FindOne(ctx, bson.D{{Key: "user", Value: req.Username}}).Decode(&account); err != nil {
		return nil, logError(status.Errorf(codes.Internal, "user not found: %v", err))
	}
	return &account, nil
}
func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}

//TODO:create a process to add child node to owners
func (auth *AuthDB) adminAccount() {

}

//TODO:create a process o add child node for employee
//this function creates a bussines owner account
func (auth *AuthDB) createAccount(account *accountDB) (*pb.CreateAccountRes, error) {

	account.ID = primitive.NewObjectID()
	account.DB = primitive.NewObjectID()
	account.Role = "owner"

	insertResult, err := auth.userCollection.InsertOne(context.Background(), account)
	if err != nil {
		return nil, logError(status.Errorf(codes.Internal, "unable to push new account to the database: %v", err))
	}
	return &pb.CreateAccountRes{Success: true, Type: &pb.CreateAccountRes_Bid{Bid: insertResult.InsertedID.(primitive.ObjectID).Hex()}}, nil
}
