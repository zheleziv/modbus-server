package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/gopcua/opcua/uacp"
)

func main() {
	var (
		endpoint = flag.String("endpoint", "opc.tcp://localhost:49321/ZHELEZNOV", "OPC UA Endpoint URL")
	)
	flag.Parse()

	ctx := context.Background()

	log.Printf("Listening on %s", *endpoint)
	l, err := uacp.Listen(*endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("1321")
	for {
		c, err := l.Accept(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("conn %d: connection from %s", c.ID(), c.RemoteAddr())
		buf, _ := c.Receive()

		// var buf []byte
		// if i, err := c.Read(buf); err != nil {
		// 	fmt.Println(err.Error() + " " + fmt.Sprintf("%d", i))
		// } else {
		fmt.Println(buf)
		// fmt.Println(i)
		// }
	}
	// log.Printf("conn %d: connection from %s", c.ID(), c.RemoteAddr())

	// 	listener, err := uacp.Listen(*endpoint, uint32(*bufsize))
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Printf("Started listening on %s.", listener.Endpoint())
	// uasc.
	// 	cfg := uasc.NewServerConfig(
	// 		"http://opcfoundation.org/UA/SecurityPolicy#None",
	// 		nil, nil, 1111, services.SecModeNone, 2222, 3600000,
	// 	)
	// for {
	// 	func() {
	// 		ctx := context.Background()
	// 		ctx, cancel := context.WithCancel(ctx)
	// 		defer cancel()

	// 		conn, err := l.Accept(ctx)
	// 		if err != nil {
	// 			log.Print(err)
	// 			return
	// 		}
	// 		defer func() {
	// 			conn.Close()
	// 			log.Println("Successfully closed connection")
	// 		}()
	// 		log.Printf("Successfully established connection with %v", conn.RemoteAddr())

	// 		secChan, err := uasc.ListenAndAcceptSecureChannel(ctx, conn, cfg)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		defer func() {
	// 			secChan.Close()
	// 			log.Printf("Successfully closed secure channel with %v", conn.RemoteAddr())
	// 		}()
	// 		log.Printf("Successfully opened secure channel with %v", conn.RemoteAddr())

	// 		sessCfg := uasc.NewServerSessionConfig(secChan)
	// 		session, err := uasc.ListenAndAcceptSession(ctx, secChan, sessCfg)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		defer func() {
	// 			session.Close()
	// 			log.Printf("Successfully closed session with %v", conn.RemoteAddr())
	// 		}()
	// 		log.Printf("Successfully activated session with %v", conn.RemoteAddr())

	// 		buf := make([]byte, 1024)
	// 		for {
	// 			n, err := session.ReadService(buf)
	// 			if err != nil {
	// 				log.Printf("Couldn't read UASC: %s", err)
	// 				continue
	// 			}
	// 			log.Printf("Successfully received message: %x\n%s", buf[:n], utils.Wireshark(0, buf[:n]))

	// 			srv, err := services.Decode(buf[:n])
	// 			if err != nil {
	// 				log.Printf("Couldn't decode received bytes as Service: %s", err)
	// 				continue
	// 			}
	// 			log.Printf("Successfully decoded as Service: %v", srv)
	// 		}
	// 	}()
	// }
}
