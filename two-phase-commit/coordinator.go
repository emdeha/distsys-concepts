package main

import (
	"net"
	"fmt"
	"io/ioutil"
	"strings"
	"bufio"
)

type Participant struct {
	connectionString string
}

type Answer int

const (
	AnswerNone = Answer(iota)

	AnswerAck
	AnswerNack
)

func main() {
	participants, err := loadParticipants()
	if err != nil {
		panic(err)
	}

	var answers []Answer

	fmt.Printf("initialized participants: %v\n", participants)

	for _, p := range participants {
		answer, err := propose(p, "test")
		if err != nil {
			answer = AnswerNack
		}
		fmt.Printf("sent proposal: %v\n", answer)
		answers = append(answers, answer)
	}

	if len(answers) < len(participants) || hasNack(answers) {
		fmt.Println("will abort proposal")
		for _, p := range participants {
			err := abort(p)
			if err != nil {
				fmt.Printf("abort failed for %v\n", p)
				continue
			}
			fmt.Printf("aborted for %v\n", p)
		}
		return
	}

	fmt.Println("will commit proposal")
	for _, p := range participants {
		err := commit(p)
		if err != nil {
			fmt.Printf("commit failed for %v\n", p)
			continue
		}
		fmt.Printf("committed for %v\n", p)
	}
}

// loadParticipants loads the configuration for the participants. This
// configuration resides in config.txt. Each line in this file represents a
// node with its type and address. The type can either be participant or
// coordinator.
//
// Example:
//     participant 127.0.0.1:1337
//     participant 127.0.0.1:1338
//     coordinator 127.0.0.1:31337
func loadParticipants() ([]Participant, error) {
	contents, err := ioutil.ReadFile("./config.txt")
	if err != nil {
		fmt.Printf("loadParticipants: error reading config: %v\n", err)
		return []Participant{}, err
	}
	fmt.Printf("loadParticipants: config.txt contents: %v\n", string(contents))

	nodes := strings.Split(string(contents), "\n")
	fmt.Printf("loadParticipants: read nodes: %v\n", nodes)
	var participants []Participant
	for i := 0; i < len(nodes) - 1; i++ {
		nodeData := strings.Split(nodes[i], " ")
		fmt.Printf("loadParticipants: nodeData %v\n", nodeData)
		nodeType := nodeData[0]
		nodeAddress := nodeData[1]
		if nodeType == "participant" {
			participants = append(
				participants, Participant{connectionString: nodeAddress})
		}
	}

	return participants, nil
}

func hasNack(answers []Answer) bool {
	for _, a := range answers {
		if a == AnswerNack {
			return true
		}
	}
	return false
}

// Phase 1. The coordinator proposes a value to be agreed upon.
func propose(p Participant, value string) (Answer, error) {
	conn, err := net.Dial("tcp", p.connectionString)
	if err != nil {
		fmt.Println("propose: can't connect to participant")
		return AnswerNone, err
	}
	fmt.Fprintf(conn, "propose %s\n", value)
	answer, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("propose: can't receive answer from participant")
		return AnswerNone, err
	}

	switch answer {
	case "Ack\n":
		return AnswerAck, nil
	case "Nack\n":
		return AnswerNack, nil
	default:
		fmt.Printf("propose: invalid answer %s", answer)
		return AnswerNone, fmt.Errorf("propse: invalid answer %s", answer)
	}
}

// Phase 2. The coordinator either commits the proposal or aborts it.
// The coordinator commits only when all the participants agree.
func commit(p Participant) error {
	conn, err := net.Dial("tcp", p.connectionString)
	if err != nil {
		fmt.Println("commit: can't connect to participant")
		return err
	}
	fmt.Fprint(conn, "commit\n")
	return nil
}

// The coordinator aborts only when some of the participants disagree.
func abort(p Participant) error {
	conn, err := net.Dial("tcp", p.connectionString)
	if err != nil {
		fmt.Println("abort: can't connect to participant")
		return err
	}
	fmt.Fprint(conn, "abort\n")
	return nil
}
