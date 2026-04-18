# План Реализации: Рефакторинг Архитектуры "The Fiasko"

## Обзор

Данный план описывает прагматичный рефакторинг псевдо-3D раннера "The Fiasko" (Go + Ebiten). Рефакторинг выполняется в 4 этапа с фокусом на простоту и управляемость кода. Каждый этап сохраняет работоспособность игры.

**Принцип:** Простота важнее идеальной архитектуры. Никакого overengineering.

**7 КРИТИЧНЫХ АРХИТЕКТУРНЫХ ИСПРАВЛЕНИЙ:**

1. **FIX #1 - Renderer removed from World, owned by Game:**
   - World НЕ имеет поля renderer
   - Game владеет ОБОИМИ: World (logic) И Renderer (presentation)
   - World НЕ имеет метода Draw()
   - Game.Draw() вызывает renderer.DrawWorld(screen, world.player, world.obstacles) напрямую
   - NewWorld(config, pool) - НЕТ параметра renderer

2. **FIX #2 - Input handling moved to World level:**
   - НЕТ HandlePlayerInput в player_logic.go
   - handleInput(w, dt) функция в world.go
   - Input обрабатывается на уровне World/Game, НЕ в entity logic

3. **FIX #3 - Obstacle.Type changed to enum:**
   - Создать ObstacleType enum (ObstacleTypeLog, ObstacleTypeRock)
   - Type поле имеет тип ObstacleType, НЕ string
   - Spawner использует ObstacleTypeLog при инициализации

4. **FIX #4 - Score tracking added:**
   - World имеет поле score (float64)
   - World.Update() увеличивает score: w.score += dt
   - Game.Draw() отображает score overlay

5. **FIX #5 - dt clamp improved:**
   - Использовать `dt = math.Min(dt, 0.05)` вместо `if dt > 0.1 { dt = 0.1 }`

6. **FIX #6 - DespawnZ added to Config:**
   - Config имеет поля: SpawnRangeX, SpawnZ, DespawnZ
   - Spawner использует config.SpawnRangeX и config.SpawnZ
   - removeOffscreenObstacles использует config.DespawnZ вместо магического -10

7. **FIX #7 - sortedObstacles capacity fixed:**
   - Renderer предаллоцирует sortedObstacles с capacity = config.MaxObstacles
   - Renderer переиспользует буфер: sortedObstacles = sortedObstacles[:0]

**КРИТИЧНЫЕ АРХИТЕКТУРНЫЕ ИЗМЕНЕНИЯ:**
- **Data vs Logic Separation**: Player/Obstacle = только данные, логика в отдельных файлах (player_logic.go, obstacle_logic.go)
- **World = Orchestrator**: Логика разбита на отдельные функции (updatePlayer, updateObstacles, handleSpawning, handleCollisions, removeOffscreenObstacles, handleInput)
- **Game владеет World И Renderer**: World = game logic, Renderer = presentation layer, World НЕ знает о Renderer
- **Input на уровне World**: handleInput() в world.go, НЕ HandlePlayerInput в player_logic.go (input = внешний источник, НЕ entity logic)
- **Obstacle.Type = enum**: ObstacleType (ObstacleTypeLog, ObstacleTypeRock), НЕ string
- **Score tracking**: World.score отслеживает прогресс (w.score += dt)
- **Spawner Fixes**: timer -= interval (НЕ timer = 0), возвращает []*Obstacle, использует Config.SpawnRangeX/SpawnZ, Type = enum
- **ObjectPool**: maxSize ограничение, НЕТ active флага, полный Reset()
- **Collision**: Early-out оптимизация (if math.Abs(obs.Position.Z - playerZ) > 5.0 { continue })
- **Renderer**: Переиспользует буфер sortedObstacles (sortedObstacles[:0]), capacity = config.MaxObstacles
- **Config**: SpawnRangeX, SpawnZ, DespawnZ (НЕ хардкод)
- **dt clamp**: math.Min(dt, 0.05) вместо if dt > 0.1 { dt = 0.1 }
- **Testing**: Manual testing (основной метод), простые unit tests (опционально), НЕТ property-based тестов

## Задачи

### ЭТАП 1: СТАБИЛИЗАЦИЯ (3-4 дня)

- [x] 1. Создать базовые структуры данных
  - [x] 1.1 Создать Vec3 с математическими операциями
    - Создать `internal/core/vector.go` с типом Vec3
    - Реализовать методы: Add, Sub, Scale, Length, Normalize
    - _Требования: 3, 4.3_
  
  - [x] 1.2 Написать unit-тесты для Vec3
    - Тестировать все математические операции
    - Тестировать граничные случаи (нулевой вектор, нормализация)

- [x] 2. Создать Config с значениями по умолчанию
  - [x] 2.1 Создать структуру GameConfig
    - Создать `internal/config/config.go` с типом Config
    - Добавить поля: ScreenWidth, ScreenHeight, TargetFPS, Gravity, JumpVelocity, BalanceSpeed, MaxBalance, SpawnInterval, ObstacleSpeed, MaxObstacles, CameraFocalLength, CameraHorizonY
    - **КРИТИЧНО (FIX #6)**: Добавить поля SpawnRangeX, SpawnZ, DespawnZ (НЕ хардкод в Spawner и World)
    - Реализовать функцию DefaultConfig() с разумными значениями по умолчанию
    - _Требования: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6, 7.7_
  
  - [x] 2.2 Написать unit-тесты для Config
    - Тестировать, что DefaultConfig() возвращает валидные значения

- [x] 3. Создать простые структуры Player и Obstacle (ТОЛЬКО данные)
  - [x] 3.1 Создать структуру Player (data-only)
    - Создать `internal/entity/player.go` с типом Player
    - **КРИТИЧНО**: Добавить ТОЛЬКО поля данных: Position (Vec3), Velocity (Vec3), Balance (float64), Width, Height, Depth (float64)
    - **НЕТ методов** кроме конструктора NewPlayer()
    - **НЕТ зависимостей** от Config, ResourceManager, Camera в структуре
    - Реализовать конструктор NewPlayer() с дефолтными значениями
    - _Требования: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_
  
  - [x] 3.2 Создать структуру Obstacle (data-only)
    - Создать `internal/entity/obstacle.go` с типом Obstacle
    - **КРИТИЧНО (FIX #3)**: Создать ObstacleType enum (ObstacleTypeLog, ObstacleTypeRock, etc.)
    - **КРИТИЧНО (FIX #3)**: Добавить ТОЛЬКО поля данных: Position (Vec3), Velocity (Vec3), Width, Height, Depth (float64), Type (ObstacleType - enum, НЕ string)
    - **НЕТ active флага** (состояние контролируется присутствием в world.obstacles)
    - Реализовать метод Reset() для полного сброса всех полей
    - **НЕТ других методов**
    - _Требования: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_

- [x] 4. Создать World как Orchestrator (НЕ God Object)
  - [x] 4.1 Создать структуру World
    - Создать `internal/game/world.go` с типом World
    - Добавить поля: player (*Player), obstacles ([]*Obstacle), spawner (*Spawner), config (*Config), pool (*ObstaclePool), state (GameState), score (float64)
    - **КРИТИЧНО (FIX #1)**: НЕТ поля renderer - World НЕ владеет Renderer! Game владеет обоими: World И Renderer
    - **КРИТИЧНО**: Добавить GameState enum (StatePlaying, StateGameOver)
    - **КРИТИЧНО (FIX #4)**: Добавить поле score для отслеживания прогресса
    - Реализовать конструктор NewWorld(config *Config, pool *ObstaclePool) - НЕТ параметра renderer!
    - _Требования: 4.1, 4.2, 4.3, 4.4, 4.9_
  
  - [x] 4.2 Реализовать World.Update(dt float64) - базовая версия с делегированием
    - **КРИТИЧНО**: World ОРКЕСТРИРУЕТ, не выполняет всю логику сам
    - **КРИТИЧНО (FIX #2)**: Вызывать handleInput(w, dt) для обработки ввода (НЕ HandlePlayerInput из entity_logic)
    - Вызывать updatePlayer(w, dt) для физики игрока
    - Вызывать updateObstacles(w, dt) для обновления препятствий
    - **КРИТИЧНО (FIX #4)**: Увеличивать score: w.score += dt
    - Пока НЕ добавлять spawner, коллизии, удаление - это будет в Этапе 2
    - _Требования: 4.3, 4.5_
  
  - [x] 4.3 World НЕ имеет метода Draw()
    - **КРИТИЧНО (FIX #1)**: World НЕ рисует сам, НЕ делегирует Renderer
    - World предоставляет только данные (player, obstacles) для рендеринга
    - Рендеринг выполняется через Game, который владеет Renderer
    - _Требования: 4.4_

- [x] 5. Создать отдельные файлы логики (Data vs Logic Separation)
  - [x] 5.1 Создать player_logic.go с функциями логики игрока
    - Создать `internal/entity/player_logic.go`
    - **КРИТИЧНО**: Реализовать UpdatePlayer(p *Player, config *Config, dt float64)
      - Применение гравитации (если Position.Y > 0)
      - Обновление позиции на основе скорости и dt
      - Ограничение позиции (Position.Y не может быть < 0)
    - **КРИТИЧНО (FIX #2)**: НЕ реализовывать HandlePlayerInput - input обрабатывается на уровне World через handleInput()
    - **КРИТИЧНО**: Реализовать CheckPlayerFall(p *Player, config *Config) bool
      - Проверка условия падения (math.Abs(p.Balance) > config.MaxBalance)
    - _Требования: 1.1, 1.2, 1.4, 1.6, 1.7, 1.8_
  
  - [x] 5.2 Создать obstacle_logic.go с функциями логики препятствий
    - Создать `internal/entity/obstacle_logic.go`
    - **КРИТИЧНО**: Реализовать UpdateObstacle(o *Obstacle, dt float64)
      - Обновление позиции на основе скорости и dt
      - Препятствия движутся к игроку (o.Position.Z -= o.Velocity.Z * dt)
    - _Требования: 1.1, 1.2, 1.4, 1.5, 1.6, 1.7, 3.5_

- [x] 6. Обновить Game для владения World И Renderer
  - [x] 6.1 Упростить структуру Game и добавить GameStateType
    - Модифицировать Game структуру: добавить поля world (*World) и renderer (*Renderer)
    - **КРИТИЧНО (FIX #1)**: Game владеет ОБОИМИ: World (logic) И Renderer (presentation)
    - **КРИТИЧНО**: Добавить поле state (GameStateType) с enum (StateMenu, StatePlaying, StateGameOver, StatePaused)
    - Удалить прямые ссылки на player, obstacles, ui
    - Добавить поле lastTime (time.Time) для вычисления dt
    - _Требования: 2.1, 2.2, 2.3_
  
  - [x] 6.2 Реализовать Game.Update() с вычислением dt и управлением состояниями
    - Вычислять dt как время между кадрами в секундах
    - **КРИТИЧНО (FIX #5)**: Ограничивать dt: `dt = math.Min(dt, 0.05)` (НЕ `if dt > 0.1 { dt = 0.1 }`)
    - **КРИТИЧНО**: Использовать switch по state для управления логикой
    - В StatePlaying: вызывать world.Update(dt), проверять world.IsGameOver()
    - **КРИТИЧНО**: Использовать state (НЕ error) для game over
    - _Требования: 1.6, 1.7, 1.8, 2.4_
  
  - [x] 6.3 Реализовать Game.Draw(screen *ebiten.Image) с управлением состояниями
    - Использовать switch по state для управления отрисовкой
    - **КРИТИЧНО (FIX #1)**: Вызывать renderer.DrawWorld(screen, world.player, world.obstacles) напрямую
    - **КРИТИЧНО (FIX #1)**: Game владеет renderer и вызывает его, НЕ world.Draw()
    - **КРИТИЧНО (FIX #4)**: Отрисовывать score overlay используя world.score
    - _Требования: 2.4_

- [x] 7. Checkpoint - Базовая стабилизация
  - Запустить игру и проверить, что она работает
  - Убедиться, что Update содержит только логику
  - Убедиться, что Draw не изменяет состояние
  - Убедиться, что все Update методы используют dt

### ЭТАП 2: ГЕЙМПЛЕЙ (4-5 дней)

- [ ] 9. Создать простой ResourceManager
  - [x] 9.1 Реализовать ResourceManager с кэшированием
    - Создать `internal/resources/manager.go` с типом ResourceManager
    - Добавить поле images (map[string]*ebiten.Image)
    - Реализовать метод LoadImage(path string) (*ebiten.Image, error)
    - Проверять кэш перед загрузкой
    - Возвращать ошибку при неудачной загрузке (НЕ fallback)
    - Логировать загрузку каждой текстуры
    - _Требования: 8.1, 8.2, 8.3, 8.4, 8.5, 8.6, 8.7, 8.8_
  
  - [~] 9.2 Написать простые unit-тесты для ResourceManager
    - Тестировать кэширование (LoadImage дважды возвращает тот же указатель)
    - Тестировать обработку ошибок загрузки

- [ ] 10. Создать простой ObjectPool для препятствий с критичными исправлениями
  - [x] 10.1 Реализовать ObstaclePool с ограничением размера
    - Создать `internal/pools/obstacle_pool.go` с типом ObstaclePool
    - **КРИТИЧНО**: Добавить поля: obstacles ([]*Obstacle), maxSize (int)
    - Реализовать метод Get() *Obstacle (возвращает из пула или создаёт новое)
    - **КРИТИЧНО**: Реализовать метод Put(obs *Obstacle)
      - Вызывать obs.Reset() для полного сброса состояния
      - Возвращать в пул ТОЛЬКО если len(obstacles) < maxSize
      - Иначе отбрасывать (GC соберёт)
    - НЕ использовать generics, только []*Obstacle
    - _Требования: 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.8_
  
  - [x] 10.2 Реализовать полный сброс состояния в Obstacle.Reset()
    - **КРИТИЧНО**: Сбрасывать ВСЕ поля: Position, Velocity, Type в нулевые значения
    - **НЕТ active флага** (состояние контролируется присутствием в world.obstacles)
    - _Требования: 6.6_
  
  - [~] 10.3 Написать простые unit-тесты для ObstaclePool
    - Тестировать Get() из пустого пула (создаёт новое)
    - Тестировать Get() из непустого пула (возвращает существующее)
    - Тестировать Put() сбрасывает состояние
    - Тестировать maxSize ограничение

- [x] 11. Создать Spawner для генерации препятствий с критичными исправлениями
  - [x] 11.1 Реализовать Spawner с таймером
    - Создать `internal/game/spawner.go` с типом Spawner
    - Добавить поля: timer (float64), interval (float64), pool (*ObstaclePool), config (*Config)
    - Реализовать конструктор NewSpawner(config *Config, pool *ObstaclePool)
    - _Требования: 5.1, 5.2, 5.5, 5.6_
  
  - [x] 11.2 Реализовать Spawner.Update(dt float64) []*Obstacle с критичными исправлениями
    - Увеличивать timer на dt
    - **КРИТИЧНО**: Использовать цикл for и timer -= interval (НЕ timer = 0!)
      - Это сохраняет точность таймера
      - Поддерживает burst spawn при больших dt
    - **КРИТИЧНО**: Возвращать []*Obstacle (НЕ *Obstacle)
      - Позволяет вернуть несколько препятствий при лагах
    - Получать препятствие из pool.Get()
    - **КРИТИЧНО (FIX #6)**: Инициализировать позицию из Config (НЕ хардкод)
      - X: (rand.Float64() - 0.5) * config.SpawnRangeX
      - Y: 0
      - Z: config.SpawnZ
    - Инициализировать скорость из config.ObstacleSpeed
    - **КРИТИЧНО (FIX #3)**: Устанавливать Type = ObstacleTypeLog (enum, НЕ string "log")
    - _Требования: 5.3, 5.4, 5.5, 5.6_
  
  - [x] 11.3 Написать простые unit-тесты для Spawner
    - Тестировать, что Update() возвращает пустой slice до истечения интервала
    - Тестировать, что Update() возвращает препятствия после интервала
    - **КРИТИЧНО**: Тестировать burst spawn (dt > interval * 2 должен вернуть 2+ препятствия)

- [x] 12. Интегрировать Spawner в World с обработкой burst spawn
  - [x] 12.1 Добавить Spawner в World.Update(dt) через handleSpawning()
    - **КРИТИЧНО**: Создать отдельную функцию handleSpawning(w *World, dt float64)
    - Вызывать spawner.Update(dt) (возвращает []*Obstacle)
    - **КРИТИЧНО**: Итерировать по всем возвращённым препятствиям
    - Для каждого препятствия: проверять лимит MaxObstacles
    - Если лимит не достигнут, добавлять в obstacles
    - Если лимит достигнут, возвращать препятствие в pool.Put()
    - _Требования: 4.5, 4.11_

- [x] 13. Реализовать удаление препятствий за границами экрана
  - [x] 13.1 Реализовать removeOffscreenObstacles() как отдельную функцию
    - **КРИТИЧНО**: Создать функцию removeOffscreenObstacles(w *World)
    - **КРИТИЧНО (FIX #6)**: Проверять каждое препятствие: если Position.Z < config.DespawnZ, удалять из списка
    - Возвращать удалённые препятствия в pool.Put()
    - Использовать эффективный алгоритм удаления (in-place filtering)
    - _Требования: 4.6_
  
  - [x] 13.2 Вызывать removeOffscreenObstacles() в World.Update(dt)
    - Вызывать после обновления всех препятствий
    - _Требования: 4.6_

- [x] 14. Реализовать простую проверку коллизий с early-out оптимизацией
  - [x] 14.1 Создать функцию CheckAABBCollision в collision.go
    - Создать `internal/game/collision.go`
    - **КРИТИЧНО**: Реализовать CheckAABBCollision(player *Player, obstacle *Obstacle) bool как отдельную функцию
    - Использовать простую AABB коллизию
    - Проверять пересечение по X, Y, Z осям
    - _Требования: 4.7, 4.10_
  
  - [x] 14.2 Добавить проверку коллизий в World.Update(dt) через handleCollisions()
    - **КРИТИЧНО**: Создать отдельную функцию handleCollisions(w *World) bool
    - **КРИТИЧНО**: Добавить early-out оптимизацию:
      - if math.Abs(obs.Position.Z - playerZ) > 5.0 { continue }
      - Пропускает далёкие объекты, значительно ускоряет проверку
    - Для каждого близкого препятствия вызывать CheckAABBCollision(player, obstacle)
    - Если коллизия обнаружена, возвращать true
    - **КРИТИЧНО**: Использовать GameState (НЕ error) для game over
    - _Требования: 4.7_
  
  - [x] 14.3 Написать простые unit-тесты для CheckAABBCollision
    - Тестировать обнаружение пересечения AABB
    - Тестировать отсутствие пересечения

- [x] 15. Вынести все магические числа в Config
  - [x] 15.1 Заменить хардкод на значения из Config
    - Найти все магические числа в Player, Obstacle, World, Spawner
    - Заменить на config.FieldName
    - Убедиться, что Config передаётся через конструкторы (НЕ глобальная переменная)
    - _Требования: 7.5, 7.6, 7.7_

- [ ] 16. Checkpoint - Геймплей работает
  - Запустить игру и проверить спавн препятствий
  - Проверить, что препятствия переиспользуются (пул работает)
  - Проверить, что коллизии обнаруживаются
  - Проверить, что лимит препятствий соблюдается
  - Проверить, что препятствия удаляются за границами экрана

### ЭТАП 3: ВИЗУАЛ И UX (3-4 дня)

- [ ] 17. Создать Renderer отдельно от Entity с критичными оптимизациями
  - [~] 17.1 Создать Renderer структуру с переиспользуемым буфером
    - Создать `internal/render/renderer.go` с типом Renderer
    - **КРИТИЧНО**: Добавить поля:
      - camera (*Camera)
      - resMgr (*ResourceManager)
      - playerTexture (*ebiten.Image)
      - obstacleTexture (*ebiten.Image)
      - sortedObstacles ([]*Obstacle) - переиспользуемый буфер
    - Реализовать конструктор NewRenderer(config *Config, resMgr *ResourceManager)
    - Preload критических текстур в конструкторе
    - **КРИТИЧНО (FIX #7)**: Предаллоцировать sortedObstacles буфер с capacity = config.MaxObstacles
  
  - [~] 17.2 Реализовать Renderer.DrawWorld() с переиспользованием буфера
    - **КРИТИЧНО (FIX #1)**: Реализовать DrawWorld(screen *ebiten.Image, player *Player, obstacles []*Obstacle)
    - **КРИТИЧНО (FIX #7)**: Переиспользовать буфер: sortedObstacles = sortedObstacles[:0]
      - НЕ создавать новый slice каждый кадр
      - Сбрасывать длину, сохранять capacity
    - Копировать obstacles в sortedObstacles
    - Сортировать по Z (painter's algorithm)
    - Рисовать препятствия (дальние первыми)
    - Рисовать игрока
  
  - [~] 17.3 Реализовать Renderer.DrawPlayer() и Renderer.DrawObstacle()
    - **КРИТИЧНО (FIX #1)**: Entity НЕ имеют методов Draw()
    - Реализовать DrawPlayer(screen *ebiten.Image, p *Player)
      - Проецировать через camera.Project()
      - Создавать вершины для DrawTriangles
      - Отрисовывать playerTexture
    - Реализовать DrawObstacle(screen *ebiten.Image, o *Obstacle)
      - Проецировать через camera.Project()
      - Создавать вершины для DrawTriangles
      - Отрисовывать obstacleTexture

- [ ] 18. Обновить Camera для использования Config
  - [~] 18.1 Модифицировать Camera структуру
    - Обновить `internal/render/camera.go`
    - Добавить поле config (*Config)
    - Использовать config.CameraFocalLength и config.CameraHorizonY
    - Сохранить существующую логику проекции Project(point Vec3) (float64, float64, float64)
  
  - [~] 18.2 Интегрировать Camera в Renderer
    - Создавать Camera в NewRenderer через NewCamera(config)
    - Использовать camera в DrawPlayer и DrawObstacle

- [ ] 19. Создать простой UI
  - [~] 19.1 Создать MenuUI структуру
    - Создать `internal/ui/menu.go` с типами MenuUI и Button
    - Добавить поля: buttons ([]Button)
    - Реализовать методы Update() и Draw(screen *ebiten.Image)
    - НЕ использовать сложные UI фреймворки
    - _Требования: 9.1, 9.2, 9.3, 9.4, 9.5_
  
  - [~] 19.2 Реализовать простую отрисовку кнопок
    - Использовать ebitenutil.DrawRect для кнопок
    - Использовать ebitenutil.DebugPrintAt для текста
    - _Требования: 9.3_
  
  - [~] 19.3 Написать unit-тесты для MenuUI
    - Тестировать обработку кликов по кнопкам

- [ ] 20. Добавить базовые звуковые эффекты
  - [~] 20.1 Создать SoundManager структуру
    - Создать `internal/audio/sound.go` с типом SoundManager
    - Добавить поля: sounds (map[string]*audio.Player), resMgr (*ResourceManager)
    - Реализовать методы LoadSound(name, path string) error и Play(name string)
    - НЕ использовать сложные аудио-системы
    - _Требования: 10.1, 10.2, 10.3, 10.4, 10.5_
  
  - [~] 20.2 Интегрировать звуки в игру
    - **КРИТИЧНО (FIX #2)**: Добавить звук прыжка в handleInput() (в world.go) при нажатии W
    - Добавить звук столкновения в handleCollisions() при коллизии
    - Добавить фоновую музыку в Game.Update()
    - _Требования: 10.2_
  
  - [~] 20.3 Написать unit-тесты для SoundManager
    - Тестировать загрузку звуков
    - Тестировать воспроизведение несуществующего звука (не должно паниковать)

- [~] 21. Checkpoint - Визуал и UX готовы
  - Запустить игру и проверить отрисовку
  - Проверить, что текстуры загружаются один раз (проверить логи)
  - **КРИТИЧНО (FIX #7)**: Проверить, что sortedObstacles буфер переиспользуется (нет allocation каждый кадр)
  - **КРИТИЧНО (FIX #1)**: Проверить, что Game владеет Renderer, World не знает о нём
  - **КРИТИЧНО (FIX #4)**: Проверить, что score отображается
  - Проверить, что UI отображается
  - Проверить, что звуки воспроизводятся

### ЭТАП 4: ЧИСТКА (2 дня)

- [x] 22. Удалить мусорный код
- [ ] 22. Удалить мусорный код
  - [~] 22.1 Удалить пустые директории UI
    - Удалить `internal/ui/animations`
    - Удалить `internal/ui/components`
    - Удалить `internal/ui/core`
    - Удалить `internal/ui/styles`
    - Удалить `internal/ui/widgets`
    - _Требования: 11.1_
  
  - [~] 22.2 Удалить весь закомментированный код
    - Найти и удалить все закомментированные блоки кода
    - _Требования: 11.2_
  
  - [~] 22.3 Удалить неиспользуемые геттеры и сеттеры
    - Найти избыточные геттеры (Go идиома - прямой доступ к полям)
    - Заменить геттеры прямым доступом
    - Удалить неиспользуемые методы
    - _Требования: 11.3, 11.4_
  
  - [~] 22.4 Проверить, что все файлы используются
    - Убедиться, что каждый оставшийся файл активно используется
    - _Требования: 11.5_

- [ ] 23. Добавить простое логирование
  - [~] 23.1 Добавить логирование ключевых событий
    - Использовать стандартный log.Println()
    - Логировать: "Player jumped", "Obstacle spawned", "Collision detected"
    - Добавить префиксы [INFO], [WARNING], [ERROR]
    - _Требования: 12.1, 12.2, 12.3, 12.4_
  
  - [~] 23.2 Добавить флаг DEBUG для отключения логов
    - Создать константу DEBUG bool в main.go
    - Обернуть debug логи в if DEBUG { ... }
    - _Требования: 12.5_

- [ ] 24. Финальная проверка
  - [~] 24.1 Запустить игру и проверить все механики
    - Проверить балансирование (клавиши A/D)
    - Проверить прыжки (клавиша W)
    - Проверить визуальный стиль (псевдо-3D, текстуры)
    - Проверить подсчёт очков
    - Проверить условия game over (падение, столкновение)
    - Проверить управление
    - _Требования: 13.1, 13.2, 13.3, 13.4, 13.5, 13.6_
  
  - [~] 24.2 Запустить все тесты (если есть)
    - Выполнить `go test ./...`
    - Убедиться, что все unit тесты проходят
  
  - [~] 24.3 Проверить производительность
    - Убедиться, что FPS стабильный
    - Проверить, что текстуры загружаются один раз (проверить логи)
    - Проверить, что препятствия переиспользуются (пул работает)
    - **КРИТИЧНО**: Проверить, что sortedObstacles буфер переиспользуется

- [~] 25. Checkpoint - Рефакторинг завершён
  - Убедиться, что игра работает идентично до рефакторинга
  - Убедиться, что код чистый и управляемый
  - Убедиться, что все требования выполнены
  - **КРИТИЧНО**: Убедиться, что все 7 архитектурных исправлений соблюдены:
    - **FIX #1**: Game владеет World И Renderer (World НЕ знает о Renderer, НЕТ World.Draw())
    - **FIX #2**: Input handling в world.go (handleInput), НЕ HandlePlayerInput в player_logic.go
    - **FIX #3**: Obstacle.Type = enum (ObstacleType), НЕ string
    - **FIX #4**: Score tracking работает (world.score += dt, отображается в Game.Draw())
    - **FIX #5**: dt clamp использует math.Min(dt, 0.05) вместо if dt > 0.1
    - **FIX #6**: Config содержит SpawnRangeX, SpawnZ, DespawnZ (используются в Spawner и World)
    - **FIX #7**: Renderer переиспользует sortedObstacles буфер (capacity = config.MaxObstacles)
  - **Дополнительные принципы**:
    - Data vs Logic Separation (player.go/player_logic.go, obstacle.go/obstacle_logic.go)
    - World = Orchestrator (отдельные функции для каждой ответственности)
    - Spawner использует timer -= interval и возвращает []*Obstacle
    - ObjectPool имеет maxSize, НЕТ active флага
    - Collision имеет early-out оптимизацию

## Примечания

- Задачи, помеченные `*`, являются опциональными и могут быть пропущены
- Каждая задача ссылается на конкретные требования для отслеживаемости
- Checkpoints обеспечивают инкрементальную валидацию на каждом этапе
- **КРИТИЧНО**: Тестирование в основном manual - запускаем игру и проверяем (это ГЛАВНЫЙ метод)
- Простые unit тесты только для критической логики (опционально)
- **НЕТ property-based тестов, НЕТ gopter, НЕТ сложных генераторов** - это overengineering для игры
- Рефакторинг выполняется инкрементально с сохранением работоспособности игры
- Фокус на простоте: никаких ECS, event bus, spatial hash, batch renderer

**Критичные архитектурные принципы:**
- **Data vs Logic Separation**: Entity = только данные, логика в отдельных файлах
- **World = Orchestrator**: Координирует, не выполняет всё сам (отдельные функции)
- **Game владеет World И Renderer**: Разделение logic и presentation, предотвращает God Object 2.0
- **World НЕ владеет Renderer**: World = game logic, Renderer = presentation layer
- **Input на уровне World/Game**: handleInput в world.go, НЕ в entity_logic
- **Obstacle.Type = enum**: ObstacleType (compile-time проверка), НЕ string
- **Score tracking**: World.score отслеживает прогресс
- **Spawner**: timer -= interval (НЕ timer = 0), возвращает []*Obstacle, Type = enum
- **ObjectPool**: maxSize ограничение, НЕТ active флага, полный Reset()
- **Collision**: Early-out оптимизация (if math.Abs(obs.Position.Z - playerZ) > 5.0)
- **Renderer**: Переиспользует буфер sortedObstacles (sortedObstacles[:0]), capacity = config.MaxObstacles
- **Config**: SpawnRangeX, SpawnZ, DespawnZ (НЕ хардкод)
- **dt clamp**: math.Min(dt, 0.05) вместо if dt > 0.1 { dt = 0.1 }
- **GameState**: Использовать state enum (НЕ error) для game over
