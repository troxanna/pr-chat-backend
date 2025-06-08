package bot

import (
 "bytes"
 "encoding/json"
 "fmt"
 "io"
 "mime/multipart"
 "net/http"
)

func sendVoiceInWhisper(file io.Reader, filename string) (string, error) {
 body := &bytes.Buffer{}
 writer := multipart.NewWriter(body)

 part, err := writer.CreateFormFile("file", filename)
 if err != nil {
  return "", err
 }

 _, err = io.Copy(part, file)
 if err != nil {
  return "", err
 }

 writer.Close()

 req, err := http.NewRequest("POST", "http://10.10.169.1:9000/transcribe/", body)
 if err != nil {
  return "", err
 }
 req.Header.Set("Content-Type", writer.FormDataContentType())

 client := &http.Client{}
 resp, err := client.Do(req)
 if err != nil {
  return "", err
 }
 defer resp.Body.Close()

 // Читаем тело ответа
 responseBody, err := io.ReadAll(resp.Body)
 if err != nil {
  return "", err
 }

 if resp.StatusCode >= 300 {
  return "", fmt.Errorf("ошибка от сервиса: %s\n%s", resp.Status, string(responseBody))
 }

 // Структура для JSON
 var result struct {
  Text string `json:"text"`
 }

 if err := json.Unmarshal(responseBody, &result); err != nil {
  return "", fmt.Errorf("не удалось распарсить JSON: %w", err)
 }

 return result.Text, nil
}