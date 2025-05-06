# biathlon-system-prototype
Прототип системы для соревнований по биатлону

## Инструкция по запуску

Для запуска программы и тестов используется `Makefile`. Убедитесь, что у вас установлен `make`.

### Сборка и запуск программы

Для сборки и запуска программы выполните команду:

```sh
make start
```

### Тесты

Для запуска тестов выполните команду:

```sh
make test
```

### Сборка программы

Для сборки программы выполните команду:

```sh
make build
```

### Запуск программы

Для запуска программы выполните команду:

```sh
make run
```

## Поясняющие комментарии

При выполнении данного задания обнаружилось достаточно много несоответствий между условиями задания и приведенным примером, поэтому приведу некоторые разнеснения моих рассуждений.

### Общее время

Если участник финишировал, то это интервал между:
- запланированным временем старта
- временем окончания прохождения последнего круга

В другом случае ставится отметка NotStarted или NotFinished.

### Время штрафного круга

Сумма временных интервалов между:
- заходом на штрафные круги
- выходом из штрафных кругов

### Огневые рубежи

Исходя их ТЗ поле FiringLines - это количество огневых рубежей на каждом основном круге.
Но в примере FiringLines - это общее количество огневых рубежей. 

Я решил использовать FiringLines как общее количество огневых рубежей.

### Округление

При вычислении средней скорости на основных и штрафных кругах я решил округлять тысячные в меньшуюю сторону (т.е 4.4356 = 4.435). 

### Время и скорость

На сколько я понял время и средняя скорость для основных кругов рассчитывается для каждого отдельно, а для штрафных вычисляется общая.

# Техническое задание

## Configuration (json)

- **Laps**        - Amount of laps for main distance
- **LapLen**      - Length of each main lap
- **PenaltyLen**  - Length of each penalty lap
- **FiringLines** - Number of firing lines per lap
- **Start**       - Planned start time for the first competitor
- **StartDelta**  - Planned interval between starts

## Events
All events are characterized by time and event identifier. Outgoing events are events created during program operation. Events related to the "incoming" category cannot be generated and are output in the same form as they were submitted in the input file.

- All events occur sequentially in time. (***Time of event N+1***) >= (***Time of event N***)
- Time format ***[HH:MM:SS.sss]***. Trailing zeros are required in input and output

#### Common format for events:
[***time***] **eventID** **competitorID** extraParams

```
Incoming events
EventID | extraParams | Comments
1       |             | The competitor registered
2       | startTime   | The start time was set by a draw
3       |             | The competitor is on the start line
4       |             | The competitor has started
5       | firingRange | The competitor is on the firing range
6       | target      | The target has been hit
7       |             | The competitor left the firing range
8       |             | The competitor entered the penalty laps
9       |             | The competitor left the penalty laps
10      |             | The competitor ended the main lap
11      | comment     | The competitor can`t continue
```
An competitor is disqualified if he/she does not start during his/her start interval. This marked as **NotStarted** in final report.
If the competitor can`t continue it should be marked in final report as **NotFinished**

```
Outgoing events
EventID | extraParams | Comments
32      |             | The competitor is disqualified
33      |             | The competitor has finished
```

## Final report
The final report should contain the list of all registered competitors
sorted by ascending time.
- Total time includes the difference between scheduled and actual start time or **NotStarted**/**NotFinished** marks
- Time taken to complete each lap
- Average speed for each lap [m/s]
- Time taken to complete penalty laps
- Average speed over penalty laps [m/s]
- Number of hits/number of shots

Examples:

`Config.conf`
```json
{
    "laps" : 2,
    "lapLen": 3651,
    "penaltyLen": 50,
    "firingLines": 1,
    "start": "09:30:00",
    "startDelta": "00:00:30"
}
```

`IncomingEvents`

```
[09:05:59.867] 1 1
[09:15:00.841] 2 1 09:30:00.000
[09:29:45.734] 3 1
[09:30:01.005] 4 1
[09:49:31.659] 5 1 1
[09:49:33.123] 6 1 1
[09:49:34.650] 6 1 2
[09:49:35.937] 6 1 4
[09:49:37.364] 6 1 5
[09:49:38.339] 7 1
[09:49:55.915] 8 1
[09:51:48.391] 9 1
[09:59:03.872] 10 1
[09:59:03.872] 11 1 Lost in the forest

```

`Output log`
```
[09:05:59.867] The competitor(1) registered
[09:15:00.841] The start time for the competitor(1) was set by a draw to 09:30:00.000
[09:29:45.734] The competitor(1) is on the start line
[09:30:01.005] The competitor(1) has started
[09:49:31.659] The competitor(1) is on the firing range(1)
[09:49:33.123] The target(1) has been hit by competitor(1)
[09:49:34.650] The target(2) has been hit by competitor(1)
[09:49:35.937] The target(4) has been hit by competitor(1)
[09:49:37.364] The target(5) has been hit by competitor(1)
[09:49:38.339] The competitor(1) left the firing range
[09:49:55.915] The competitor(1) entered the penalty laps
[09:51:48.391] The competitor(1) left the penalty laps
[09:59:03.872] The competitor(1) ended the main lap
[09:59:05.321] The competitor(1) can`t continue: Lost in the forest
```

`Resulting table`
```
[NotFinished] 1 [{00:29:03.872, 2.093}, {,}] {00:01:44.296, 0.481} 4/5