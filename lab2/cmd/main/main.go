package main

import (
	"fmt"
	"lab2/internal/domain"
)

func main() {
	//создаем банк
	bank := domain.NewBank("B001", "Мой Банк")
	fmt.Printf("Создан банк: %s (ID: %s)\n", bank.GetName(), bank.GetID())

	//создаем отделение банка
	branch := domain.NewBranch("BR001", "Главный офис", "ул. Ленина, 10", "+375298546688")
	bank.AddBranch(branch)
	fmt.Println("Добавлено отделение", branch.GetID(), "по адресу", branch.GetAddress())

	//создаем сотрудников
	emp1 := domain.NewEmployee("E001", "Иван", "Петров", "+375298654499", "Менеджер")
	emp2 := domain.NewEmployee("E002", "Анна", "Сидорова", "+375333789944", "Кассир")

	branch.AddEmployee(emp1)
	branch.AddEmployee(emp2)
	fmt.Printf("Добавлены сотрудники: %s и %s\n", emp1.GetFullName(), emp2.GetFullName())

	//создаем банкомат
	atm := domain.NewATM("ATM001", 1000000.00, branch.GetID())
	branch.AddATM(atm)
	fmt.Printf("Добавлен банкомат с балансом: %.2f бун.\n", atm.GetBalance())

	//создаем клиента
	client := domain.NewClient("C001", "Алексей", "Смирнов", "+375298456611")
	fmt.Println("Создан клиент:", client.GetFullName(), "телефон:", client.GetPhoneNum())

	//сткрываем счет клиенту
	account := domain.NewAccount("40817810000000000001", client.GetID(), "BYN", 50000.00)
	fmt.Printf("Открыт счет %s с балансом %.2f %s\n",
		account.GetNumber(), account.GetBalance(), account.GetCurrency())

	//создаем карту для счета
	card := domain.NewCard("1234567890123456", account.GetNumber(), 12, 2025, 123)
	account.AddCard(card)
	fmt.Printf("Выпущена карта %s для счета %s\n",
		card.GetCardNum(), card.GetAccountNum())

	//демонстрация работы с транзакциями
	transactions := []*domain.Transaction{
		domain.NewTransaction("T001", account.GetNumber(), "40817810000000000002", "BYN", 10000.00),
		domain.NewTransaction("T002", account.GetNumber(), "40817810000000000003", "BYN", 15000.00),
	}

	for _, t := range transactions {
		fmt.Printf("Транзакция %s: сумма %.2f %s, статус: %s\n",
			t.GetID(), t.GetMoney(), "BYN", t.GetStatus())

		if account.WithdrawMoney(t.GetMoney()) {
			t.Complete()
			fmt.Printf("Транзакция выполнена успешно\n")
		} else {
			t.Fail()
			fmt.Printf("Ошибка: недостаточно средств\n")
		}
	}

	//демонстрация работы с дженерик функцией SearchByID
	clients := []domain.Entity{client}
	employees := []domain.Entity{emp1, emp2}

	if found, ok := domain.SearchByID(clients, "C001"); ok {
		if c, ok := found.(*domain.Client); ok {
			fmt.Printf("Найден клиент: %s\n", c.GetFullName())
		}
	}

	if found, ok := domain.SearchByID(employees, "E002"); ok {
		if e, ok := found.(*domain.Employee); ok {
			fmt.Printf("Найден сотрудник: %s, должность: %s\n", e.GetFullName(), e.GetPosition())
		}
	}

	//демонстрация работы с депозитом
	deposit := domain.NewDeposit("D001", client.GetID(), 100000.00, "BYN", 5.5)
	fmt.Printf("Открыт депозит на сумму %.2f %s, ставка %.1f%%\n", deposit.GetBalance(), "BYN", 5.5)
	fmt.Printf("Проценты за период: %.2f %s\n", deposit.CalculateInterest(), "BYN")

	//демонстрация работы с кредитом
	loan := domain.NewLoan("L001", client.GetID(), "BYN", 300000.00)
	fmt.Printf("Выдан кредит на сумму %.2f %s\n", loan.GetRemainingDebt(), "BYN")

	PrintBalance := func (product domain.FinancialProduct) {
		fmt.Printf("Баланс: %.2f %s\n", product.GetBalance(), product.GetCurrency())
	}

	fmt.Print("Банковский счёт:")
	PrintBalance(account)
	fmt.Print("Кредит:")
	PrintBalance(loan)
	fmt.Print("Депозит:")
	PrintBalance(deposit)

	payment := 50000.00
	loan.MakePayment(payment)
	fmt.Printf("Внесен платеж %.2f %s, остаток долга: %.2f %s\n",
		payment, "BYN", loan.GetRemainingDebt(), "BYN")

	//демонстрация состояния банка
	fmt.Printf("Банк: %s\n", bank.GetName())
	fmt.Printf("Всего отделений: %d\n", len(bank.GetAllBranch()))

	for _, b := range bank.GetAllBranch() {
		fmt.Printf("  Отделение %s:\n", b.GetID())
		fmt.Printf("    - Адрес: %s\n", b.GetAddress())
		fmt.Printf("    - Телефон: %s\n", b.GetPhone())
		fmt.Printf("    - Открыто: %v\n", b.IsOpen())
		fmt.Printf("    - Сотрудников: %d\n", len(b.GetAllEmployee()))
		fmt.Printf("    - Банкоматов: %d\n", len(b.GetAllATM()))

		for i, emp := range b.GetAllEmployee() {
			status := "работает"
			if !emp.IsEmployed() {
				status = "не работает"
			}
			fmt.Printf("      %d) %s - %s (%s)\n", i+1, emp.GetFullName(), emp.GetPosition(), status)
		}
	}

	//демонстрация работы с картой
	fmt.Println("\nПроверка карты")
	if card.IsExpired() {
		fmt.Printf("Карта %s просрочена\n", card.GetCardNum())
	} else {
		fmt.Printf("Карта %s действительна\n", card.GetCardNum())
	}

	//демонстрация снятия наличных с банкомата
	withdrawAmount := 20000.00
	if atm.WithdrawMoney(withdrawAmount) {
		fmt.Printf("Снято %.2f бун. из банкомата. Остаток: %.2f бун.\n", withdrawAmount, atm.GetBalance())
	} else {
		fmt.Printf("Недостаточно средств в банкомате\n")
	}

	//демонстрация изменения данных
	fmt.Printf("Старый телефон клиента: %s\n", client.GetPhoneNum())
	client.SetPhone("+37529555555555")
	fmt.Printf("Новый телефон клиента: %s\n", client.GetPhoneNum())

	StateBranch := func() {
		if (branch.IsOpen()) {
			fmt.Printf("Отделение %s открыто\n", branch.GetID())
		} else {
			fmt.Printf("Отделение %s закрыто\n", branch.GetID())
		}
	}
	StateBranch()
	branch.SetClose()
	StateBranch()
	branch.SetOpen()
	StateBranch()
}
