package main

import (
	"fmt"
	"lab1/internal/domain"
)

func main() {
	//банк
	bank := domain.NewBank("1", "СберБанк")
	fmt.Println("Создан банк:", bank.GetName())

	//клиент
	client := domain.NewClient("1", "Иван", "Петров", "+375 29 855 26 97")
	fmt.Println("Создан клиент:", client.GetFullName())

	//счет
	account := domain.NewAccount("40817810000000000001", client.GetID(), "BYN", 50000.0)
	fmt.Println("Счет ", account.GetNumber(), ": баланс", account.GetBalance(), account.GetCurrency())

	//карта
	card := domain.NewCard("1234567890123456", account.GetNumber(), 12, 2025, 123)
	fmt.Println("Выпущена карта: ", card.GetCardNum())

	if !card.IsExpired() {
		fmt.Println("Карта действительна")
	}

	// Создание депозита
	deposit := domain.NewDeposit("DEP001", client.GetID(), 100000.0, "BYN", 8.5)
	fmt.Println("Начислены проценты по депозиту: ", deposit.CalculateInterest())

	//кредит
	loan := domain.NewLoan("LOAN001", client.GetID(), "BYN", 500000.0, 12)
	fmt.Println("Выдан кредит: ", loan.GetRemainingDebt(), "BYN")
	loan.MakePayment(15000.0)
	fmt.Println("Остаток по кредиту после платежа: ", loan.GetRemainingDebt())

	//отделения банка
	branch := domain.NewBranch("1", "Главный офис", "ул. Ленина, 1", "+375 29 955 28 97")
	fmt.Println("Отделение: ", branch.GetAddress())

	//банкомат
	atm := domain.NewATM("ATM001", 100000.0, "ATM009")
	fmt.Println("Банкомат ", atm.GetID(),": баланс ", atm.GetBalance())

	//транзакция
	transaction := domain.NewTransaction("TRX001", account.GetNumber(), "40817810000000000002", "BYN", 5000.0)
	fmt.Println("Создана транзакция со статусом:", transaction.GetStatus())

	//снятие со счета
	if account.WithdrawMoney(5000.0) {
		transaction.Complete()
		fmt.Println("Снятие успешно. Новый баланс:", account.GetBalance())
		fmt.Println("Статус транзакции:", transaction.GetStatus())
	} else {
		transaction.Fail()
		fmt.Println("Недостаточно средств")
	}

	//снятие с банкомата
	if atm.WithdrawMoney(10000.0) {
		fmt.Println("Снятие с банкомата успешно. Остаток в банкомате:", atm.GetBalance())
	}

	//сотрудник
	employee := domain.NewEmployee("1", "Анна Сидорова", "кассир")
	if employee.IsEmployed() {
		fmt.Println("Сотрудник", employee.GetFullName(),"работает", employee.GetPosition())
	} else {
		fmt.Println("Сотрудник", employee.GetFullName(),"уволен")
	}
	employee.StopWork()
	if employee.IsEmployed() {
		fmt.Println("Сотрудник", employee.GetFullName(),"должность:", employee.GetPosition())
	} else {
		fmt.Println("Сотрудник", employee.GetFullName(),"уволен")
	}
}
