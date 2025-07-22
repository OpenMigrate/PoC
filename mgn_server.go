package main
import (
	"fmt"
	"io"
	"os"
	"net"
)

const (
	BlockSize = 1024 * 1024
)

//Serer
type Server struct{
	ln net.Listener
	quitch chan struct {}
}

//Start the server
func (s *Server) StartServer() error{
	ln,err := net.Listen("tcp",":4000")
	if err!=nil{
		return err
	}
	defer ln.Close()
	s.ln=ln

	fmt.Println("Server Started...")
	
	go s.AcceptLoop()
	
	<-s.quitch
	return nil
}

//Accept connections to server
func (s *Server) AcceptLoop(){
	for {
		conn,err:=s.ln.Accept()
		if err!=nil{
			fmt.Println("Error Accepting connection")
			continue
		}
		fmt.Println("Connection Succesfull to Source %s",conn.RemoteAddr)
		
		go s.StartReplication(conn)

	}
}

//start replication
func (s *Server) StartReplication(conn net.Conn){
	//Close the connection when exiting
	defer conn.Close()

	//Open device to write
	dest, err := os.OpenFile("/dev/nvme3n1", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err!=nil{
		fmt.Println(err)
		return
	}

	fmt.Printf("Starting Replication")

	//Receive data
	replicator := &Replicator{
		Reader: conn,
	}

	buffer := make([]byte,512)

	copied,err := io.CopyBuffer(dest,replicator,buffer)
	if err!=nil{
		fmt.Println(err)
	}

	fmt.Println("Replication Completed")
	fmt.Printf("Total copied: %d bytes (%.2f GB)\n", copied, float64(copied)/(1024*1024*1024))
}

type Replicator struct{
	Reader io.Reader
}

func (rp *Replicator) Read(p []byte) (int, error){
	n,err := rp.Reader.Read(p)
	return n,err
}

func main(){
	server := Server{
		quitch: make(chan struct{}),
	}
	err := server.StartServer()
	if err!=nil{
		fmt.Println(err)
	}
}