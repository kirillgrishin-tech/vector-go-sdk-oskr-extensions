package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
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

func TriggerWakeWord() {
	var cloudSock ipc.Conn

	log.Println("Creating Test Client Socket to send messages to vic-cloud")
	cloudSock = getSocketWithRetry(ipc.GetSocketPath("cloud_sock"), "cp_test")
	defer cloudSock.Close()
	log.Println("Socket created")

	location_currentzone, _ := time.LoadLocation("Local")
	log.Println("Triggering hotword: " + location_currentzone.String())
	hw := cloud.Hotword{Mode: cloud.StreamType_Normal, Locale: "en-US", Timezone: location_currentzone.String(), NoLogging: true}
	message := cloud.NewMessageWithHotword(&hw)

	log.Println("Creating sender")
	testSender := voice.IPCMsgSender{Conn: cloudSock}
	log.Println("Sending message")
	testSender.Send(message)
	log.Println("DONE")
}

func main() {
	/*
		var err error
		log.Println("Creating Test Server Socket to receive from vic-cloud")
		testSock, err = ipc.NewUnixgramServer(ipc.GetSocketPath("cp_test"))
		if err != nil {
			log.Println("Server create error:", err)
		}
		defer testSock.Close()
	*/

	//var testSend voice.MsgIO
	//testSend, testRecv = voice.NewMemPipe()

	//go testReader(testSock, testSend)

	/*
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
	*/
	log.Println("Now triggering wake word")
	TriggerWakeWord()
}

func getSocketWithRetry(name string, client string) ipc.Conn {
	for {
		sock, err := ipc.NewUnixgramClient(name, client)
		if err != nil {
			log.Println("Couldn't create socket", name, "- retrying:", err)
			time.Sleep(5 * time.Second)
		} else {
			return sock
		}
	}
}

func testReader(serv ipc.Server, send voice.MsgSender) {
	for conn := range serv.NewConns() {
		go func(conn ipc.Conn) {
			for {
				msg := conn.ReadBlock()
				if msg == nil || len(msg) == 0 {
					conn.Close()
					return
				}
				var cmsg cloud.Message
				if err := cmsg.Unpack(bytes.NewBuffer(msg)); err != nil {
					log.Println("Test reader unpack error:", err)
					continue
				}
				send.Send(&cmsg)
			}
		}(conn)
	}
}
