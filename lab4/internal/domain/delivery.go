package domain

import (
	"fmt"
	"strings"
)

type DeliveryResult struct {
	CargoList     []Cargo
	Transport     Transport
	Distance      float64
	CargoCost     float64
	DeliveryCost  float64
	TotalCost     float64
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

	sb.WriteString("\nРЕЗУЛЬТАТ РАСЧЕТА ДОСТАВКИ\n")
	sb.WriteString(fmt.Sprintf("Способ доставки: %s\n", d.Transport.Name))
	sb.WriteString(fmt.Sprintf("Расстояние: %.1f км\n", d.Distance))
	sb.WriteString("\nГрузы\n")

	for _, cargo := range d.CargoList {
		sb.WriteString(fmt.Sprintf("  %s: %.1f кг (%d) - %.2f руб.\n",
			cargo.Name, cargo.TotalWeight(), cargo.Quantity, cargo.TotalCost()))
	}

	sb.WriteString(fmt.Sprintf("\nСтоимость\n"))
	sb.WriteString(fmt.Sprintf("  Перевозка грузов: %.2f руб.\n", d.CargoCost))
	sb.WriteString(fmt.Sprintf("  Накладные расходы: %.2f руб.\n", d.DeliveryCost))
	sb.WriteString(fmt.Sprintf("  ИТОГО: %.2f руб.\n", d.TotalCost))
	sb.WriteString(fmt.Sprintf("\nВремя\n"))
	sb.WriteString(fmt.Sprintf("  %s\n", d.FormatDuration()))

	return sb.String()
}
