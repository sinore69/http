package main

import (
	"fmt"
	"net"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func parseHTTPRequest(data string) (*Request, error) {
	lines := strings.Split(data, "\r\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("empty request")
	}
	requestLine := strings.Split(lines[0], " ")
	if len(requestLine) != 3 {
		return nil, fmt.Errorf("malformed request line")
	}
	req := &Request{
		Method:  requestLine[0],
		Path:    requestLine[1],
		Version: requestLine[2],
		Headers: make(map[string]string),
	}
	i := 1
	for ; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			i++ 
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			req.Headers[key] = val
		}
	}
	if i < len(lines) {
		req.Body = strings.Join(lines[i:], "\r\n")
	}
	return req, nil
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Listening on port 8080...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go func(c net.Conn) {
			defer c.Close()

			buffer := make([]byte, 4096)
			n, err := c.Read(buffer)
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			raw := string(buffer[:n])
			req, err := parseHTTPRequest(raw)
			if err != nil {
				fmt.Println("Parse error:", err)
				return
			}
			fmt.Printf("Received request: %+v\n", req)
			resp := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello, World!"
			c.Write([]byte(resp))
		}(conn)
	}
}
