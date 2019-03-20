package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type output struct {
	writer io.Writer
	lock   sync.Mutex
}

func main() {
	var wg sync.WaitGroup
	scanner := bufio.NewScanner(os.Stdin)
	pongMsg := []byte("PONG :tmi.twitch.tv\r\n")
	output := output{
		os.Stdout,
		sync.Mutex{},
	}

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
			log.Println("Error on sending join command", "channels:", channels, "err:", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			buffer := make([]byte, 0xFFFF)
			for {
				n, err := socket.Read(buffer)
				if err != nil {
					log.Println("Error on reading from connection", "channels handled by this socket:", channels, "buffer:", string(buffer), "n:", n, "err:", err)
					break
				}

				removedCount := deleteLastLine(buffer, n)
				buffer[n-removedCount] = '\n'
				output.lock.Lock()
				output.writer.Write(buffer[0 : n-removedCount+1])
				output.lock.Unlock()

				time.Sleep(1000 * time.Millisecond)

				_, err = socket.Write(pongMsg)
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
			removedCount++
			if buffer[i-1] == '\r' {
				removedCount++
			}
			break
		}

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
