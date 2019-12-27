package main

import (
    "bufio"
    "os"
    "log"
	"fmt"
    "flag"
    "io"
    "encoding/csv"
    "strings"
    "strconv"
)

// WrongAnswerInfo contains question and answers info
type WrongAnswerInfo struct {
    Question string
    PlayerAnswer string
    CorrectAnswer string
}

// ProgramDefaults represents default program settings
type ProgramDefaults struct {
    DefaultFile string
    DefaultTimeLimit int
}

func getDefaults() ProgramDefaults {
    defaultFile := "questions.csv"
    defaultTimeLimit := 30

    setDefaults := func (f *flag.Flag) {
        switch f.Name {
            case "csv" :
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
    flag.Int("limit", defaultTimeLimit, "The time limit for answering a single question")
    flag.Parse()

    flag.Visit(setDefaults)

    return ProgramDefaults {
        DefaultFile: defaultFile,
        DefaultTimeLimit: defaultTimeLimit,
     }
}

func main() {
    defaults := getDefaults()
    defaultFile, defaultTimeLimit := defaults.DefaultFile, defaults.DefaultTimeLimit

    fmt.Println("Time limit:", defaultTimeLimit)

    csvBuffer, err := os.Open(defaultFile)
    if err != nil {
        log.Fatal(err)
    }

    questionsCount := 0
    var wrongAnswersInfo []WrongAnswerInfo

    cvsReader := csv.NewReader(bufio.NewReader(csvBuffer))
    for {
        line, error := cvsReader.Read()

        if error == io.EOF {
            break
        } else if error != nil {
            log.Fatal(error)
        }

        questionsCount++
        question := line[0]
        answer := strings.TrimSpace(line[1])

        buf := bufio.NewReader(os.Stdin)
        fmt.Print(question, "> ")
        userAnswer, err := buf.ReadString('\n')
        parsedUserAnswer := strings.TrimSpace(userAnswer)

        if err != nil {
            fmt.Println(err)
        } else {
            if  parsedUserAnswer == answer {
                fmt.Print("Well done!\n\n")
            } else {
                fmt.Print("Wrong answer\n\n")
                answerInfo := WrongAnswerInfo{ Question: question, PlayerAnswer: parsedUserAnswer, CorrectAnswer: answer }
                wrongAnswersInfo = append(wrongAnswersInfo, answerInfo)
            }
        }
    }

    wrongAnswersCount := len(wrongAnswersInfo)

    fmt.Println("Final score:", questionsCount - wrongAnswersCount, "/", questionsCount)

    if wrongAnswersCount > 0 {
        fmt.Print("\nWrong answers - details\n\n")
        for _, info := range(wrongAnswersInfo) {
            fmt.Println("Question:", info.Question)
            fmt.Println("Your answer:", info.PlayerAnswer)
            fmt.Print("Correct answer:", info.CorrectAnswer, "\n\n")
        }
    }
}
