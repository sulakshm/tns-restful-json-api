package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	"errors"
	"os"
	"os/signal"
)

/**** overload listener interface for stopping it *******/
var stoppedError = errors.New("Listener stopped")
type stoppableListener struct {
        *net.TCPListener          //Wrapped listener
        cancel             chan int //Channel used only to indicate listener should shutdown
}

func NewstoppableListener(proto, netaddr string) *stoppableListener {
        originalListener, err := net.Listen(proto, netaddr)
        if err != nil {
                panic(err)
        }

        sl, err := makeNew(originalListener)
        if err != nil {
                panic(err)
        }

        return sl
}

func makeNew(l net.Listener) (*stoppableListener, error) {
        tcpL, ok := l.(*net.TCPListener)
        if !ok {
                return nil, errors.New("Cannot wrap listener")
        }

        retval := &stoppableListener{}
        retval.TCPListener = tcpL
        retval.cancel = make(chan int)

        return retval, nil
}

func (sl *stoppableListener) Accept() (net.Conn, error) {
        for {
                //Wait up to one second for a new connection
                sl.SetDeadline(time.Now().Add(time.Second))

                newConn, err := sl.TCPListener.Accept()

                //Check for the channel being closed
                select {
                case <-sl.cancel:
                        return nil, stoppedError
                default:
                        //If the channel is still open, continue as normal
                }

                if err != nil {
                        netErr, ok := err.(net.Error)

                        //If this is a timeout, then continue to wait for
                        //new connections
                        if ok && netErr.Timeout() && netErr.Temporary() {
                                continue
                        }
                }

                return newConn, err
        }
}

func (sl* stoppableListener) stop() {
	close(sl.cancel)
}


var sl *stoppableListener
var wg		sync.WaitGroup

func InitNetwork() {
        netaddr := ":8080" // localhost:8080
        log.Printf("Setting up network serving at %s\n", netaddr)
        sl = NewstoppableListener("tcp", netaddr)
}

func StartServer() {
	router := NewRouter()
        server := http.Server{Handler: router}
	wg.Add(1)
	go func() {
		defer wg.Done()
		server.Serve(sl)
	}()
}

func StopServer() {
	sl.stop()
	wg.Wait()
}

func main() {
	InitNetwork()
	StartServer()

        log.Printf("Wait for quit signal\n")
        stop := make(chan os.Signal)
        //signal.Notify(stop, syscall.SIGINT)
        signal.Notify(stop, os.Interrupt)
        signal := <- stop
        fmt.Printf("Got signal:%v\n", signal)
        StopServer()

	for k, v := range dbMap {
		fmt.Printf("Commit DB: %s\n", k);
		v.Commit()
	}
        fmt.Printf("Good, Bye!")
}
