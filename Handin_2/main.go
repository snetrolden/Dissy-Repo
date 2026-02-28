package main

import (
	"time" // used for time out in main
)

func main() {

	p1 := newPeer(3000)
	p2 := newPeer(3001)
	p3 := newPeer(3002)

	//start each Peer on it's own thread so they don't block each other
	go p1.Start()
	go p2.Start()
	go p3.Start()
	time.Sleep(1 * time.Second)

	// p1 and p2 P2P established
	p1.Connect("127.0.0.1", 3001) //localhost IP
	time.Sleep(1 * time.Second)

	//p2 and p3 P2P established
	p2.Connect("127.0.0.1", 3002)
	time.Sleep(1 * time.Second)

	testMsg1 := &Message{
		Type:  "Ping",
		MsgID: "1", //keep ID simple in test
		From:  "127.0.0.1:3000",
	}

	testMsg2 := &Message{
		Type:  "Ping",
		MsgID: "2", //keep ID simple in test
		From:  "127.0.0.1:3002",
	}

	// P1 -> P2
	p1.Send("127.0.0.1:3001", testMsg1)
	// P3 -> P1
	p3.Send("127.0.0.1:3001", testMsg2)

	//sleep to let the network do it's thing before terminating
	time.Sleep(3 * time.Second)

}
