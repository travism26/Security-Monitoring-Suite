package sender

import (
    "bytes"
    "net/http"
    "log"
)

// SendData sends collected data to the web dashboard.
func SendData(data interface{}) {
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Printf("Error marshalling data: %v", err)
        return
    }

    resp, err := http.Post("http://your-dashboard-url.com/metrics", "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        log.Printf("Error sending data: %v", err)
        return
    }
    defer resp.Body.Close()

    log.Printf("Data sent successfully, status code: %d", resp.StatusCode)
}
