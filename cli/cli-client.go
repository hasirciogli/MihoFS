package cli

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
)

func SendCommand(command string) {
	// Custom Unix socket transport
	tr := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}

	client := &http.Client{Transport: tr}

	// URL'yi oluştur ve isteği gönder
	url := fmt.Sprintf("http://unix/%s", command)
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Yanıtı oku
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	fmt.Println("Response:", string(body))
}
