package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type RateData struct {
	Kucoin string `json:"kucoin"`
	HTX    string `json:"htx"`
	Bybit  string `json:"bybit"`
}

type ExchangeInfo struct {
	Name string
	Rate float64
	Link string
}

var exchangeLinks = map[string]string{
	"kucoin": "https://www.kucoin.com/ru/otc/sell/USDT-RUB",
	"htx":    "https://www.htx.com/ru-ru/fiat-crypto/trade/sell-usdt-rub/",
	"bybit":  "https://www.bybit.com/ru-RU/fiat/trade/otc/?actionType=0&token=USDT&fiat=RUB&paymentMethod=581",
}

func parseRateString(s string) (float64, error) {
	s = strings.ReplaceAll(s, "RUB", "")
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")
	return strconv.ParseFloat(s, 64)
}

func getBestRate(jsonFilePath string) (ExchangeInfo, error) {
	// Открываем файл
	file, err := os.Open(jsonFilePath)
	if err != nil {
		return ExchangeInfo{}, err
	}
	defer file.Close()

	// Парсим JSON
	var data RateData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return ExchangeInfo{}, err
	}

	// Составим слайс доступных курсов
	exchanges := []ExchangeInfo{}

	// Добавим Kucoin
	if rate, err := parseRateString(data.Kucoin); err == nil {
		exchanges = append(exchanges, ExchangeInfo{
			Name: "kucoin",
			Rate: rate,
			Link: exchangeLinks["kucoin"],
		})
	}

	// Добавим HTX
	if rate, err := parseRateString(data.HTX); err == nil {
		exchanges = append(exchanges, ExchangeInfo{
			Name: "htx",
			Rate: rate,
			Link: exchangeLinks["htx"],
		})
	}

	// Добавим Bybit
	if rate, err := parseRateString(data.Bybit); err == nil {
		exchanges = append(exchanges, ExchangeInfo{
			Name: "bybit",
			Rate: rate,
			Link: exchangeLinks["bybit"],
		})
	}

	if len(exchanges) == 0 {
		return ExchangeInfo{}, fmt.Errorf("нет данных о курсах")
	}

	// Предположим, что лучший курс = самый высокий (обычно при продаже это выгоднее).
	best := exchanges[0]
	for _, ex := range exchanges {
		if ex.Rate > best.Rate {
			best = ex
		}
	}

	// Возвращаем структуру с лучшим курсом
	return best, nil
}

func main() {
	botToken := "7708967393:AAEXMPcBNgejHRQ_AjN6Ovru3oalsps2Fbs"
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			userMsg := update.Message.Text

			// Например, при команде /start отправим приветственное сообщение с кнопкой
			if userMsg == "/start" {
				msg := tgbotapi.NewMessage(chatID, "Привет! Нажми кнопку, чтобы узнать лучший курс.")

				// Создаём кнопки (ReplyKeyboardMarkup)
				keyboard := tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Получить лучший курс"),
					),
				)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
				continue
			}

			// Обработка нажатия кнопки
			if userMsg == "Получить лучший курс" {
				best, err := getBestRate("../parser/prices.json")
				if err != nil {
					bot.Send(tgbotapi.NewMessage(chatID, "Ошибка при получении курса: "+err.Error()))
					continue
				}

				// Формируем ответ
				// best.Rate - float64, нужно вывести с нужной точностью
				text := fmt.Sprintf(
					"Лучший курс сейчас: %.2f\nБиржа: %s\nСсылка: %s",
					best.Rate,
					strings.Title(best.Name),
					best.Link,
				)

				bot.Send(tgbotapi.NewMessage(chatID, text))
				continue
			}
		}
	}
}
