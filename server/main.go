package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"watcharis/go-poc-grpc/go-poc-grpc/pb-example/proto"
)

var (
	port         = ":50051"
	StorePersons = map[int32]*proto.PersonResponse{
		1: {
			Name:  "John Doe",
			Id:    1,
			Email: "john@example.com",
		},
	}
	StoreDBPersons = map[int32]*proto.InsertPersonResponse{
		1: {
			Id:      1,
			Name:    "John Doe",
			Email:   "john@example.com",
			PhoneNo: "0998383333",
		},
	}
)

type server struct {
	proto.UnimplementedUserServiceServer
}

// implement GetUser ตามที่ตกลงไว้ใน .proto
func (s *server) GetUser(ctx context.Context, req *proto.PersonRequest) (*proto.PersonResponse, error) {
	log.Printf("request GetUser : %+v\n", req.ProtoReflect().Descriptor())

	result, exists := StorePersons[req.Id]
	if !exists {
		errStr := "userId : %d, not found"
		return nil, status.Errorf(codes.NotFound, errStr, req.Id)
	}

	return result, nil
}

func (s *server) InsertUser(ctx context.Context, req *proto.InsertPersonRequest) (*proto.InsertPersonResponse, error) {
	log.Printf("request InsertUser : %+v\n", req.ProtoReflect().Descriptor())

	id := int32(len(StoreDBPersons) + 1)

	if result, exists := StoreDBPersons[id]; exists {

		var errStr string
		if req.Name == result.Name {
			errStr = "request field name is duplicate : %s"
			return nil, status.Errorf(codes.NotFound, errStr, req.Name)

		} else if req.Email == result.Email {
			errStr = "request field email is duplicate : %s"
			return nil, status.Errorf(codes.NotFound, errStr, req.Email)

		} else if req.PhoneNo == result.PhoneNo {
			errStr = "request field phone_no is duplicate : %s"
			return nil, status.Errorf(codes.NotFound, errStr, req.PhoneNo)
		} else {
			id += 1
		}

	} else {
		person := &proto.InsertPersonResponse{
			Id:      id,
			Name:    req.Name,
			Email:   req.Email,
			PhoneNo: req.PhoneNo,
		}

		StoreDBPersons[id] = person
	}

	resp := &proto.InsertPersonResponse{
		Id:      id,
		Name:    req.Name,
		Email:   req.Email,
		PhoneNo: req.PhoneNo,
	}

	return resp, nil
}

// และสร้างเป็น Server gRPC ขึ้นมา
func main() {
	// สร้าง TCP listener ที่ฟังการเชื่อมต่อบน port ที่กำหนด
	lis, err := net.Listen("tcp", port)
	if err != nil {
		// ถ้าเกิดข้อผิดพลาดขณะสร้าง listener, ให้ log error แล้วออกจากโปรแกรม
		log.Fatalf("failed to listen: %v", err)
	}

	// สร้าง gRPC server instance ใหม่
	grpcServer := grpc.NewServer()

	// ลงทะเบียน User service กับ server
	proto.RegisterUserServiceServer(grpcServer, &server{})

	// Log ว่า server กำลังฟังการเชื่อมต่อบน port ที่กำหนด
	log.Printf("Server is listening on port %v", port)

	// เริ่มฟังการเชื่อมต่อและให้บริการ gRPC request
	if err := grpcServer.Serve(lis); err != nil {
		// ถ้าเกิดข้อผิดพลาดขณะให้บริการ, ให้ log error แล้วออกจากโปรแกรม
		log.Fatalf("failed to serve: %v", err)
	}
}
