# Implementation Plan: Устранение глобальных переменных конфигурации

## Overview

Данный план описывает пошаговую реализацию рефакторинга для устранения глобальных переменных `config.Camera`, `config.Game` и `config.Physics`. Рефакторинг применяет паттерн Dependency Injection, передавая конфигурации явно через конструкторы компонентов. Каждая задача строится на предыдущих, обеспечивая инкрементальный прогресс без нарушения работоспособности кода.

## Tasks

- [x] 1. Подготовка: Сохранение старой реализации для тестирования
  - Создать копию функции `render.Project()` как `render.ProjectOld()` для property-based тестов
  - Эта копия будет использоваться для сравнения результатов после рефакторинга
  - _Requirements: 3.4, 8.3_

- [x] 2. Рефакторинг структуры Camera
  - [x] 2.1 Добавить поле config в структуру Camera
    - Добавить поле `config config.CameraConfig` в структуру `Camera` в `internal/render/camera.go`
    - Обновить конструктор `NewCamera` для приёма параметра `cfg config.CameraConfig`
    - Сохранить переданную конфигурацию в поле `config`
    - Использовать `cfg.DefaultPositionY` и `cfg.DefaultPositionZ` для инициализации позиции камеры
    - _Requirements: 2.1, 2.4, 3.2_

  - [x] 2.2 Реализовать метод Camera.Project()
    - Создать метод `func (c *Camera) Project(point core.Vec3) (float64, float64, float64)`
    - Скопировать логику из функции `render.Project()` в новый метод
    - Заменить все обращения к `config.Camera` на `c.config`
    - Использовать поля камеры `c.screenW`, `c.screenH`, `c.focalLength`, `c.horizonY`, `c.position`
    - _Requirements: 3.1, 3.2, 3.3, 3.4_

  - [ ]* 2.3 Написать property-based тест для эквивалентности проекции
    - **Property 1: Projection calculation equivalence**
    - **Validates: Requirements 3.4, 8.3**
    - Использовать библиотеку gopter для генерации случайных входных данных
    - Генерировать случайные точки (Vec3), конфигурации камеры, размеры экрана
    - Сравнивать результаты `Camera.Project()` и `render.ProjectOld()` с epsilon = 1e-10
    - Минимум 100 итераций теста
    - _Requirements: 3.4, 8.3_

- [x] 3. Рефакторинг структуры Player
  - [x] 3.1 Добавить поле physics в структуру Player
    - Добавить поле `physics config.PhysicsConfig` в структуру `Player` в `internal/entity/player.go`
    - Добавить поле `Physics config.PhysicsConfig` в структуру `PlayerConfig`
    - Обновить конструктор `NewPlayer` для сохранения `cfg.Physics` в поле `physics`
    - _Requirements: 2.2, 2.4, 5.1_

  - [x] 3.2 Обновить метод Player.Update для использования сохранённой конфигурации
    - В методе `Update` заменить `config.Physics.Gravity` на `p.physics.Gravity`
    - Убедиться, что логика прыжка использует сохранённое значение гравитации
    - _Requirements: 5.2, 5.3_

  - [x] 3.3 Обновить метод Player.Draw для использования Camera.Project()
    - Заменить вызовы `render.Project(point, cam)` на `cam.Project(point)`
    - Обновить вызовы для точек `bottom` и `top`
    - _Requirements: 3.5, 3.6_

  - [ ]* 3.4 Написать property-based тест для независимости гравитации
    - **Property 3: Gravity configuration independence**
    - **Validates: Requirements 5.2, 8.2**
    - Генерировать случайные значения Gravity (0.1 - 1.0)
    - Создавать Player с различными значениями гравитации
    - Симулировать прыжок с фиксированной начальной скоростью
    - Проверять, что траектория соответствует формуле `velocity -= gravity`
    - Проверять edge case: gravity = 0 (бесконечный прыжок)
    - _Requirements: 5.2, 8.2_

- [ ] 4. Checkpoint - Проверка компонентов Camera и Player
  - Убедиться, что все тесты проходят
  - Проверить, что Camera и Player корректно используют переданные конфигурации
  - Спросить пользователя, если возникли вопросы

- [x] 5. Рефакторинг структуры Manager
  - [x] 5.1 Добавить поле gameConfig в структуру Manager
    - Добавить поле `gameConfig config.GameConfig` в структуру `Manager` в `internal/state/manager.go`
    - Обновить конструктор `NewManager` для приёма параметра `gameCfg config.GameConfig`
    - Сохранить переданную конфигурацию в поле `gameConfig`
    - _Requirements: 2.2, 2.4_

  - [x] 5.2 Добавить метод Manager.GameConfig()
    - Создать метод `func (m *Manager) GameConfig() config.GameConfig`
    - Метод должен возвращать сохранённую конфигурацию игры
    - Это позволит состояниям получать доступ к конфигурации через manager
    - _Requirements: 2.2, 7.1_

- [x] 6. Рефакторинг структуры GameState
  - [x] 6.1 Обновить сигнатуру конструктора NewGameState
    - Изменить сигнатуру на `NewGameState(manager *Manager, gameCfg config.GameConfig, cameraCfg config.CameraConfig, physicsCfg config.PhysicsConfig)`
    - Добавить поле `gameConfig config.GameConfig` в структуру `GameState`
    - Сохранить переданную конфигурацию в поле `gameConfig`
    - _Requirements: 2.3, 2.4, 4.1_

  - [x] 6.2 Обновить создание компонентов в NewGameState
    - Передать `gameCfg.ScreenWidth` в `NewSkyLayer` вместо `config.Game.ScreenWidth`
    - Создать `Camera` с параметрами `(float64(gameCfg.ScreenWidth), float64(gameCfg.ScreenHeight), cameraCfg)`
    - Создать `Player` с `PlayerConfig{Physics: physicsCfg}`
    - Создать `BalanceBarLayer` с `float64(gameCfg.ScreenWidth)` и `float64(gameCfg.ScreenHeight)`
    - _Requirements: 2.3, 2.5, 6.3_

  - [x] 6.3 Обновить метод GameState.Update для использования сохранённой конфигурации
    - Заменить `config.Game.DriftThreshold` на `g.gameConfig.DriftThreshold`
    - Обновить вызов `NewGameOverState` для передачи `g.gameConfig`
    - _Requirements: 4.2, 4.3, 4.4_

  - [x] 6.4 Обновить метод GameState.Draw для использования сохранённой конфигурации
    - Заменить `config.Game.DriftThreshold` на `g.gameConfig.DriftThreshold` в отображении текста
    - _Requirements: 4.2, 4.4_

  - [ ]* 6.5 Написать property-based тест для независимости drift threshold
    - **Property 2: Drift threshold configuration independence**
    - **Validates: Requirements 4.3**
    - Генерировать случайные значения DriftThreshold (50 - 500)
    - Создавать GameState с различными порогами
    - Симулировать достижение порога
    - Проверять, что drift mechanics активируется при настроенном пороге
    - _Requirements: 4.3_

- [x] 7. Рефакторинг структуры GameOverState
  - [x] 7.1 Обновить сигнатуру конструктора NewGameOverState
    - Изменить сигнатуру на `NewGameOverState(manager *Manager, score float64, gameCfg config.GameConfig)`
    - Добавить поле `gameConfig config.GameConfig` в структуру `GameOverState`
    - Сохранить переданную конфигурацию в поле `gameConfig`
    - _Requirements: 2.2, 6.1_

  - [x] 7.2 Обновить метод GameOverState.Update
    - При создании нового `GameState` передать все три конфигурации
    - Использовать `config.DefaultCameraConfig()` и `config.DefaultPhysicsConfig()` для новой игры
    - Передать сохранённый `g.gameConfig` в `NewGameState`
    - _Requirements: 6.1, 7.1_

  - [x] 7.3 Обновить метод GameOverState.Draw
    - Заменить `config.Game.ScreenWidth` на `g.gameConfig.ScreenWidth`
    - Заменить `config.Game.ScreenHeight` на `g.gameConfig.ScreenHeight`
    - _Requirements: 6.2, 6.4_

- [ ] 8. Checkpoint - Проверка состояний игры
  - Убедиться, что все тесты проходят
  - Проверить, что GameState и GameOverState корректно используют конфигурации
  - Спросить пользователя, если возникли вопросы

- [x] 9. Обновление других состояний и компонентов
  - [x] 9.1 Обновить MenuState для передачи конфигураций
    - Если `MenuState` создаёт `GameState`, обновить вызов для передачи всех конфигураций
    - Получить `GameConfig` из `manager.GameConfig()`
    - Создать `CameraConfig` и `PhysicsConfig` через Default функции
    - _Requirements: 2.3, 7.1_

  - [x] 9.2 Найти и обновить все вызовы render.Project в world слоях
    - Использовать grep для поиска всех вызовов `render.Project(`
    - Заменить `render.Project(point, cam)` на `cam.Project(point)` во всех найденных местах
    - Проверить файлы в `internal/world/` на использование функции проекции
    - _Requirements: 3.5, 3.6, 6.3_

- [x] 10. Обновление main.go
  - [x] 10.1 Создать конфигурации в main.go
    - Создать `gameCfg := config.DefaultGameConfig()`
    - Создать `cameraCfg := config.DefaultCameraConfig()`
    - Создать `physicsCfg := config.DefaultPhysicsConfig()`
    - _Requirements: 2.5, 7.1, 7.2_

  - [x] 10.2 Обновить создание Manager и начального состояния
    - Передать `gameCfg` в `NewManager(nil, gameCfg)`
    - Обновить создание начального состояния (MenuState) если требуется
    - _Requirements: 2.5, 7.1_

  - [x] 10.3 Обновить метод Layout для использования конфигурации из Manager
    - Изменить `Layout` для получения размеров через `g.manager.GameConfig()`
    - Вернуть `cfg.ScreenWidth, cfg.ScreenHeight`
    - _Requirements: 7.2, 7.3_

- [x] 11. Удаление глобальных переменных из config пакета
  - [x] 11.1 Удалить глобальные переменные из internal/config/constants.go
    - Удалить объявление `var Camera = DefaultCameraConfig()`
    - Удалить объявление `var Game = DefaultGameConfig()`
    - Удалить объявление `var Physics = DefaultPhysicsConfig()`
    - Оставить типы и Default функции без изменений
    - _Requirements: 1.1, 1.2, 1.3_

  - [x] 11.2 Удалить файл internal/render/projection.go
    - Удалить файл `internal/render/projection.go` полностью
    - Функция `Project` теперь является методом `Camera`
    - _Requirements: 3.5_

- [ ] 12. Финальная проверка и валидация
  - [x] 12.1 Проверить отсутствие ссылок на глобальные конфигурации
    - Выполнить поиск по всему проекту: `config\.(Camera|Game|Physics)\.`
    - Убедиться, что поиск не находит совпадений (кроме комментариев)
    - _Requirements: 1.4, 6.4_

  - [ ]* 12.2 Запустить все unit тесты
    - Выполнить `go test ./...` для запуска всех тестов
    - Убедиться, что все тесты проходят успешно
    - _Requirements: 8.1, 8.2, 8.3, 8.4_

  - [ ]* 12.3 Запустить все property-based тесты
    - Выполнить property-based тесты для проекции, гравитации и drift threshold
    - Убедиться, что все свойства выполняются на 100+ итерациях
    - _Requirements: 8.1, 8.2, 8.3, 8.4_

- [ ] 13. Checkpoint - Финальная проверка
  - Убедиться, что приложение компилируется без ошибок
  - Убедиться, что все тесты проходят
  - Убедиться, что визуальный вывод и поведение игры идентичны оригиналу
  - Спросить пользователя, если возникли вопросы

## Notes

- Задачи, помеченные `*`, являются опциональными и могут быть пропущены для более быстрой реализации
- Каждая задача ссылается на конкретные требования для отслеживаемости
- Checkpoint задачи обеспечивают инкрементальную валидацию
- Property-based тесты валидируют универсальные свойства корректности
- Unit тесты валидируют конкретные примеры и граничные случаи
- Рефакторинг выполняется инкрементально: сначала добавляются новые параметры, затем обновляется использование, и только в конце удаляются глобальные переменные
