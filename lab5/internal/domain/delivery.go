package domain

import (
	"fmt"
	"strings"
)

type DeliveryResult struct {
	CargoList     []Cargo
	Transport     Transport
	Distance      float64
	TransportCost float64 // ∑ c_i * m_i - стоимость перевозки грузов
	DeliveryCost  float64 // r * p_delivery - накладные расходы
	TotalCost     float64 // итоговая стоимость доставки
	DeliveryHours float64
}

func (d *DeliveryResult) FormatDuration() string {
	hours := int(d.DeliveryHours)
	mins := int((d.DeliveryHours - float64(hours)) * 60)

	if hours > 24 {
		days := hours / 24
		hours %= 24
		return fmt.Sprintf("%dд %dч", days, hours)
	}
	return fmt.Sprintf("%dч %dмин", hours, mins)
}

func (d *DeliveryResult) String() string {
	var sb strings.Builder

	sb.WriteString("         РЕЗУЛЬТАТ РАСЧЕТА ДОСТАВКИ\n")
	sb.WriteString(fmt.Sprintf("Способ доставки: %s\n", d.Transport.Name))
	sb.WriteString(fmt.Sprintf("Расстояние: %.1f км\n", d.Distance))
	sb.WriteString(fmt.Sprintf("Скорость: %.1f км/ч\n\n", d.Transport.SpeedPerKmH))

	sb.WriteString("ГРУЗЫ:\n")
	for _, cargo := range d.CargoList {
		sb.WriteString(fmt.Sprintf("  - %s\n", cargo.Name))
		sb.WriteString(fmt.Sprintf("    Количество: %d шт\n", cargo.Quantity))
		sb.WriteString(fmt.Sprintf("    Общий вес: %.1f кг\n", cargo.TotalWeight()))
		sb.WriteString(fmt.Sprintf("    Стоимость перевозки: %.2f руб.\n\n", cargo.TotalTransportCost()))
	}

	sb.WriteString("РАСЧЕТ СТОИМОСТИ:\n")
	sb.WriteString(fmt.Sprintf("  Перевозка грузов: %.2f руб.\n", d.TransportCost))
	sb.WriteString(fmt.Sprintf("  Накладные расходы (%d км × %.2f руб/км): %.2f руб.\n",
		int(d.Distance), d.Transport.CostPerKm, d.DeliveryCost))
	sb.WriteString(fmt.Sprintf("  ИТОГО к оплате: %.2f руб.\n\n", d.TotalCost))

	sb.WriteString(fmt.Sprintf("️  Время в пути: %s\n", d.FormatDuration()))

	return sb.String()
}
