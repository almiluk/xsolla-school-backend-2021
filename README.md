# xsolla-school-backend-2021
Тестовое задание для xsolla backend school 2021.  
Приложение является сервисом для управления (CRUD) товарами электронной торговой площадки.

## Задание
<details>
    <summary>Текст задания</summary>
    
### Обязательная часть
**Задача**: реализация системы управления товарами для площадки электронной коммерции (от англ. *e-commerce*).

Представь, что ты решил основать компанию. Вы занимаетесь реализацией решений, которые помогают разработчикам и издателям игр (ваша целевая аудитория). Главная задача представителей целевой аудитории — это продажа таких товаров, как игры, мерч, виртуальная валюта и др. Таким образом, ваша первая задача — дать возможность управлять товарами с помощью [RESTful API](https://searchapparchitecture.techtarget.com/definition/RESTful-API).

Для реализации прототипа системы напишите:
* Методы API для управления товарами — [операции CRUD](https://ru.wikipedia.org/wiki/CRUD). Товар определяется уникальным идентификатором и обязательно должен иметь [SKU](https://ru.wikipedia.org/wiki/SKU), имя, тип, стоимость. Предполагается наличие следующих [REST-методов](https://restfulapi.net/http-methods):
    * **Создание товара**. Метод генерирует и возвращает уникальный идентификатор товара.
    * **Редактирование товара**. Метод изменяет все данные о товаре по его идентификатору или SKU.
    * **Удаление товара по его идентификатору или SKU**.
    * **Получение информации о товаре по его идентификатору или SKU**.
    * <a name="get_items_catalog"></a>**Получение каталога товаров**. Метод возвращает список всех добавленных товаров.  
    Обратите внимание, что товаров может быть много. Необходимо предусмотреть возможность снижения нагрузки на сервис. **Вариант реализации**: возвращайте список товаров по частям.
* [Документацию в README](https://medium.com/xsolla-tech/tips-to-help-developer-improve-their-test-tasks-69d5a3b948d3). Обязательно укажите последовательность действий для запуска и локального тестирования API.

### Дополнительная часть
**Задача**: доработка системы управления товарами.

Мы предлагаем следующий список доработок:
* Фильтрация товаров по их типу и/или стоимости в методе получения каталога товаров.
* Спецификация OpenAPI [2.0](https://swagger.io/specification/v2/) или [3.0](https://swagger.io/specification/) (бывший Swagger). Документация для разработанного REST API.
* [Dockerfile](https://www.youtube.com/watch?v=QF4ZF857m44) для создания образа приложения системы. Желательно наличие файла Docker-compose.
* Модульные и функциональные тесты.
* Развертывание приложения на любом публичном хостинге, например, на [heroku](https://www.heroku.com/).

Выполнение пунктов из дополнительной части не является обязательным условием для прохождения тестового задания. Однако выполнение даже нескольких пунктов поможет нам составить более четкую картину о твоих знаниях и навыках.
 
</details>

### Реализовано
* API методы для операций CRUD
* Получение списка продуктов по частям
* Спецификация OpenAPI 2.0 (docs/swagger.*)
* Интерактивная документация swaggerUI


## Сборка
### Зависимости
* golang (v1.15)
* [github.com/gin-gonic/gin](github.com/gin-gonic/gin) v1.6.3
* [github.com/mattn/go-sqlite3](github.com/mattn/go-sqlite3) v1.14.7
* [github.com/swaggo/files](github.com/swaggo/files) v0.0.0-20190704085106-630677cd5c14
* [github.com/swaggo/gin-swagger](github.com/swaggo/gin-swagger) v1.3.0
* [github.com/swaggo/swag](github.com/swaggo/swag) v1.7.0

### Запуск сборки
    go build
    
    
## Запуск
### Windows
    XsollaSchoolBE.exe
    
### Linux
    ./XsollaSchoolBE
    
По-умолчанию приложение запускается в отладочном режиме (реализованно в github.com/gin-gonic/gin), для запуска в режиме релиза, установите значение переменной среды GIN_MODE равным "release".

## Описание API 
**После запуска приложения, по адресу [http://localhost:8080/swagger/index.html/](http://localhost:8080/swagger/index.html/) доступна интерактивная документация swaggerUI**.  
Также в папке /docs находятся фалы документации swagger 2.0: swagger.json и swagger.yaml 
     
Для запросов и ответов используются формат json, используются следующие основные объекты:
* Product - продукт находящийся в базе данных приложения:
```
{  
    "id": int64,  
    "sku": string,  
    "name": string,  
    "type": string,  
    "cost": uint32  
}
```
* InputProduct - продукт, добавляемый в базу данных приложения:  
```
{  
    "sku": string,  
    "name": string,  
    "type": string,  
    "cost": uint32  
}
```

### Методы API
* /products/
    * Метод GET. 
     
    Получение продукта с указанным SKU или ID, указанной части продуктов или всех продуктов.  
    URL query component параметры:  
    
    | Имя       | Тип    | Описание                                          |  
    |-----------|--------|---------------------------------------------------|  
    | sku       | string | sku искомого продукта                             |  
    | id        | int64  | id искомого продукта                              |  
    | groupSize | uint32 | Размер группы запрашиваемых продуктов             |  
    | groupNum  | uint32 | Номер запрашиваемой группы продуктов, начиная с 1 |  
    
    Использование параметров происходит в указанном в таблице порядке, т.е., если указан sku, выполняется поиск продукт с указанным sku, иначе аналогично для id, иначе для группы продуктов (в этом случае оба параметра groupSize и groupNum должны быть указаны), если не указан ни один параметр, метод вернёт все продукты.  
    Возможные ответы:  
    
    | Когда возвращается                       | Http код | Объект в теле ответа                                                                |
    |------------------------------------------|----------|-------------------------------------------------------------------------------------|
    | Успешное выполнение                      | 200      | **Массив** объектов Product (если запрашивался один продукт, массив из одного элемента) |
    | Некорректный запрос                      | 400      | string (описание ошибки)                                                            |
    | Продукт с указанным sku или id не найден | 404      | string (описание ошибки)                                                            |
    | Внутренняя ошибка сервера                | 500      | string (описание ошибки)                                                            |
      
    * Метод POST
  
    Добавление нового продукта.  
    Тело запроса - объект типа InputProduct с данными добавляемого продукта.  
    Возможные ответы:
    
    | Когда возвращается                        | Http код | Объект в теле ответа                            |
    |-------------------------------------------|----------|-------------------------------------------------|
    | Успешное выполнение                       | 201      | Product, описывающий добавленный продукт |
    | Некорректный запрос                       | 400      | string (описание ошибки)                 |
    | Продукт с таким SKU уже содержится в базе | 409      | string (описание ошибки и продукт в БД, вызвавший конфликт)
    | Внутренняя ошибка сервера                 | 500      | string (описание ошибки)                 |
    
    * Метод DELETE
    
    Удаление продукта из базы данных приложения.  
    URL query component параметры:  
    
    | Имя       | Тип    | Описание                                          |  
    |-----------|--------|---------------------------------------------------|  
    | sku       | string | sku искомого продукта                             |  
    | id        | int64  | id искомого продукта                              |  
  
    Если в запросе указан параметр sku, поиск продукта для удаления  будет производиться по нему, иначе по id, при этом, хотя бы один из двух параметров обязан присутствовать в запросе.  
    Возможные ответы:  
      
    | Когда возвращается                       | Http код | Объект в теле ответа                                                                |
    |------------------------------------------|----------|-------------------------------------------------------------------------------------|
    | Успешное выполнение                      | 200      | -                                                                            |
    | Некорректный запрос                      | 400      | string (описание ошибки)                                                                       |
    | Продукт с указанным sku или id не найден | 404      | string (описание ошибки)                                                                        |
    | Внутренняя ошибка сервера                | 500      | string (описание ошибки)                                                                        |
    
    * Метод PUT
    
    Замена всех данных (кроме id) об определённом продукте.
    Тело запроса - объект типа InputProduct с новыми данными продукта.  
    URL query component параметры:  
        
    | Имя       | Тип    | Описание                                          |  
    |-----------|--------|---------------------------------------------------|  
    | sku       | string | sku искомого продукта                             |  
    | id        | int64  | id искомого продукта                              |  
    
    Если в запросе указан параметр sku, поиск продукта для замены будет производиться по нему, иначе по id, при этом, хотя бы один из двух параметров обязан присутствовать в запросе.  
    Возможные ответы:  
    
    | Когда возвращается                       | Http код | Объект в теле ответа                                |
    |------------------------------------------|----------|-----------------------------------------------------|
    | Успешное выполнение                      | 200      | Product, описывающий продукт после изменения |
    | Некорректный запрос                      | 400      | string (описание ошибки)                                 |
    | Продукт с таким SKU уже содержится в базе| 409      | string (описание ошибки и продукт в БД, вызвавший конфликт) |
    | Продукт с указанным sku или id не найден | 404      | string (описание ошибки)                                        |
    | Внутренняя ошибка сервера                | 500      | string (описание ошибки)                                        |
       
* /products/{SKU}
    * Метод GET
    
    Получение продукта с указанным sku.  
    Возможные ответы:
    
    | Когда возвращается                       | Http код | Объект в теле ответа                                                                                          |
    |------------------------------------------|----------|---------------------------------------------------------------------------------------------------------------|
    | Успешное выполнение                      | 200      | Массив объектов Product, состояний из одного найденного продукта (для унификации типов возвращаемых значений) |
    | Продукт с указанным sku или id не найден | 404      | string (описание ошибки)                                                                                                         |
    | Внутренняя ошибка сервера                | 500      | string (описание ошибки)                                                                                                         |
    
    * Метод DELETE
    
    Удаление продукта, с указанным sku.  
    Возможные ответы:  
          
    | Когда возвращается                       | Http код | Объект в теле ответа                                                                |
    |------------------------------------------|----------|-------------------------------------------------------------------------------------|
    | Успешное выполнение                      | 200      | -                                                                            |
    | Продукт с указанным sku или id не найден | 404      | string (описание ошибки)                                                                        |
    | Внутренняя ошибка сервера                | 500      | string (описание ошибки)                                                                        |
    
    * Метод PUT 
      
    Тело запроса - объект типа InputProduct с новыми данными продукта.  
    Возможные ответы:  
        
    | Когда возвращается                       | Http код | Объект в теле ответа                                        |
    |------------------------------------------|----------|-------------------------------------------------------------|
    | Успешное выполнение                      | 200      | Product, описывающий продукт после изменения                |
    | Некорректный запрос                      | 400      | string (описание ошибки)                                    |
    | Продукт с указанным sku или id не найден | 404      | string (описание ошибки)                                    |
    | Продукт с таким SKU уже содержится в базе| 409      | string (описание ошибки и продукт в БД, вызвавший конфликт) |
    | Внутренняя ошибка сервера                | 500      | string (описание ошибки)                                        |
