package chatbot

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"golang.org/x/net/context"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

// [END import_libraries]

func main2() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -project-id <PROJECT ID> -session-id <SESSION ID> -language-code <LANGUAGE CODE> <OPERATION> <INPUTS>\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "<PROJECT ID> must be your Google Cloud Platform project id\n")
		fmt.Fprintf(os.Stderr, "<SESSION ID> must be a Dialogflow session ID\n")
		fmt.Fprintf(os.Stderr, "<LANGUAGE CODE> must be a language code from https://dialogflow.com/docs/reference/language; defaults to en\n")
		fmt.Fprintf(os.Stderr, "<OPERATION> must be one of text, audio, stream\n")
		fmt.Fprintf(os.Stderr, "<INPUTS> can be a series of text inputs if <OPERATION> is text, or a path to an audio file if <OPERATION> is audio or stream\n")
	}

	var projectID, sessionID, languageCode string
	flag.StringVar(&projectID, "project-id", "", "Google Cloud Platform project ID")
	flag.StringVar(&sessionID, "session-id", "", "Dialogflow session ID")
	flag.StringVar(&languageCode, "language-code", "en", "Dialogflow language code from https://dialogflow.com/docs/reference/language; defaults to en")

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	operation := args[0]
	inputs := args[1:]

	switch operation {
	case "text":
		fmt.Printf("Responses:\n")
		for _, query := range inputs {
			fmt.Printf("\nInput: %s\n", query)
			response, err := DetectIntentText(projectID, sessionID, query, languageCode)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Output: %s\n", response)
		}

	default:
		flag.Usage()
		os.Exit(1)
	}
}

// [START dialogflow_detect_intent_text]
func DetectIntentText(projectID, sessionID, text, languageCode string) (string, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return "", err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return "", errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return "", err
	}

	queryResult := response.GetQueryResult()
	fulfillmentText := queryResult.GetFulfillmentText()
	return fulfillmentText, nil
}
