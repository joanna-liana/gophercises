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

type wrongAnswerInfo struct {
	Question      string
	PlayerAnswer  string
	CorrectAnswer string
}

type programDefaults struct {
	DefaultFile      string
	DefaultTimeLimit int
}

func getDefaults() programDefaults {
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

	return programDefaults{
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

func runQuiz(problems [][]string, timer *time.Timer, userAnswerChannel chan string) {
	var wrongAnswersInfo []wrongAnswerInfo

	questionsCount := len(problems)
	correctCount := 0

ProblemLoop:
	for _, problem := range problems {
		question := problem[0]
		answer := strings.TrimSpace(problem[1])

		fmt.Print(question, "> ")

		go getUserAnswer(userAnswerChannel)

		select {
		case <-timer.C:
			fmt.Println("\nThe time is up!")
			break ProblemLoop

		case input := <-userAnswerChannel:
			parsedUserAnswer := strings.TrimSpace(input)

			if parsedUserAnswer == answer {
				fmt.Print("Well done!\n\n")
				correctCount++
			} else {
				fmt.Print("Wrong answer\n\n")
				answerInfo := wrongAnswerInfo{Question: question, PlayerAnswer: parsedUserAnswer, CorrectAnswer: answer}
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

func runTimedQuiz(problems [][]string, timeLimit int, userAnswerChannel chan string) {
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	runQuiz(problems, timer, userAnswerChannel)
}

func main() {
	defaults := getDefaults()
	filepath, timeLimit := defaults.DefaultFile, defaults.DefaultTimeLimit
	problems := prepareProblems(filepath)

	fmt.Println("Time limit: " + strconv.Itoa(timeLimit) + " seconds")
	fmt.Println("Press enter to start the quiz")
	fmt.Scanf("\n")

	userAnswerChannel := make(chan string)

	runTimedQuiz(problems, timeLimit, userAnswerChannel)

	fmt.Println("Wanna try again?")

	replayAnswer := <-userAnswerChannel

	parsedReplayAnswer := strings.TrimSpace(strings.ToLower(replayAnswer))

	if parsedReplayAnswer == "y" {
		runTimedQuiz(problems, timeLimit, userAnswerChannel)
	} else {
		fmt.Println("See ya!")
	}
}
