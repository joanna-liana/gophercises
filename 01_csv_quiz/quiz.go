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
)

func main() {
    defaultFile := "questions.csv"

    csvFileFlag := flag.String("csv", defaultFile, "A CSV file in the format: question, answer")
    timeTimitFlag := flag.Int("limit", 30, "The time limit for answering a single question")
    flag.Parse()

    fmt.Println(csvFileFlag)
    fmt.Println(timeTimitFlag)

    csvBuffer, err := os.Open(defaultFile)
    if err != nil {
        log.Fatal(err)
    }

    cvsReader := csv.NewReader(bufio.NewReader(csvBuffer))
    for {
        line, error := cvsReader.Read()

        if error == io.EOF {
            break
        } else if error != nil {
            log.Fatal(error)
        }

        question := line[0]
        answer := line[1]

        buf := bufio.NewReader(os.Stdin)
        fmt.Print(question, "> ")
        userAnswer, err := buf.ReadString('\n')

        if err != nil {
            fmt.Println(err)
        } else {
            if strings.TrimSpace(userAnswer) == strings.TrimSpace(answer) {
                fmt.Println("Well done!")
            } else {
                fmt.Println("Wrong answer")
            }
        }
    }
}
