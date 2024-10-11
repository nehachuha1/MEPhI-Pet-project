![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) ![HTML5](https://img.shields.io/badge/html5-%23E34F26.svg?style=for-the-badge&logo=html5&logoColor=white) ![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white) ![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white) ![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white) ![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white) ![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)

# Пет-проект для НИЯУ МИФИ "Студенческое взаимодействие"

## Основная суть
Данный проект был создан в целях практики работы с различными технологиями и языком Golang (микросервисная архитектура на основе gRPC, взаимодействие с реляционной и key-value БД, создание docker-образов)

## Разделы
### - Маркетплейс (включает в себя продажу/покупку товаров, отслеживание заказов и т.д.)
### - Лента новостей (возможность публикации различных анонсов на главной странице)
### - Взаимодействие с со своим профилем и профилем пользователей

## Авторизация
В качестве алгоритма авторизации используется JWT-токен, хранящий в себе логин пользователя и ID сессии. Все сессии сохраняются в Redis и действительны в течение 3-ёх дней.

## Маркетплейс
### В маркетплейсе реализованы следующие возможности:
Выставление товаров (название, описание, цена, фотографии)
Возможность изменить данные товара, удалить его
Заказ товара и отслеживание его статуса в кабинете

### От лица продавца:
Возможность менять статус заказа, принимать его
Отменять заказ
