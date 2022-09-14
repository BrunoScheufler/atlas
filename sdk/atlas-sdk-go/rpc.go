package sdk

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brunoscheufler/atlas/core"
	"github.com/brunoscheufler/atlas/protobuf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type server struct {
	protobuf.UnimplementedAtlasfileServer
	atlasfile *atlas.Atlasfile
}

func (s server) Eval(ctx context.Context, request *protobuf.EvalRequest) (*protobuf.EvalReply, error) {
	marshaled, err := json.Marshal(s.atlasfile)
	if err != nil {
		return nil, fmt.Errorf("could not marshal atlasfile: %w", err)
	}

	return &protobuf.EvalReply{
		Output: string(marshaled),
	}, nil
}

func (s server) Ping(ctx context.Context, request *protobuf.PingRequest) (*protobuf.PingReply, error) {
	return &protobuf.PingReply{}, nil
}

func serve(ctx context.Context, logger *logrus.Entry, atlasfile *atlas.Atlasfile, port int) error {
	grpcServer := grpc.NewServer()
	atlasFileServer := &server{
		atlasfile: atlasfile,
	}

	protobuf.RegisterAtlasfileServer(grpcServer, atlasFileServer)

	logger.WithFields(logrus.Fields{
		"port": port,
	}).Traceln("starting atlasfile server")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	logger.WithFields(logrus.Fields{
		"port": port,
	}).Traceln("starting to serve")

	go func() {
		<-ctx.Done()
		logger.Traceln("shutting down atlasfile server")
		grpcServer.GracefulStop()
		logger.Traceln("shut down atlasfile server")
	}()

	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	logger.Traceln("stopped serving")

	return nil
}
