package service

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager      *JWTManager
	accessibleRoles map[string][]string
}

//NewAuthInterceptor returns a new auth interceptor
func NewAuthInterceptor(jwtManager *JWTManager, accessibleRoles map[string][]string) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, accessibleRoles}
}

//Unary returns a server interceptor function to authenticate an authorize unary rpc
func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)
		//Check the token and return the user privileges and info
		claims, err := interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}
		if info.FullMethod == "/pb.AuthService/Login" || info.FullMethod == "/pb.AuthService/AccountExistance" {
			return handler(ctx, req)
		}
		res, err := interceptor.redirectUnitary(claims, ctx, info.FullMethod, req)
		if err != nil {
			return nil, err
		}
		return res, err
	}

}
func (interceptor *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Println("--> unary interceptor: ", info.FullMethod)
		_, err := interceptor.authorize(stream.Context(), info.FullMethod)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}

}

//Checks for the JWT and returns if its valid and the user privileges
func (interceptor *AuthInterceptor) authorize(ctx context.Context, method string) (*UserClaims, error) {
	accessibleRoles, ok := interceptor.accessibleRoles[method]
	if !ok {
		//means everyone can access to that service
		return nil, nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "matadata is not provided")
	}
	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}
	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return &UserClaims{}, status.Errorf(codes.Unauthenticated, "access token is  invalid: %v", err)
	}
	//Checking if their user is able to enter to that service
	for _, role := range accessibleRoles {
		if role == claims.Role {
			return claims, nil
		}
	}
	return nil, status.Error(codes.PermissionDenied, "insuficient credentials")
}

//TODO: Change the conn to be created and pass one to the authInterceptor object from packege connections
//This is the method that is going to check the header to redirect to specific microservice
func (interceptor *AuthInterceptor) redirectUnitary(claims *UserClaims, ctx context.Context, method string, req interface{}, opts ...grpc.CallOption) (interface{}, error) {
	//TODO: Change the direction of the localhost to use envoy proxy
	//TODO: Change the grpc connection to use HTTPS
	conn, err := grpc.Dial(":7777", grpc.WithInsecure())
	if err != nil {
		//TODO: Check if the errors should be centralized as constants on one package
		return nil, logError(status.Error(codes.Internal, "Could no stablish the connection with the other service"))
	}
	defer conn.Close()
	header := metadata.New(map[string]string{"role": claims.Role})
	//TODO: CHANGE THE context for calling this service
	updatedCtx := metadata.NewIncomingContext(ctx, header)
	res := new(interface{})
	conn.Invoke(updatedCtx, method, req, res, opts...)
	if err != nil {
		return nil, logError(status.Error(codes.Internal, "connection lost with service to service"))
	}
	return res, nil
}

//
func (interceptor *AuthInterceptor) redirectStream(claim UserClaims) (*interface{}, error) {
	return nil, logError(status.Error(codes.Unimplemented, "This methods has not been implement"))
}
