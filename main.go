package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
	"strings"
	"syscall"
	"github.com/common-nighthawk/go-figure"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
	"golang.org/x/term"
)

func OpenWireshark() {
	fileName := "captured.pcap"
	cmd := exec.Command("wireshark", fileName)
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


func hello_function(){
	figure := figure.NewFigure("Give Me Packets", "", true)

	figure.Print()
	fmt.Println(" -> Created by c0mrade <-")
}

func main() {

	hello_function()

	time.Sleep(100 * time.Millisecond)
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