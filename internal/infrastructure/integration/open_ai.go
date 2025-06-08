package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ContextItem struct {
	RequestMessage  string `json:"requestMessage"`
	RequestTime     string `json:"requestTime"`
	ResponseMessage string `json:"responseMessage"`
	ResponseTime    string `json:"responseTime"`
}

type Data struct {
	UserCrocCode     string        `json:"userCrocCode"`
	DialogIdentifier string        `json:"dialogIdentifier"`
	AIModelCode      int           `json:"aiModelCode"`
	AIModelName      string        `json:"aiModelName"`
	Warning          int           `json:"warning"`
	LastMessage      string        `json:"lastMessage"`
	LastResponseTime string        `json:"lastResponseTime"`
	Context          []ContextItem `json:"context"`
}

type Status struct {
	IsSuccess   bool   `json:"isSuccess"`
	Description string `json:"description"`
}

type FullResponse struct {
	Status Status `json:"status"`
	Data   Data   `json:"data"`
}

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

func (c Client) SendPromptForQuestion(uid string) {
	url := "https://gpt.orionsoft.ru/api/External/PostNewRequest"
	// url := c.baseURL

	// Создаем объект с данными (пока не прикрутим сюда данные от бота)
	requestData := SendPromptForQuestionRequest{
		OperatingSystemCode: 12,
		APIKey:              "OrVrQoQ6T43vk0McGmHOsdvvTiX446RJ",
		UserDomainName:      "Team6QSXgoYCNNsG",
		DialogIdentifier:    uid,
		AIModelCode:         1,
		Message: `Сформулируй один открытый вопрос для собеседования, чтобы оценить уровень компетенции PosgreSQL у сотрудника. Уровень указан как 2 по следующей шкале:

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

	log.Println(req)

	// Отправка запроса
	resp, err := c.httpClient.Do(req)
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

func (c Client) GetResultForQuestionRequest(uid string) bool {
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
	// client := &http.Client{}
	log.Println(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Println(string(responseData))

	var result FullResponse

	if err := json.Unmarshal(responseData, &result); err != nil {
		log.Printf("Ошибка разбора JSON: %v", err)
		log.Println("Сырой ответ:", string(responseData))
		return false
	}

	if result.Data.LastMessage == "" {
		return false
	}

	log.Println("✅ Успешность:", result.Status.IsSuccess)
	log.Println("📄 Описание:", result.Status.Description)
	log.Println("👤 UserCrocCode:", result.Data.UserCrocCode)
	log.Println("🧾 DialogIdentifier:", result.Data.DialogIdentifier)
	log.Println("🤖 Модель:", result.Data.AIModelName)
	log.Println("🗨️ Последнее сообщение:", result.Data.LastMessage)
	log.Println("⏱ Время последнего ответа:", result.Data.LastResponseTime)
	log.Println("📚 Контекст сообщений:")

	for i, item := range result.Data.Context {
		log.Printf("  [%d] ▶️ Запрос (%s): %s", i+1, item.RequestTime, item.RequestMessage)
		log.Printf("      💬 Ответ  (%s): %s", item.ResponseTime, item.ResponseMessage)
	}
	return true

}

func (c Client) CleanContextRequest(uid string) {
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

	log.Println(req)
	resp, err := c.httpClient.Do(req)
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
