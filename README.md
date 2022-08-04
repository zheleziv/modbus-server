# modbus-server
## Как установить
* Скачать репозиторий
* Открыть консоль
* Перейти в папку ./cmd/app/ ```(cd PATH/cmd/app/)```
* Запустить команду  ```(go build .)```
## Конфигурационные файлы
Файлы с конфигурацией приложения должны располагаться в директории приложения ```app.exe```
## config.json
Файл с конфигурацией сервера опроса устройств 
#### массив ```NODES```
содержит в себе конфигурации для опроса множества устройств 
* ```name``` - имя устройства, должно быть уникальным
* ```ConnectionType``` - протокол подключения
* ```ID``` - для Modbus устройств id конечного устройства
* ```Log``` - в значении ```true``` предоставляет в консоли дополнительную информацию о состоянии подключения
* ```ConnectionTimeout``` - время между ```ConnectionAttempts``` попытками подключения в случае ошибки подключения 
* ```ConnectionAttempts``` - количество попыток подключения 
#### массив ```TAGS```
содержит в себе конфигурации для опроса тэгов для каждого устройства
* ```name``` - имя тэга, должно быть уникальным внутри каждого устройства
* ```dataType``` 
  * ```"coil"``` - дискрет
  * ```"word"``` - 16-битный аналог
  * ```"dword"``` - 32-битный аналог 
* ```scanPeriod``` - время между запросами
##### Пример конфига для опроса двух устройств с разным количеством тэгов
```JSON
{
"NODES":[
    {
        "name" : "UPS",
        "ConnectionType": "modbusTCP",
        "IP": "127.0.0.1",
        "port":502,
        "ID":1,
        "Log": false,
        "ConnectionTimeout": 10,
        "ConnectionAttempts": 5,
        "TAGS": [
            {
                "name": "ups_charge",
                "address":400008,
                "dataType": "word",
                "scanPeriod":10
            },
            {
                "name": "ups_state",
                "address": 1000001,
                "dataType": "coil",
                "scanPeriod" : 5
            }
        ]
    },
    {
        "name" : "testNode2",
        "ConnectionType": "modbusTCP",
        "IP": "172.32.32.32",
        "port":502,
        "ID":1,
        "Log": false,
        "ConnectionTimeout": 10,
        "ConnectionAttempts": 5,
        "TAGS": [
            {
                "name": "testTag1",
                "address":400008,
                "dataType": "dword",
                "scanPeriod":10
            }
        ]
    }
]
}
```
## config.json
Файл с конфигурацией проверки полученных данных
```JSON
{"NODES":[
    {
        "name":"UPS.ups_charge",
        "log": true,
        "stateCondition" : "good",
        "valueCondition" : "<10000",
        "logic": "and",
        "action": "shutdown",
        "actionTimeout": 10,
        "scanPeriod": 5
    },
    {
        "name":"UPS.ups_state",
        "log": true,
        "stateCondition" : "bad",
        "valueCondition" : "false",
        "logic": "and",
        "action": "restart",
        "actionTimeout": 10,
        "scanPeriod": 5
    },
    {
        "name":"testNode2.testTag1",
        "log": true,
        "stateCondition" : "good",
        "valueCondition" : ">=23324",
        "logic": "and",
        "action": "run __debug_bin.exe",
        "actionTimeout": 10,
        "scanPeriod": 5
    }
]}
```
