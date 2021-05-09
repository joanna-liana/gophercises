package story

import (
	"bufio"
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Option struct {
	Text string
	Arc string
}

type StoryPart struct {
	Title string;
	Story []string;
	Options []Option;
}

type ParsedStory = map[string]StoryPart

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func readFile(storyPath string) *bufio.Reader {
	f, err := os.Open(storyPath)

    check(err)

	reader := bufio.NewReader(f)

	return reader
}

func GetParsedStory(storyPath string) ParsedStory {
	reader := readFile(storyPath)

	dec := json.NewDecoder(reader)

	var story ParsedStory

	// TODO: compare this with unmarshal
	err := dec.Decode(&story)

	if err != nil {
		log.Fatal(err)
	}

	return story
}

func GetStoryIntro(parsedStory ParsedStory) StoryPart {
	// TODO: how to ensure JSON order?
	intro, exists := parsedStory["intro"]

	if !exists {
		for _, v := range(parsedStory) {
			intro = v

			break;
		}

	}

	return intro
}

func GetArcByName(parsedStory ParsedStory, arcName string) (error, StoryPart) {
	intro, exists := parsedStory[arcName]

	if !exists {
		return errors.New("The given arc does not exist"), StoryPart{}
	}

	return nil, intro
}
