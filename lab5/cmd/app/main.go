package main

import (
	"bufio"
	"fmt"
	"lab5/internal/domain"
	"lab5/internal/exporter"
	"lab5/internal/factory"
	"lab5/internal/filter"
	"lab5/internal/importer"
	"lab5/internal/service"
	"lab5/internal/sorter"
	"os"
	"strings"
)

func main() {
	// 1. Импорт данных через Фасад
	importer := importer.NewDataImporter()

	// Выбор файла для импорта
	fmt.Println("=== ЗАГРУЗКА ДАННЫХ ===")
	fmt.Println("Доступные файлы:")
	fmt.Println("1 - logistic.csv")
	fmt.Println("2 - logistic.json")
	fmt.Println("3 - logistic.xml")
	fmt.Print("Выберите файл (1-3): ")

	var fileChoice int
	fmt.Scanln(&fileChoice)

	var filePath string
	switch fileChoice {
	case 1:
		filePath = "../../materials/logistic.csv"
	case 2:
		filePath = "../../materials/logistic.json"
	case 3:
		filePath = "../../materials/logistic.xml"
	default:
		filePath = "../../materials/logistic.csv"
	}

	transports, importerCargoInfo, err := importer.ImportFromFile(filePath)
	if err != nil {
		fmt.Printf("Ошибка загрузки данных: %v\n", err)
		return
	}

	fmt.Printf("Загружено %d транспортов, %d типов грузов\n\n", len(transports), len(importerCargoInfo))

	// Преобразуем CargoInfo из importer в формат factory
	factoryCargoInfo := make(map[string]factory.CargoInfo)
	for k, v := range importerCargoInfo {
		factoryCargoInfo[k] = factory.CargoInfo{
			WeightPerKg:        v.WeightPerKg,
			TransportCostPerKg: v.TransportCostPerKg,
		}
	}

	// 2. Создаем фабрику и калькулятор
	logisticsFactory := factory.NewLogisticsFactory(transports, factoryCargoInfo)
	calculator := service.NewDeliveryCalculator(logisticsFactory)

	// 3. Выбор грузов
	fmt.Println("=== ФОРМИРОВАНИЕ ГРУЗОВ ===")
	fmt.Println("Доступные типы грузов:")
	fmt.Println("1 - Электроника")
	fmt.Println("2 - Одежда")
	fmt.Println("3 - Оборудование")
	fmt.Println("4 - Скоропортящиеся продукты")

	var cargoRequests []service.CargoRequest

	for {
		var cargoChoice, quantity int
		fmt.Print("\nВыберите тип груза (0 для завершения): ")
		fmt.Scanln(&cargoChoice)

		if cargoChoice == 0 {
			break
		}

		fmt.Print("Введите количество: ")
		fmt.Scanln(&quantity)

		var cargoType string
		switch cargoChoice {
		case 1:
			cargoType = domain.Electronic
		case 2:
			cargoType = domain.Clothing
		case 3:
			cargoType = domain.Equipment
		case 4:
			cargoType = domain.Perishable
		default:
			fmt.Println("Неверный выбор, попробуйте снова")
			continue
		}

		cargoRequests = append(cargoRequests, service.CargoRequest{
			Type:     cargoType,
			Quantity: quantity,
		})
	}

	if len(cargoRequests) == 0 {
		fmt.Println("Не выбрано ни одного груза. Выход.")
		return
	}

	// 4. Ввод расстояния
	var distance float64
	fmt.Print("\nВведите расстояние (км): ")
	fmt.Scanln(&distance)

	// 5. Выбор транспорта (с возможностью "все варианты")
	fmt.Println("\n=== ВЫБОР ТРАНСПОРТА ===")
	fmt.Println("0 - Все доступные варианты")
	fmt.Println("1 - Грузовик (Земля)")
	fmt.Println("2 - Поезд (Земля)")
	fmt.Println("3 - Танкер (Вода)")
	fmt.Println("4 - Самолет (Воздух)")
	fmt.Println("5 - Вертолет (Воздух)")

	var transportChoice int
	fmt.Print("Выберите транспорт (0-5): ")
	fmt.Scanln(&transportChoice)

	var results []domain.DeliveryResult

	if transportChoice == 0 {
		// Если транспорт не выбран — возвращаем все варианты
		fmt.Println("\nРасчет для всех доступных видов транспорта...")
		results, err = calculator.CalculateAllTransports(cargoRequests, distance)
		if err != nil {
			fmt.Printf("Ошибка расчета: %v\n", err)
			return
		}
	} else {
		// Конкретный транспорт
		var transportType, transportName string
		switch transportChoice {
		case 1:
			transportType = domain.Land
			transportName = "Грузовик (Земля)"
		case 2:
			transportType = domain.Land
			transportName = "Поезд (Земля)"
		case 3:
			transportType = domain.Water
			transportName = "Танкер (Вода)"
		case 4:
			transportType = domain.Air
			transportName = "Самолет (Воздух)"
		case 5:
			transportType = domain.Air
			transportName = "Вертолет (Воздух)"
		default:
			transportType = domain.Land
			transportName = "Грузовик (Земля)"
		}

		result, err := calculator.CalculateDelivery(cargoRequests, transportType, transportName, distance)
		if err != nil {
			fmt.Printf("Ошибка расчета: %v\n", err)
			return
		}
		results = append(results, *result)
	}

	// 6. Вывод результатов
	fmt.Println("\n=== РЕЗУЛЬТАТЫ ===")
	for i, r := range results {
		fmt.Printf("\n--- Вариант %d ---\n", i+1)
		fmt.Printf("Транспорт: %s (%s)\n", r.Transport.Name, r.Transport.Type)
		fmt.Printf("Скорость: %.1f км/ч\n", r.Transport.SpeedPerKmH)
		fmt.Printf("Стоимость перевозки грузов: %.2f руб.\n", r.TransportCost)
		fmt.Printf("Накладные расходы: %.2f руб.\n", r.DeliveryCost)
		fmt.Printf("ИТОГО: %.2f руб.\n", r.TotalCost)
		fmt.Printf("Время в пути: %s\n", r.FormatDuration())
	}

	// 7. Фильтрация (опционально)
	fmt.Println("\n=== ФИЛЬТРАЦИЯ ===")
	fmt.Print("Применить фильтр по максимальной цене? (y/n): ")
	var applyFilter string
	fmt.Scanln(&applyFilter)

	var filtered []domain.DeliveryResult = results

	if strings.ToLower(applyFilter) == "y" {
		var maxPrice float64
		fmt.Print("Введите максимальную цену: ")
		fmt.Scanln(&maxPrice)

		priceFilter := &filter.MaxPriceFilter{MaxPrice: maxPrice}
		filtered = priceFilter.Apply(filtered)

		fmt.Printf("После фильтрации осталось %d вариантов\n", len(filtered))
	}

	fmt.Print("Применить фильтр по названию транспорта? (y/n): ")
	var applyNameFilter string
	fmt.Scanln(&applyNameFilter)

	if strings.ToLower(applyNameFilter) == "y" {
		fmt.Print("Введите подстроку для поиска: ")
		reader := bufio.NewReader(os.Stdin)
		substring, _ := reader.ReadString('\n')
		substring = strings.TrimSpace(substring)

		nameFilter := &filter.TransportNameFilter{Substring: substring}
		filtered = nameFilter.Apply(filtered)

		fmt.Printf("После фильтрации осталось %d вариантов\n", len(filtered))
	}

	if len(filtered) == 0 {
		fmt.Println("Нет результатов после фильтрации. Выход.")
		return
	}

	// 8. Сортировка
	fmt.Println("\n=== СОРТИРОВКА ===")
	fmt.Println("Выберите тип сортировки (можно комбинировать через запятую, например: 1,2):")
	fmt.Println("1 - По названию транспорта")
	fmt.Println("2 - По цене (возрастание)")
	fmt.Println("3 - По скорости (убывание)")
	fmt.Print("Ваш выбор: ")

	var sortChoice string
	fmt.Scanln(&sortChoice)

	var strategies []sorter.SortStrategy
	choices := strings.Split(sortChoice, ",")

	for _, c := range choices {
		switch strings.TrimSpace(c) {
		case "1":
			strategies = append(strategies, &sorter.ByTransportName{})
		case "2":
			strategies = append(strategies, &sorter.ByPrice{})
		case "3":
			strategies = append(strategies, &sorter.BySpeed{})
		}
	}

	if len(strategies) > 0 {
		compositeSorter := sorter.NewCompositeSorter(strategies...)
		compositeSorter.Sort(filtered)

		fmt.Println("\nРезультаты после сортировки:")
		for i, r := range filtered {
			fmt.Printf("  %d. %-25s | Цена: %10.2f | Скорость: %6.1f км/ч | Время: %s\n",
				i+1, r.Transport.Name, r.TotalCost, r.Transport.SpeedPerKmH, r.FormatDuration())
		}
	}

	// 9. Экспорт
	fmt.Println("\n=== ЭКСПОРТ ===")
	fmt.Println("Выберите формат экспорта:")
	fmt.Println("1 - JSON")
	fmt.Println("2 - CSV")
	fmt.Println("3 - JSON + шифрование")
	fmt.Println("4 - JSON + ZIP")
	fmt.Println("5 - JSON + шифрование + ZIP")
	fmt.Println("6 - CSV + ZIP")
	fmt.Print("Ваш выбор: ")

	var exportChoice int
	fmt.Scanln(&exportChoice)

	var writer exporter.ExportWriter

	switch exportChoice {
	case 1:
		writer = &exporter.JSONWriter{}
	case 2:
		writer = &exporter.CSVWriter{}
	case 3:
		writer = exporter.NewEncryptionDecorator(
			&exporter.JSONWriter{},
			[]byte("my-secret-key12"),
		)
	case 4:
		writer = exporter.NewZipDecorator(&exporter.JSONWriter{})
	case 5:
		jsonWriter := &exporter.JSONWriter{}
		encryptedWriter := exporter.NewEncryptionDecorator(jsonWriter, []byte("my-secret-key12"))
		writer = exporter.NewZipDecorator(encryptedWriter)
	case 6:
		writer = exporter.NewZipDecorator(&exporter.CSVWriter{})
	default:
		writer = &exporter.JSONWriter{}
	}

	data, ext, err := writer.Write(filtered)
	if err != nil {
		fmt.Printf("Ошибка экспорта: %v\n", err)
		return
	}

	outputFile := "delivery_result" + ext
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Результат успешно сохранен в файл: %s\n", outputFile)
}
