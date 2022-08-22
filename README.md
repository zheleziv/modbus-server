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
Содержит в себе конфигурации для опроса множества устройств 
* ```name``` - имя устройства, должно быть уникальным
* ```ConnectionType``` - протокол подключения
* ```ID``` - для Modbus устройств id конечного устройства
* ```Log``` - в значении ```true``` предоставляет в консоли дополнительную информацию о состоянии подключения
* ```ConnectionTimeout``` - время между ```ConnectionAttempts``` попытками подключения в случае ошибки подключения 
* ```ConnectionAttempts``` - количество попыток подключения 
#### массив ```TAGS```
Содержит в себе конфигурации для опроса тэгов для каждого устройства
* ```name``` - имя тэга, должно быть уникальным внутри каждого устройства
* ```address``` - адресс тэга
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
## win_config.json
Файл с конфигурацией проверки полученных данных
* ```name``` - индентификатор тэга: имя узла + "." + имя тэга
* ```stateCondition``` - ```good``` или  ```bad``` условие состояние подключения 
* ```valueCondition``` - принимает массив условий (>```value```, <```value```, >=```value```, <=```value```, ==```value```, !=```value```), также доступны условия побитовых сравнений (условие, что 2 бит находится в сосотяние 1, будет выглядет следующим образом ```bit(2)``` или же ```!bit(2)```, если 2 бит должен находится в состоянии 0)
* ```logic``` - бинарный оператор, применяемый к двум вышеописанным условиям. Например, в конфигурационном файле мы указали
```JSON
{"NODES":[
    {
        "name":"UPS.ups_charge",
        "log": true,
        "stateCondition" : "bad",
        "valueCondition" : ">400",
        "logic": "and",
        "action": "shutdown",
        "actionTimeout": 10,
        "scanPeriod": 5
    },
]}
```
также мы имеем аналоговый регистр с значением 250 и хорошим сигналом. В ходе времени, значение регистра изменилось на 401, значит ```valueCondition" : ">400"``` становится ```true```, но выполнение программы не начинается так как у нас есть сигнал и условие ```"stateCondition" : "bad"``` = false. Затем сигнал пропал и началось выполнение команды ```"action": "shutdown"``` с задержкой ```"actionTimeout": 10```.
```JSON
{"NODES":[
    {
        "name":"UPS.ups_charge",
        "log": true,
        "stateCondition" : "good",
        "valueCondition" : "<10000,>400",
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
