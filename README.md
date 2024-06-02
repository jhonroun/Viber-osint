
### Об утилитке

Не люблю Python, не умею в Go ;)

Небольшое ПО, компилится, как и положено Go под всё, что угодно. Отлично рабоает в Termux.

Решает одну простую задачу: автоматизация отправки номера телефона в качестве контакта в канал в Viber. Если контакт существует в данном мессенджере, то автоматически добавляется автарка и подпись пользователя. 
Если повешать правильный веб-хук, то по изменению сведений о номере можно отловить и положить в базу, например sqlite.
В моем кейсе собрать аккаунты с ~ 17к номеров обошолся простым прогоном без базы.

построено на основе репозитория: https://github.com/sgxgsx/ViberOSINT.git

при пустом запуске даст пример использования или можно посмотреть в оригинальной репе. от себя добавил только проход по csv с номерами, без указания международного формата, т.е. без "+7".

### Оссновано на этой репе и пересписано (криво) на Go

```
https://github.com/sgxgsx/ViberOSINT.git
```

#### Как получить токен

* Создаем Viber аккаунт
* Создаем пустой канал 
* Тажимаем настройки канала, выбираем "Для разработчиков" получить токен. Сведения об API:  [https://developers.viber.com/docs/tools/channels-post-api/](https://developers.viber.com/docs/tools/channels-post-api/).
* Редактируем config.json меняем только [TOKEN] на полученый из предидущего шага.
* Запускаем

#### Что дальше?
А дальше, как выше писал, по описанию https://www.alexbilz.com/post/2021-01-29-forensic-artifacts-viber-desktop/ выгрузил историю канала из базы и объединил с данными CSV. На выходе получил структуру номер->ФИО->пользователь VIBER. Таким образом получилось установить админов некоторых интересных каналов и удиторию.

## Описание родного репозитория на Python.

Лень расписывать, кто хочет, тот поймет.

### About

This is a small script that helps you to automatically send phone numbers in bulk to your VIBER Channel as contacts and thus reveal whether a certain phone number is associated with a registered VIBER account.
There's no need in sharing your contacts with Viber.
Overall it helps you with your Open Source Intelligence (OSINT) workflow on Messangers - specifically Viber.

**It's not a full automation, you will need to do some things manually, but overall it saves more time by automating the most time consuming and dull task - adding phone numbers to the contact book**

### Based on this repo and rewrited to Go

```
https://github.com/sgxgsx/ViberOSINT.git
```

#### Token

* Create a VIBER account
* Create a VIBER channel
* Use the following documentation to find your API token - [https://developers.viber.com/docs/tools/channels-post-api/](https://developers.viber.com/docs/tools/channels-post-api/). You can find your token by entering your Channel’s info screen-> scroll down to “Developer Tools” -> copy token -> use the token for posting via API.
* Edit config.json in the repository and change [TOKEN] to your copied token
* Next time you are running the script, it'll automatically update the config file to include your UID and also would activate your channel api token by setting up a webhook.


### Examples


* send a single phone number to a viber channel as a contact **(this phone number has a registered Viber account)**

```

python3 viber_contacts.py --phone +79124538669

```

* send a single phone number to a viber channel as a contact **(this phone number doesn't have a registered Viber account associated with it)**

```
python3 viber_contacts.py --phone +79124538670

```


* send phones in bulk

```

python3 viber_contacts --list phones.txt

```



### After you ran a script:

* Open Viber Desktop and go to your VIBER channel
* Mostly always VIBER doesn't want to resolve contacts automatically, that's why you might need to do it manually or wait
* **Account is unresolved if its name is still "a"**

![alt text](https://github.com/sgxgsx/ViberOSINT/blob/main/images/notshown.png?raw=true)

* In order to manually resolve the accounts you need to:
* Click on the "a" name on a popup contact. It'll take you to another view where the name of this contact would change
* **If the name is a phone number and profile picture hasn't changed - it means that this phone number is not on Viber**

![alttext](https://github.com/sgxgsx/ViberOSINT/blob/main/images/notonviber.png?raw=true)

* **Otherwise, you'll see that the name of the contact is updated or you'll see a profile picture**

![alt text](https://github.com/sgxgsx/ViberOSINT/blob/main/images/onviber.png?raw=true)


