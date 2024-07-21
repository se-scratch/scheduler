package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"scheduler/models"
	"strconv"
	"strings"
	"time"
)

const (
	DateLayout         = "20060102"
	repeatRulePrefixes = "dywm"
)

func NextDate(now time.Time, date, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	startTime, err := time.Parse(DateLayout, date)
	if err != nil {
		return "", fmt.Errorf("failed to parse start date: %w", err)
	}

	repeatRule := strings.Split(repeat, " ")
	repeatRulePrefix := repeatRule[0]
	if !strings.ContainsAny(repeatRulePrefix, repeatRulePrefixes) {
		return "", fmt.Errorf("invalid repeat rule: %s", repeatRule)
	}

	switch repeatRulePrefix {
	case "d":
		if len(repeatRule) != 2 {
			return "", fmt.Errorf("invalid repeat rule: %s", repeatRule)
		}
		intervalOfDays, err := strconv.Atoi(repeatRule[1])
		if err != nil {
			return "", fmt.Errorf("interval of days must be  number: %s", repeatRule)
		}
		if intervalOfDays > 400 {
			return "", fmt.Errorf("interval of days must not be greater than 400: %d", intervalOfDays)
		}
		for {
			startTime = startTime.AddDate(0, 0, intervalOfDays)
			if startTime.After(now) {
				return startTime.Format(DateLayout), nil
			}
		}
	case "w":
		return "", fmt.Errorf("invalid repeat rule: %s", repeatRule)
	case "m":
		return "", fmt.Errorf("invalid repeat rule: %s", repeatRule)
	case "y":
		for {
			startTime = startTime.AddDate(1, 0, 0)
			if startTime.After(now) {
				return startTime.Format(DateLayout), nil
			}
		}
	}
	return "", nil
}

func DecodePostTask(buf io.Reader) (models.Task, error) {
	var task models.Task
	if err := json.NewDecoder(buf).Decode(&task); err != nil {
		return task, fmt.Errorf("unmarshal json error: %w", err)
	}
	if task.Title == "" {
		return task, fmt.Errorf("task title is empty")
	}

	now := time.Now()
	localZeroTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	if task.Date == "" {
		task.Date = localZeroTime.Format(DateLayout)
	}

	date, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return task, fmt.Errorf("parse date error: %w", err)
	}

	if date.Before(localZeroTime) {
		if task.Repeat == "" {
			task.Date = localZeroTime.Format(DateLayout)
		} else {
			task.Date, err = NextDate(localZeroTime, task.Date, task.Repeat)
			if err != nil {
				return task, err
			}
		}
	}
	return task, nil
}
