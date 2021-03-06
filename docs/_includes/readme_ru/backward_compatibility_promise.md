<!-- THIS FILE IS AUTOGENERATED BY werf docs COMMAND! DO NOT EDIT! -->

> _Note:_ Настоящее обещание относится к werf, начиная с версии 1.0, и не относится к предыдущим версиям или версиям dapp

werf использует [семантическое версионирование](https://semver.org/lang/ru/). Это значит, что мажорные версии (1.0, 2.0) могут быть обратно не совместимыми между собой. В случае werf это означает, что обновление на следующую мажорную версию _может_ потребовать полного передеплоя приложений, либо других ручных операций.

Минорные версии (1.1, 1.2, etc) могут добавлять новые "значительные" изменения, но без существенных проблем обратной совместимости в пределах мажорной версии. В случае werf это означает, что обновление на следующую минорную версию в большинстве случаев будет беспроблемным, но _может_ потребоваться запуск предоставленных скриптов миграции.

Патч-версии (1.1.0, 1.1.1, 1.1.2) могут добавлять новые возможности, но без каких-либо проблем обратной совместимости в пределах минорной версии (1.1.x).
В случае werf это означает, что обновление на следующий патч (следующую патч-версию) не должно вызывать проблем и требовать каких-либо ручных действий.

Все изменения проходят полный цикл по каналам стабильности:

- Канал обновлений `alpha` может содержать новые возможности и быть нестабильным. Релизы выполняются с высокой периодичностью.
  Мы **не гарантируем** обратную совместимость между версиями канала обновлений `alpha`.
- Канал обновлений `beta` предназначен для более детального тестирования новых возможностей.
  Мы **не гарантируем** обратную совместимость между версиями канала обновлений `beta`.
- Канал обновлений `ea` безопасно использовать в некритичных окружениях и при локальной разработке.
  Мы **не гарантируем** обратную совместимость между версиями канала обновлений `ea`.
- Канал обновлений `stable` считается безопасным и рекомендуемым для всех окружений.
  Мы **гарантируем**, что версия канала обновлений `ea` перейдет в канал обновлений `stable` не ранее чем через неделю после внутреннего тестирования.
  Мы **гарантируем** обратную совместимость между версиями канала обновлений `stable` в пределах минорной версии (1.1.x).
- Канал обновлений `rock-solid` рекомендуется использовать в критичных окружениях с высоким SLA.
  Мы **гарантируем**, что версия из канала обновлений `stable` перейдет в канал обновлений `rock-solid` не ранее чем через 2 недели плотного тестирования.
  Мы **гарантируем** обратную совместимость между версиями канала обновлений `rock-solid` в пределах минорной версии (1.1.x).

Соответствие каналов и релизов описывается в файле [multiwerf.json](https://github.com/werf/werf/blob/multiwerf/multiwerf.json), а использование актуальной версии werf в рамках канала должно быть организовано с помощью утилиты [multiwerf](https://github.com/werf/multiwerf).

Каналы стабильности и частые релизы позволяют получать непрерывную обратную связь по новым изменениям, выполнять быстрый откат проблемных изменений, а также обеспечивать высокую степень стабильности и при этом приемлемую скорость разработки.
