# avitoQA-task

Приветствую! Мое решение тестового задания направления QA

## Задание 1


Файл [task_1.md](/task_1.md) содержит решение первого тестового задания, в директории [task_1_screenshots](/task_1_screenshots) скрины к нему.


## Задание 2

В директории [task_1.md](/task_2) содержится файл [TESTCASES.md](/task_2/TESTCASES.MD) с тест-кейсами, тесты для эндпоинтов:

GET /api/1/item/{id}     

[get_item_by_id_test.go](/task_2/get_item_by_id_test.go)

GET /api/1/statistic/{id}  

[get_statistics_test.go](/task_2/get_statistics_test.go)

GET /api/1/{sellerID}/item 

[get_item_by_seller_test.go](/task_2/get_item_by_seller_test.go)

POST /api/1/item            
[save_item_test.go](/task_2/save_item_test.go)
    

И [BUGS.md](/task_2/BUGS.MD) с баг-репортами

## Инструкция по запуску тестов

1. Клонируйте репозиторий
2. Установите Go (при написании использовался версии 1.24.4) 
3. Перейдите в директорию /task_2. Запустите тесты:

go test -v
или go test -v -run TestCreateItem_Positive для конкретного теста