package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

type Prices struct {
	Kucoin string `json:"kucoin"`
	HTX    string `json:"htx"`
	Bybit  string `json:"bybit"`
}

func main() {
	// Настройка параметров браузера
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Убираем headless режим для видимого окна
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("start-maximized", true),
	)

	for {
		// Создание и выделение ресурсов для браузера
		allocatorCtx, cancelAllocator := chromedp.NewExecAllocator(context.Background(), opts...)
		ctx, cancel := chromedp.NewContext(allocatorCtx)

		// Таймаут для контекста
		ctx, cancelTimeout := context.WithTimeout(ctx, 500*time.Second)

		// Переменные для хранения данных
		var kucoinPrice, htxPrice, bybitPrice string

		// Работа с Bybit
		err := fetchBybit(ctx, &bybitPrice)
		if err != nil {
			log.Printf("Error fetching Bybit: %v", err)
		} else {
			log.Printf("Bybit Price: %s", bybitPrice)
		}
		// Работа с KuCoin
		err = fetchKucoin(ctx, &kucoinPrice)
		if err != nil {
			log.Printf("Error fetching KuCoin: %v", err)
		} else {
			log.Printf("KuCoin Price: %s", kucoinPrice)
		}

		// Работа с HTX
		err = fetchHTX(ctx, &htxPrice)
		if err != nil {
			log.Printf("Error fetching HTX: %v", err)
		} else {
			log.Printf("HTX Price: %s", htxPrice)
		}

		// Преобразование данных
		kucoinPrice = processPrice(kucoinPrice)
		htxPrice = processPrice(htxPrice)
		bybitPrice = processPrice(bybitPrice)

		// Создание JSON
		prices := Prices{
			Kucoin: kucoinPrice,
			HTX:    htxPrice,
			Bybit:  bybitPrice,
		}
		err = savePricesToJSON("prices.json", prices)
		if err != nil {
			log.Fatalf("Failed to save JSON: %v", err)
		}

		fmt.Printf("Prices saved to JSON:\nKuCoin: %s\nHTX: %s\nBybit: %s\n", kucoinPrice, htxPrice, bybitPrice)

		// Завершаем текущий контекст вручную
		cancelTimeout()
		cancel()
		cancelAllocator()

		// Таймаут перед повторением
		log.Println("Waiting for 30 minutes before the next iteration...")
		time.Sleep(30 * time.Minute)
	}
}

func processPrice(price string) string {
	price = strings.TrimSpace(price)
	price = strings.ReplaceAll(price, "RUB", "")
	price = strings.ReplaceAll(price, " ", "")
	price = strings.ReplaceAll(price, ".", ",")
	return price
}

func savePricesToJSON(filename string, prices Prices) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(prices)
}

func fetchKucoin(ctx context.Context, price *string) error {
	url := "https://www.kucoin.com/ru/otc/sell/USDT-RUB"
	return chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#root > div > div > div.page-body.body_hXuuO > div > div > div.container_aQpvt.containerOrderList_BXxhh > div > div > div > div.topFilterContainer_y7A5l > div.orderListHeader_lgyW0 > div.filterContainer_v4W5Z > div.lrtcss-1hmera3 > div > div > div`),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#root > div > div > div.page-body.body_hXuuO > div > div > div.container_aQpvt.containerOrderList_BXxhh > div > div > div > div.topFilterContainer_y7A5l > div.orderListHeader_lgyW0 > div.filterContainer_v4W5Z > div.lrtcss-1hmera3 > div > div.KuxDropDown-popper.KuxDropDown-open > div > div > div.lrtcss-6u20ey > div:nth-child(7)`),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#root > div > div > div.page-body.body_hXuuO > div > div > div.container_aQpvt.containerOrderList_BXxhh > div > div > div > div.topFilterContainer_y7A5l > div.orderListHeader_lgyW0 > div.filterContainer_v4W5Z > div.lrtcss-1hmera3 > div > div.KuxDropDown-popper.KuxDropDown-open > div > div > div.lrtcss-lfpii1 > button.KuxButton-root.KuxButton-contained.KuxButton-containedPrimary.KuxButton-sizeBasic.KuxButton-containedSizeBasic.lrtcss-oiu68b`),
		chromedp.Sleep(3*time.Second),
		chromedp.SendKeys(`#root > div > div > div.page-body.body_hXuuO > div > div > div.container_aQpvt.containerOrderList_BXxhh > div > div > div > div.topFilterContainer_y7A5l > div.orderListHeader_lgyW0 > div.filterContainer_v4W5Z > div.amountInput_Te8IU > div > input`, "10000"),
		chromedp.Sleep(3*time.Second),
		chromedp.Text(`#root > div > div > div.page-body.body_hXuuO > div > div > div.container_aQpvt.containerOrderList_BXxhh > div > div > div > div.orderListTableWrapper_H7Sii > div > div > div > div > div > div > table > tbody > tr:nth-child(1) > td.priceColumn_fkDBG.lrtcss-fnv5fp > div > span:nth-child(1)`, price),
		chromedp.Sleep(3*time.Second),
	)
}

func fetchHTX(ctx context.Context, price *string) error {
	url := "https://www.htx.com/ru-ru/fiat-crypto/trade/sell-usdt-rub/"
	return chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`body > div.hb-modal.en > div.hb-dialog__wrapper.pc-modal > div > div > span`), // Закрытие popup
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#app > div > div.flex-auto > div > div.trade-responsive > div.bg-grey > div.trade-list-title.flex-between-center > div.flex-y-center.pc-common-list > div:nth-child(2) > div > div > div.refresh-trade-list.refresh-filtrate > div`), // Клик на Payments
		chromedp.Sleep(3*time.Second),
		chromedp.Evaluate(`document.querySelector('#app > div > div.flex-auto > div > div.trade-responsive > div.bg-grey > div.trade-list-title.flex-between-center > div.flex-y-center.pc-common-list > div:nth-child(2) > div.flex-y-center.spe-flex-wrap > div > div.refresh-trade-list.refresh-filtrate > div.setting-inner-new.filtrateSearch.font12').scrollBy(0, 450);`, nil),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#app > div > div.flex-auto > div > div.trade-responsive > div.bg-grey > div.trade-list-title.flex-between-center > div.flex-y-center.pc-common-list > div:nth-child(2) > div.flex-y-center.spe-flex-wrap > div > div.refresh-trade-list.refresh-filtrate > div.setting-inner-new.filtrateSearch.font12 > div:nth-child(5) > div.filtrate-item-content > ul > li:nth-child(6)`),
		chromedp.Sleep(3*time.Second),
		chromedp.SendKeys(`#app > div > div.flex-auto > div > div.trade-responsive > div.bg-grey > div.trade-list-title.flex-between-center > div.flex-y-center.pc-common-list > div:nth-child(2) > div > div > div.search-amount-fiat-container.refresh-setting > div.input-container.flex-y-center > div > div > input`, "10000"),
		chromedp.Sleep(5*time.Second),
		chromedp.Text(`#app > div > div.flex-auto > div > div.trade-responsive > div.bg-grey > div.trade-content > div.content > div:nth-child(1) > div > div > div > div.width210.price.average.mr-24 > div`, price),
		chromedp.Sleep(3*time.Second),
	)
}

func fetchBybit(ctx context.Context, price *string) error {
	url := "https://www.bybit.com/ru-RU/fiat/trade/otc/?actionType=0&token=USDT&fiat=RUB&paymentMethod=581"
	return chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`button.ant-btn.css-7o12g0.ant-btn-primary.css-7o12g0.ant-btn-custom.ant-btn-custom-middle.ant-btn-custom-primary.bds-theme-component-light`, chromedp.NodeVisible),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#paywayAnchorList`),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#lists > li:nth-child(5)`),
		chromedp.Sleep(3*time.Second),
		chromedp.Click(`#paywayList > div > div > section > button.by-button.btn-confirm`),
		chromedp.Sleep(5*time.Second),
		chromedp.SendKeys(`#guide-step-two > div:nth-child(1) > div > div:nth-child(1) > input`, "10000"),
		chromedp.Sleep(5*time.Second),
		chromedp.Text(`#root > div.trade-list > div.trade-list__main > div.trade-list__wrapper > div.trade-list__content > div > div > div > table > tbody.trade-table__tbody > tr:nth-child(2) > td:nth-child(2) > div > div`, price),
		chromedp.Sleep(3*time.Second),
	)
}
