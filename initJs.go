package main

import (
	"bufio"
	"log"
	"os"
)

func initQuestionAnswerJs() string {
	js, err := os.Open("executionJs/question_answer.js")
	if err != nil {
		return ""
	}
	defer func(js *os.File) {
		err := js.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(js)
	read := bufio.NewScanner(js)
	return read.Text()
}

func initZhuanLanJs() string {
	js, err := os.Open("executionJs/question_answer.js")
	if err != nil {
		return ""
	}
	defer func(js *os.File) {
		err := js.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(js)
	read := bufio.NewScanner(js)
	return read.Text()
}
