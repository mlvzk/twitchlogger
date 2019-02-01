package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
	"time"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)

	for {
		end, channels := readBulk(scanner, 400)
		log.Println("channels", len(channels), "end", end)
		joinCommand := "JOIN #" + strings.Join(channels, ",#") + "\r\n"

		socket, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
		if err != nil {
			log.Println("Error on connecting to twitch irc:", err)
			continue
		}

		_, err = socket.Write([]byte("PASS 123\r\nNICK justinfan12345\r\n"))
		if err != nil {
			log.Println("Error on authentication:", err)
			continue
		}

		_, err = socket.Write([]byte(joinCommand))
		if err != nil {
			log.Println("Error on sending join command:", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			buffer := make([]byte, 0xFFFF)
			for {
				n, err := socket.Read(buffer)
				if err != nil {
					log.Println("Error on reading from connection:", err)
					break
				}

				removedCount := deleteLastLine(buffer, n)
				buffer[n-removedCount] = '\n'
				log.Println("n:", n)
				os.Stdout.Write(buffer[0:n-removedCount+1])

				for i := 0; i < n; i++ {
					buffer[i] = 0
				}

				time.Sleep(1000 * time.Millisecond)

				_, err = socket.Write([]byte("PONG :tmi.twitch.tv\r\n"))
				if err != nil {
					log.Println("Error on sending PONG:", err)
					break
				}
			}
		}()

		if end {
			break
		}

		time.Sleep(30 * time.Second)
	}

	wg.Wait()
}

func deleteLastLine(buffer []byte, length int) int {
	removedCount := 0

	for i := length - 1; i >= 0; i-- {
		if buffer[i] == '\n' {
			buffer[i] = 0
			removedCount++
			if buffer[i-1] == '\r' {
				buffer[i-1] = 0
				removedCount++
			}
			break
		}

		buffer[i] = 0
		removedCount++
	}

	return removedCount
}

type Scanner interface {
	Scan() bool
	Text() string
}

func readBulk(scanner Scanner, length int) (bool, []string) {
	a := make([]string, length)

	for i := 0; i < length; i++ {
		if !scanner.Scan() {
			return true, a
		}

		a[i] = scanner.Text()
	}

	return false, a
}
