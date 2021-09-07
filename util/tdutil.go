package util

import (
	"time"
)

// GetDate 
func GetDate() string {
	const formatDate = "2006-01-02"
	aTime := time.Now()
	return aTime.Format(formatDate)
}

// GetDateTime
func GetDateTime() string {
	const formatDateTime = "2006-01-02 15:04:05"
	aTime := time.Now()
	return aTime.Format(formatDateTime)
}

// GetLastDateKeyHash Chave de Hash com a data anterior
func GetAfterDate() string {
	const formatDate = "2006-01-02"
	var time1 = time.Now()
	time1 = time1.AddDate(0, 0, +1)
	return time1.Format(formatDate)
}
