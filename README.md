# Give_Me_Packets  
Hello! This program is a simple packet sniffer written in the Go programming language. It allows the user to capture and save network packets passing through a selected network device.    
  
Who might need this program? It is useful for developers, system administrators, and anyone who wants to study and analyze network traffic. The program provides a simple and convenient way to capture packets and save them in the pcap format, which can be further analyzed using the Wireshark program.  
  
Advantages of this program:  
  
1. Ease of use: The program provides a simple command-line interface that allows you to select the desired network device and start packet capture with a single command.  
 
2. Flexibility: Users can choose any available network device to capture packets. The program also provides the option to save captured packets in a pcap file for further analysis.  
  
3. Integration with Wireshark: After completing the packet capture, the program automatically opens the pcap file in Wireshark, allowing users to study and analyze network traffic in more detail.  


# Compile  
1. build -o Give_me_Packets.exe ./main.go



# Example   
1) Capture packets and open in WireShark  
![image](https://github.com/c0mrade12211/Give_Me_Packets/assets/132468035/6abfd1ff-d422-4efd-adaf-f4fed3feefd1)   
![image](https://github.com/c0mrade12211/Give_Me_Packets/assets/132468035/71acc74c-8326-450e-9b4a-ffea2e58b2af)
2) Also you're can to send pcap file in the server (command for save file: nc -l -p <port> > output.pcap)  
![image](https://github.com/c0mrade12211/Give_Me_Packets/assets/132468035/a79bc27a-c73d-4e5b-a85d-369c485cdd63)
![image](https://github.com/c0mrade12211/Give_Me_Packets/assets/132468035/1a1cf675-b136-4061-b9c6-1c9ee3018c58)
