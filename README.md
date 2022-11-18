# BalanceService

## Тестовое задание Avito_tech 
**Микросервис для работы с балансом пользователей**
Более детальное описание задания в файле [task.md](https://github.com/ksean42/BalanceService/blob/main/task.md)

**Основной функционал:**
- [x]  Метод начисления средств на баланс.
- [x]  Метод резервирования средств с основного баланса на отдельном счете.
- [x]  Метод признания выручки. – списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. 
- [x]  Метод получения баланса пользователя.

**Доп. функционал:**
- [x] Метод для получения месячного отчета в разрезе услуг.
- [x] Метод получения списка транзакций пользователя.
- [x] Метод перевода средств между пользователями.
- [x] Метод отмены резервирования денег.

**Запуск:**

````
docker-compose up
````

Сервер запускается на 8071 порту.

Конфигурация приложения в файле config.toml в корне проекта. (Для запуска сервера на локальной машине установите параметр db_host в значение "localhost")

**Swagger** документация к API доступна по адресу: http://localhost:8071/swagger/index.html

SQL файлы с созданием всех необходимых таблиц в БД в папке migrations

## Примеры запросов/ответов:

- **Метод начисления средств на баланс.**

Запрос:

````
curl -X POST localhost:8071/api/add 
-H "Content-Type: application/json"

-d '{
    "id": 1, 
    "amount": 500
}'
````
Ответ:

````
Status 200
{
    "Result":"OK"
}
````

Невалидные данные в полях:

Запрос:
````
curl -X POST localhost:8071/api/add  -H "Content-Type: application/json" -d '{ "id": 0, "amount": -500 }'
````
Ответ:
````
Status 400
{
    "Error":"id or amount is not correct"
}
````

- **Метод получения баланса пользователя:**

Запрос:

````
curl -X GET localhost:8071/api/balance -H "Content-Type: application/json" -d '{"id": 1 }'
````
Ответ:

````
Status 200
{
    "Result": 3270
}
````

Невалидное поле

Запрос:

````
curl -X GET localhost:8071/api/balance -H "Content-Type: application/json" -d '{"id": -1 }'
````
Ответ:

````
Status 400
{
    "Error":"id is incorrect"
}
````

Несуществующий id пользователя

Запрос:

````
curl -X GET localhost:8071/api/balance -H "Content-Type: application/json" -d '{"id": 31424131 }'
````
Ответ:

````
Status 400
{
    "Error":"not found""
}
````

- **Метод резервирования средств.** 

Принимает id пользователя, ИД услуги, ИД заказа, стоимость

Номер заказа на каждую транзакцию уникальный!

Запрос:

````
curl -X POST localhost:8071/api/reserve -H "Content-Type: application/json" -d '{"amount": 300, "id": 1, "order_id": 131, "service_id": 1 }'
````
Ответ:

````
Status 200
{
    "Result":"OK"
}
````

В случае если такой order_id уже существует

Ответ:

````
Status 400
{
    "Error":"order already exist"
}
````

Запрос большего кол-ва денег чем имеется у пользователя:

````
curl -X POST localhost:8071/api/reserve -H "Content-Type: application/json" -d '{"amount": 123300, "id": 1, "order_id": 132, "service_id": 1 }' 
````

Ответ:

````
Status 400
{
    "Error":"insufficient funds"
}
````

Запрос с невалидными данными:

````
curl -X POST localhost:8071/api/reserve -H "Content-Type: application/json" -d '{"amount": 10, "id": -1, "order_id": "dasdas", "service_id": 0 }'
````

Ответ:

````
Status 400
{
    "Error":"id is incorrect"
}
````
- **Метод признания выручки.**

Списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. Принимает id пользователя, ИД услуги, ИД заказа, сумму.

Запрос:

````
curl -X POST localhost:8071/api/approve -H "Content-Type: application/json" -d '{"amount": 300, "id": 1, "order_id": 132, "service_id": 1 }'
````

Ответ:

````
Status 200
{
    "Result":"OK"
}
````

Если резерв не найден:

````
curl -X POST localhost:8071/api/approve -H "Content-Type: application/json" -d '{"amount": 300, "id": 1, "order_id": 1321, "service_id": 1 }'
````

Ответ:

````
Status 400
{
    "Error":"not found"
}
````

Если сумма подтверждения отличается от суммы резерва:

````
curl -X POST localhost:8071/api/approve -H "Content-Type: application/json" -d '{"amount": 30, "id": 1, "order_id": 131, "service_id": 1 }' {"Error":"amount is incorrect"}````
````
Ответ:

````
Status 400
{
    "Error":"amount is incorrect"
}
````

- **Метод перечисления средств другому пользователю**

Принимает id отправителя, id получателя и сумму перевода.

Запрос:

````
curl -X POST localhost:8071/api/transfer -H "Content-Type: application/json" -d '{ "amount": 100,  "dest_id": 2, "src_id": 1}'
````
Ответ:

````
Status 200
{
    "Result":"OK"
}
````

При недостатке средств у отправителя:

````
Status 400
{
    "Error":"insufficient funds"
}
````

При невалидности полей запроса:
````
Status 400
{
    "Error":"invalid request, id and amount cant be less of equal than zero"
}
````

- **Метод отмены резерва**

Принимает ID пользователя и ID заказа.

Запрос:

````
curl -X POST localhost:8071/api/reject -H "Content-Type: application/json" -d '{"id": 1, "order_id": 131}'
````

Ответ:
````
Status 200
{
    "Result":"OK"
}
````

Обработка невалидных и несуществующих запросов - по аналогии с выше указанными примерами.


- **Метод для получения месячного отчета в разрезе услуг.**

Принимает месяц в формате "2022-11"

Возвращает путь к CSV файлу c отчетом

Запрос:

````
curl -X GET localhost:8071/api/report -H "Content-Type: application/json" -d '{"date": "2022-11"}'
````

Ответ:
````
Status 200
{
    "Result":"/Users/ksean/GolandProjects/avito_test_task/reports/revenue_report_2022-11.csv"}
}
````

**Пример отчета:**

````
service_id,revenue
1,3300
2,60
3,300

````

- **Метод получения списка транзакций пользователя.**

Принимает ID пользователя.
Возвращает 3 массива с данными о переводах другим пользователям, о подтвержденных транзакциях, о резервах.

Запрос:

````
curl -X GET localhost:8071/api/userReport -H "Content-Type: application/json" -d '{"id": 1}'````
````

Ответ:
````
Status 200
{
	"transfers": [
		{
			"src_id": 1,
			"dest_id": 5,
			"amount": 30,
			"date": "2022-11-17T22:14:13.082174Z"
		},
		{
			"src_id": 1,
			"dest_id": 2,
			"amount": 100,
			"date": "2022-11-18T04:02:21.415187Z"
		},
		{
			"src_id": 1,
			"dest_id": 1,
			"amount": 9,
			"date": "2022-11-18T04:04:35.10698Z"
		}
	],
	"transactions": [
		{
			"service_id": 2,
			"order_id": 2,
			"amount": 30,
			"date": "2022-11-17T03:06:56.15561Z"
		},
		{
			"service_id": 2,
			"order_id": 3,
			"amount": 30,
			"date": "2022-11-17T03:35:47.720993Z"
		},
		{
			"service_id": 1,
			"order_id": 132,
			"amount": 300,
			"date": "2022-11-18T03:54:37.939579Z"
		}
	],
	"reserves": [
		{
			"service_id": 5,
			"order_id": 23,
			"amount": 30
		},
		{
			"service_id": 3,
			"order_id": 234,
			"amount": 30
		},
		{
			"service_id": 5,
			"order_id": 22,
			"amount": 3000
		}
	]
}
````