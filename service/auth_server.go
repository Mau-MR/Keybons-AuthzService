package service

import (
	"context"

	"github.com/Mau-MR/authzService/pb"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

type AuthServer struct {
	jwtManager *JWTManager
	authDB     *AuthDB
}

func NewAuthServer(db *DB, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{
		jwtManager, NewAuthDB(db),
	}
}
func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := server.authDB.login(req)
	//TODO: finish this part of the authentication code
	if err != nil || !IsCorrectPassword(user, req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}
	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "unable to generete token")
	}
	res := &pb.LoginResponse{AccessToken: token}
	return res, nil

}

//TODO: Change the response just to show errors here at the server
//TODO: test for this
func (server *AuthServer) AccountExistance(ctx context.Context, req *pb.AccountExistanceReq) (*pb.AccountExistanceRes, error) {
	err := server.authDB.checkAccount(req.User)
	if err != nil {
		return &pb.AccountExistanceRes{Exist: false}, err
	}
	return &pb.AccountExistanceRes{
		Exist: true,
	}, nil
}

//CreateAccount inserts bussines account to db
func (server *AuthServer) CreateAccount(ctx context.Context, req *pb.CreateAccountReq) (*pb.CreateAccountRes, error) {
	if err := server.authDB.checkAccount(req.User); err == nil {
		return nil, logError(status.Error(codes.AlreadyExists, "username already exist"))
	}
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, logError(status.Errorf(codes.Internal, "unable to hash password: %v", err))
	}
	newAccount := accountDB{Name: req.Name, Surname: req.Surname, User: req.User, Password: hashedPassword, Phone: req.Phone, Bussines: req.Bussines}
	status, err := server.authDB.createAccount(&newAccount)
	if err != nil {
		return nil, err
	}
	return status, nil
}

//TODO:make createuser for normal clients
