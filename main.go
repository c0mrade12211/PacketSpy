package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
	"net"
	"github.com/common-nighthawk/go-figure"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"golang.org/x/term"
	"runtime"
)

func OpenWireshark() {
	fileName := "captured.pcap"
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", "wireshark", fileName)
	} else {
		cmd = exec.Command("wireshark", fileName)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
func ExitAndOpenWireshark() {
	fmt.Println("Exiting program...")
	OpenWireshark()
	os.Exit(0)
}

func hello_function() {
	figure := figure.NewFigure("Give Me Packets", "", true)
	figure.Print()
	fmt.Println(" -> Created by c0mrade <-")
}

func CapturePackets() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[----------------------------------------------------------------]")
	fmt.Println("[+] Available devices")
	for i, device := range devices {
		fmt.Printf("%d. %s\n", i+1, device.Name)
	}

	var choice int
	fmt.Print("Enter device number: ")
	_, err = fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(devices) {
		log.Fatal("Invalid interface choice")
	}

	handle, err := pcap.OpenLive(devices[choice-1].Name, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	file, err := os.Create("captured.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	w := pcapgo.NewWriter(file)
	w.WriteFileHeader(65536, layers.LinkTypeEthernet)
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	fmt.Println("Press 'q' to exit and open Wireshark")
	termios, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(int(syscall.Stdin), termios)

	go func() {
		for packet := range packetSource.Packets() {
			w.WritePacket(packet.Metadata().CaptureInfo, packet.Data())
			fmt.Println("[+] Packet captured")
		}
	}()

	var b []byte = make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err != nil {
			log.Fatal(err)
		}
		if strings.ToLower(string(b)) == "q" {
			ExitAndOpenWireshark()
		}
	}
}

func SendPcapToServer() {
	var server_ip string
	var server_port int

	fmt.Print("[+] Enter server IP: \n")
	_, err := fmt.Scanln(&server_ip)
	if err != nil {
		fmt.Println("[-] Invalid server IP")
		log.Fatal(err)
	}
	fmt.Print("[+] Enter server port: \n")
	_, err = fmt.Scanf("%d", &server_port)
	if err != nil {
		fmt.Println("[-] Error: Port")
		log.Fatal(err)
	}
	fmt.Println("[+] Connecting to server...")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server_ip, server_port))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	file, err := os.Open("captured.pcap")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	fileSize := make([]byte, 8)
	fileSize[0] = byte(fileInfo.Size() >> 56)
	fileSize[1] = byte(fileInfo.Size() >> 48)
	fileSize[2] = byte(fileInfo.Size() >> 40)
	fileSize[3] = byte(fileInfo.Size() >> 32)
	fileSize[4] = byte(fileInfo.Size() >> 24)
	fileSize[5] = byte(fileInfo.Size() >> 16)
	fileSize[6] = byte(fileInfo.Size() >> 8)
	fileSize[7] = byte(fileInfo.Size())

	_, err = conn.Write(fileSize)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 4096)
	for {
		bytesRead, err := file.Read(buffer)
		if err != nil {
			break
		}
		_, err = conn.Write(buffer[:bytesRead])
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("[+] File sent successfully!")

}

func main() {
	hello_function()
	time.Sleep(100 * time.Millisecond)

	fmt.Println("[1] Capture Packets and save to file (and open Wireshark)")
	fmt.Println("[2] Capture Packets and send to server")

	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || (choice != 1 && choice != 2) {
		log.Fatal("Invalid choice")
	}

	if choice == 1 {
		CapturePackets()
	} else if choice == 2 {
		SendPcapToServer()
	}
}
