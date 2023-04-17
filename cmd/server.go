package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"vector-go-sdk-oskr-extensions/pkg/oskrpb"
)

const (
	WIFI_MAX = 100
)

type server struct {
	oskrpb.UnimplementedOSKRServiceServer
}

func (s *server) GetWifiSignalStrength(ctx context.Context, req *oskrpb.WifiSignalStrengthRequest) (*oskrpb.WifiSignalStrengthResponse, error) {
	signalStrength := GetWifiSignalStrengthInt()
	return &oskrpb.WifiSignalStrengthResponse{SignalStrength: int32(signalStrength)}, nil
}

func GetWifiSignalStrengthInt() int {
	cmd := exec.Command("sh", "-c", "iwconfig wlan0 | grep -i \"Signal level\"")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return 0
	}

	signalRegexp := regexp.MustCompile(`Signal level=(-?\d+) dBm`)
	matches := signalRegexp.FindStringSubmatch(string(output))
	if len(matches) < 2 {
		fmt.Println("Error parsing signal level")
		return 0
	}

	signal, err := strconv.Atoi(matches[1])
	if err != nil {
		fmt.Println("Error converting signal level to integer:", err)
		return 0
	}

	// Convert dBm to a percentage value between 0 and WIFI_MAX
	percentage := int((float64(signal+100) / 70) * float64(WIFI_MAX))
	if percentage < 0 {
		percentage = 0
	} else if percentage > WIFI_MAX {
		percentage = WIFI_MAX
	}

	return percentage
}

func main() {
	fmt.Println("OSKR starting listening on port 50051")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		fmt.Println("failed to listen: %v", err)
	}

	fmt.Println("New server")
	grpcServer := grpc.NewServer()
	oskrpb.RegisterOSKRServiceServer(grpcServer, &server{})
	reflection.Register(grpcServer)
	fmt.Println("New server registered")

	fmt.Println("Serving...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
		fmt.Println("failed to serve: %v", err)
	}
}
