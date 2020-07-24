package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// WrongAnswerInfo contains question and answers info
type WrongAnswerInfo struct {
	Question      string
	PlayerAnswer  string
	CorrectAnswer string
}

// ProgramDefaults represents default program settings
type ProgramDefaults struct {
	DefaultFile      string
	DefaultTimeLimit int
}

func getDefaults() ProgramDefaults {
	defaultFile := "questions.csv"
	defaultTimeLimit := 10

	setDefaults := func(f *flag.Flag) {
		switch f.Name {
		case "csv":
			defaultFile = f.Value.String()

		case "limit":
			intLimit, err := strconv.Atoi(f.Value.String())
			if err != nil {
				log.Fatal(err)
				return
			}
			defaultTimeLimit = intLimit
		}
	}

	flag.String("csv", defaultFile, "A CSV file in the format: question, answer")
	flag.Int("limit", defaultTimeLimit, "The time limit for answering a single question (in seconds)")
	flag.Parse()

	flag.Visit(setDefaults)

	return ProgramDefaults{
		DefaultFile:      defaultFile,
		DefaultTimeLimit: defaultTimeLimit,
	}
}

func prepareProblems(filepath string) [][]string {
	csvBuffer, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
	}

	cvsReader := csv.NewReader(bufio.NewReader(csvBuffer))

	problems, error := cvsReader.ReadAll()

	if error != nil {
		log.Fatal(error)
	}

	return problems
}

func getUserAnswer(answerChan chan string) {
	buf := bufio.NewReader(os.Stdin)
	userAnswer, err := buf.ReadString('\n')

	if err != nil {
		fmt.Println(err)
	}

	answerChan <- userAnswer
}

func runQuiz(problems [][]string, timer *time.Timer) {
	var wrongAnswersInfo []WrongAnswerInfo

	questionsCount := len(problems)
	correctCount := 0

ProblemLoop:
	for _, problem := range problems {
		question := problem[0]
		answer := strings.TrimSpace(problem[1])
		userAnswerChannel := make(chan string)

		fmt.Print(question, "> ")

		go getUserAnswer(userAnswerChannel)

		select {
		case <-timer.C:
			fmt.Println("\nThe time is up!")
			break ProblemLoop

		case input := <-userAnswerChannel:
			close(userAnswerChannel)
			parsedUserAnswer := strings.TrimSpace(input)

			if parsedUserAnswer == answer {
				fmt.Print("Well done!\n\n")
				correctCount++
			} else {
				fmt.Print("Wrong answer\n\n")
				answerInfo := WrongAnswerInfo{Question: question, PlayerAnswer: parsedUserAnswer, CorrectAnswer: answer}
				wrongAnswersInfo = append(wrongAnswersInfo, answerInfo)
			}
		}
	}

	wrongAnswersCount := len(wrongAnswersInfo)

	fmt.Println("Final score:", correctCount, "/", questionsCount)

	if wrongAnswersCount > 0 {
		fmt.Print("\nWrong answers - details\n\n")
		for _, info := range wrongAnswersInfo {
			fmt.Println("Question:", info.Question)
			fmt.Println("Your answer:", info.PlayerAnswer)
			fmt.Print("Correct answer:", info.CorrectAnswer, "\n\n")
		}
	}

	fmt.Println("Thanks for playing!")
}

func runTimedQuiz(problems [][]string, timeLimit int) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)
	runQuiz(problems, timer)
}

func main() {
	defaults := getDefaults()
	filepath, timeLimit := defaults.DefaultFile, defaults.DefaultTimeLimit
	problems := prepareProblems(filepath)

	fmt.Println("Time limit:" + strconv.Itoa(timeLimit) + " seconds")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press enter to start the quiz")
	_, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	runTimedQuiz(problems, timeLimit)
	fmt.Println("Wanna try again?")
	// TODO: this does not work when inputting y/n
	// you have to press enter again and then the program ALWAYS exist ("See ya!")
	answer, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	if strings.TrimSpace(answer) == "y" {
		runTimedQuiz(problems, timeLimit)
	} else {
		fmt.Println("See ya!")
	}
}
