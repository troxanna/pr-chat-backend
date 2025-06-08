package integration

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io"
	"log"
	"github.com/google/uuid"
)

var (
	uid = uuid.NewString()
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewClient(
	httpClient *http.Client,
	baseURL string,
	apiKey string,
) Client {
	return Client{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

type CompleteSessionRequest struct {
    OperatingSystemCode int    `json:"operatingSystemCode"`
    APIKey              string `json:"apiKey"`
    DialogIdentifier    string `json:"dialogIdentifier"`
}

type SendPromptForQuestionRequest struct {
    OperatingSystemCode int    `json:"operatingSystemCode"`
    APIKey              string `json:"apiKey"`
    UserDomainName      string `json:"userDomainName"`
    DialogIdentifier    string `json:"dialogIdentifier"`
    AIModelCode         int    `json:"aiModelCode"`
    Message             string `json:"Message"`
}

type GetResultForQuestionRequest struct {
    OperatingSystemCode int    `json:"operatingSystemCode"`
    APIKey              string `json:"apiKey"`
    DialogIdentifier    string `json:"dialogIdentifier"`
}


func (c Client) SendPromptForQuestion() {
	url := "https://gpt.orionsoft.ru/api/External/PostNewRequest"
	// url := c.baseURL

    // Создаем объект с данными (пока не прикрутим сюда данные от бота)
    requestData := SendPromptForQuestionRequest{
        OperatingSystemCode: 12,
        APIKey:              c.apiKey,
        UserDomainName:      "Team6QSXgoYCNNsG",
        DialogIdentifier:    uid,
        AIModelCode:         1,
        Message: `Сформулируй один открытый вопрос для собеседования, чтобы оценить уровень компетенции PosgreSQL у сотрудника. Уровень указан как {level} по следующей шкале:

0 — Нет желания изучать
1 — Нет экспертизы. Не изучал и не применял на практике
2 — Средняя экспертиза. Изучал самостоятельно, практики было мало
3 — Хорошая экспертиза. Регулярно применяет на практике
4 — Эксперт. Знает тонкости, делится лайфхаками
5 — Гуру. Готов выступать на конференциях

Построй вопрос так, чтобы он был релевантен именно для уровня 3 и позволял раскрыть глубину знаний сотрудника. Используй профессиональный стиль.`,
    }

    // Сериализация структуры в JSON
    jsonData, err := json.Marshal(requestData)
    if err != nil {
        panic(err)
    }

    // Создание POST-запроса
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Отправка запроса
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Чтение и вывод ответа
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	
	log.Println("Status:", resp.Status)
	log.Println("Response:", string(responseData))
}

func (c Client) GetResultForQuestionRequest() {
	url := "https://gpt.orionsoft.ru/api/External/GetNewResponse"

    // Данные запроса
    requestData := GetResultForQuestionRequest{
        OperatingSystemCode: 12,
        APIKey:              "OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ",
        DialogIdentifier:    uid,
    }

	jsonData, err := json.Marshal(requestData)
    if err != nil {
        panic(err)
    }

    // Создание HTTP-запроса
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Отправка запроса
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    // Чтение ответа
    responseData, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    log.Println("Status:", resp.Status)
    log.Println("Response:", string(responseData))

}

func (c Client) CleanContextRequest() {
	url := "https://gpt.orionsoft.ru/api/External/CompleteSession"

    // Данные запроса
    requestData := CompleteSessionRequest{
        OperatingSystemCode: 12,
        APIKey:              "OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ",
        DialogIdentifier:    uid,
    }

    jsonData, err := json.Marshal(requestData)
    if err != nil {
        panic(err)
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        panic(err)
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

	log.Println("Status:", resp.Status)
    log.Println("Response:", string(body))

}