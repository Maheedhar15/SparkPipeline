package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/joho/godotenv"
)

type payload struct {
	Text string `json:"text"`
}

func createChangeReq() string {

	body := `{"result":{"error":"","conflict_details":[{"event_name":"","event_start_time":"","event_end_time":"","applies_to":"","condition":"","blackout_schedule":"","blackout_schedule_type":""}]}}`

	//fmt.Println("Response Body:", string(body))
	response_data := string(body)
	if strings.Contains(response_data, "error") || strings.Contains(response_data, "fatal") {
		message := fmt.Sprintf("Change creation failure occured in prod pipeline. Here are the details: \n \n \n Build ID: %s\n Build Pipeline Name: %s\n Build Job name: %s\n, Build Name: %s\n\n Response : %s", os.Getenv("BUILD_ID"), os.Getenv("BUILD_PIPELINE_NAME"), os.Getenv("BUILD_JOB_NAME"), os.Getenv("BUILD_NAME"), response_data)
		err := sendSlackNotification(message)
		if err != nil {
			fmt.Printf("Error sending slack notification: %s\n", err)
		}
		return "Error Occured during change creation"
	} else {
		re := regexp.MustCompile(`CHG\d+`)
		extractedString := re.FindString(string(body))
		fmt.Println("Extracted String:", extractedString)
		return extractedString
	}
}

func sendSlackNotification(message string) error {
	slackBody, _ := json.Marshal(payload{Text: message})
	req, err := http.NewRequest(http.MethodPost, os.Getenv("TEST_SLACK"), bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 response: %s", resp.Status)
	}

	return nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	change_req_nubmber := createChangeReq()
	//time.Sleep(60 * time.Second)
	if strings.Contains(change_req_nubmber, "CHG") {
		fmt.Println(change_req_nubmber)
	} else {
		fmt.Println("Change creation failed")
	}
}
