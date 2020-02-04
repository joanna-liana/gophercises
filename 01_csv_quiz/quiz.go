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
    "time"
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
    defaultTimeLimit := 10

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
    flag.Int("limit", defaultTimeLimit, "The time limit for answering a single question (in seconds)")
    flag.Parse()

    flag.Visit(setDefaults)

    return ProgramDefaults {
        DefaultFile: defaultFile,
        DefaultTimeLimit: defaultTimeLimit,
     }
}

func runQuiz(filepath string, channel chan string) {
    csvBuffer, err := os.Open(filepath)
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

    channel <- "Thanks for playing!"
}

func runTimedQuiz(filepath string, timeLimit int) {
    quizRun := make(chan string, 1)
    go runQuiz(filepath, quizRun)

    select {
    case res := <- quizRun:
        fmt.Println(res)
    case <-time.After(time.Duration(timeLimit) * time.Second):
        // TODO: stop the code inside c1
        fmt.Println("\nThe time is up!")
    }
}

func main() {
    defaults := getDefaults()
    filepath, timeLimit := defaults.DefaultFile, defaults.DefaultTimeLimit

    fmt.Println("Time limit:" + strconv.Itoa(timeLimit) + " seconds")

    buf := bufio.NewReader(os.Stdin)
    fmt.Println("Press enter to start the quiz")
    _, err := buf.ReadString('\n')

    if err != nil {
        fmt.Println(err)
        return
    }

    runTimedQuiz(filepath, timeLimit)
    fmt.Println("Wanna try again?")
    answer, err := buf.ReadString('\n')

    if err != nil {
        fmt.Println(err)
        return
    }

    if (answer == "y") {
        runTimedQuiz(filepath, timeLimit)
    }

    fmt.Println("See ya!")
}
