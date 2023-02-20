package main

import (
	"os"
)

func initQuestionAnswerJs() string {
	js, err := os.ReadFile("executionJs/question_answer.js")
	if err != nil {
		return ""
	}
	return string(js)
}

func initZhuanLanJs() string {
	js, err := os.ReadFile("executionJs/question_answer.js")
	if err != nil {
		return ""
	}
	return string(js)
}
