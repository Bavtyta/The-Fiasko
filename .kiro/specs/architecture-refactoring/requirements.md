# Документ Требований: Рефакторинг Архитектуры "The Fiasko"

## Введение

Данный документ описывает требования к прагматичному рефакторингу псевдо-3D раннера "The Fiasko" (Go + Ebiten). Игра работает, но код хаотичен и сложен в поддержке.

**Цель рефакторинга:** Сделать код управляемым, чистым и готовым к расширению, сохранив текущий функционал игры.

**Принцип:** Простота важнее идеальной архитектуры. Никакого overengineering.

## Глоссарий

- **World**: Контейнер для всех игровых объектов (игрок, препятствия)
- **Player**: Игровой персонаж с позицией, скоростью и управлением
- **Obstacle**: Препятствие на пути игрока
- **Spawner**: Система генерации препятствий
- **Object_Pool**: Простой пул для переиспользования препятствий
- **Resource_Manager**: Простой менеджер для загрузки текстур один раз
- **Config**: Конфигурация игры (скорости, частота спавна и т.д.)
- **Game_Loop**: Основной цикл: Update (логика) → Draw (отрисовка)

## Требования

### ЭТАП 1: СТАБИЛИЗАЦИЯ (Приоритет: КРИТИЧНО)

#### Требование 1: Разделение Update и Draw

**User Story:** Как разработчик, я хочу чёткое разделение логики и отрисовки, чтобы избежать багов и упростить отладку.

**Acceptance Criteria:**

1. THE Update() метод SHALL содержать ТОЛЬКО игровую логику (позиции, скорости, коллизии)
2. THE Draw() метод SHALL содержать ТОЛЬКО код отрисовки
3. THE Draw() метод SHALL NOT изменять состояние игры
4. THE Update() метод SHALL NOT работать с ebiten.Image напрямую
5. FOR ALL игровых объектов, логика и рендеринг SHALL быть разделены
6. FOR ALL систем, Update SHALL принимать параметр dt (delta time) типа float64
7. FOR ALL временных вычислений (скорость, таймеры), они SHALL использовать dt для независимости от FPS
8. THE Game SHALL вычислять dt как время между кадрами (в секундах) и передавать его в World.Update(dt)

#### Требование 2: Устранение God Object

**User Story:** Как разработчик, я хочу убрать God Object (Game), чтобы упростить структуру кода.

**Acceptance Criteria:**

1. THE Game структура SHALL содержать только ссылку на World
2. THE Game структура SHALL NOT содержать прямые ссылки на player, obstacles, ui
3. THE World SHALL управлять всеми игровыми объектами
4. THE World SHALL иметь методы Update() и Draw()

#### Требование 3: Выделение Player и Obstacle

**User Story:** Как разработчик, я хочу вынести Player и Obstacle в отдельные файлы, чтобы улучшить организацию кода.

**Acceptance Criteria:**

1. THE Player SHALL быть определён в `internal/entity/player.go`
2. THE Obstacle SHALL быть определён в `internal/entity/obstacle.go`
3. THE Player SHALL содержать: Position (Vec3), Velocity (Vec3), методы управления
4. THE Obstacle SHALL содержать: Position (Vec3), Type, методы коллизии
5. THE Player и Obstacle SHALL иметь методы Update(dt float64) и Draw(screen *ebiten.Image)
6. THE Player и Obstacle SHALL NOT содержать вложенные системы (physics, renderer)
7. THE Player и Obstacle SHALL быть простыми структурами с конкретными типами полей

#### Требование 4: Создание World

**User Story:** Как разработчик, я хочу создать World для управления всеми игровыми объектами.

**Acceptance Criteria:**

1. THE World SHALL быть определён в `internal/game/world.go`
2. THE World SHALL содержать: player *Player, obstacles []*Obstacle, spawner *Spawner, config *Config
3. THE World SHALL иметь метод Update(dt float64) для обновления всех объектов
4. THE World SHALL иметь метод Draw(screen *ebiten.Image) для отрисовки всех объектов
5. THE World SHALL выполнять следующие обязанности в Update():
   - Обновлять Player (вызывать player.Update(dt))
   - Обновлять все Obstacle (вызывать obstacle.Update(dt) для каждого)
   - Вызывать Spawner.Update(dt) для генерации новых препятствий
   - Проверять коллизии между Player и каждым Obstacle
   - Удалять препятствия, вышедшие за границы экрана
6. WHEN препятствие выходит за границы экрана (z < -10), THE World SHALL удалять его из списка и возвращать в ObjectPool
7. WHEN коллизия между Player и Obstacle обнаружена, THE World SHALL переходить в состояние game over
8. THE World SHALL использовать только конкретные типы (НЕ map[string]interface{})
9. THE World SHALL получать Config при инициализации через конструктор NewWorld(config *Config)
10. THE World SHALL использовать простую проверку коллизий (AABB или distance-based)
11. THE World SHALL ограничивать максимальное количество активных препятствий (например, не более 20)


### ЭТАП 2: ГЕЙМПЛЕЙ (Приоритет: ВЫСОКИЙ)

#### Требование 5: Система Spawner

**User Story:** Как разработчик, я хочу систему Spawner для управления генерацией препятствий.

**Acceptance Criteria:**

1. THE Spawner SHALL быть определён в `internal/game/spawner.go`
2. THE Spawner SHALL иметь поле timer для отслеживания времени
3. THE Spawner SHALL иметь метод Update(dt float64) возвращающий *Obstacle или nil
4. THE Spawner SHALL генерировать препятствия с настраиваемой частотой (из Config)
5. THE Spawner SHALL использовать конфигурацию для параметров спавна
6. THE Spawner SHALL использовать ObstaclePool для получения препятствий

#### Требование 6: Object Pool для Препятствий

**User Story:** Как разработчик, я хочу простой Object Pool для переиспользования препятствий, чтобы снизить нагрузку на GC.

**Acceptance Criteria:**

1. THE ObstaclePool SHALL быть определён в `internal/pools/obstacle_pool.go`
2. THE ObstaclePool SHALL иметь простую структуру: pool []*Obstacle
3. THE ObstaclePool SHALL иметь метод Get() для получения препятствия из пула
4. WHEN пул пуст, THE ObstaclePool SHALL создавать новое препятствие
5. THE ObstaclePool SHALL иметь метод Put(obstacle *Obstacle) для возврата в пул
6. THE ObstaclePool SHALL полностью сбрасывать состояние объекта при возврате (позиция, скорость, тип, активность)
7. THE Spawner SHALL использовать ObstaclePool вместо прямого создания
8. THE ObstaclePool SHALL NOT использовать generics (простой []*Obstacle)

#### Требование 7: Конфигурация Игры

**User Story:** Как разработчик, я хочу вынести хардкод в конфигурацию, чтобы легко менять параметры игры.

**Acceptance Criteria:**

1. THE GameConfig SHALL быть определён в `internal/config/config.go`
2. THE GameConfig SHALL содержать: PlayerSpeed, SpawnRate, Gravity, JumpForce, MaxObstacles
3. THE GameConfig SHALL загружаться из значений по умолчанию
4. THE GameConfig SHALL быть доступен всем игровым объектам через передачу в конструкторы
5. THE Config SHALL NOT быть глобальной переменной
6. THE Config SHALL передаваться через конструкторы: Game → World → Player/Spawner
7. FOR ALL магических чисел в коде, они SHALL быть заменены на значения из конфигурации

### ЭТАП 3: ВИЗУАЛ И UX (Приоритет: СРЕДНИЙ)

#### Требование 8: Resource Manager (Упрощённый)

**User Story:** Как разработчик, я хочу простой Resource Manager для загрузки текстур один раз.

**Acceptance Criteria:**

1. THE ResourceManager SHALL быть определён в `internal/resources/manager.go`
2. THE ResourceManager SHALL иметь map[string]*ebiten.Image для кэширования
3. THE ResourceManager SHALL иметь метод LoadImage(path string) (*ebiten.Image, error)
4. WHEN текстура уже загружена, THE ResourceManager SHALL возвращать закешированную версию
5. THE ResourceManager SHALL загружать каждую текстуру не более одного раза
6. WHEN загрузка изображения не удалась, THE ResourceManager SHALL возвращать ошибку (НЕ fallback)
7. THE ResourceManager SHALL использовать относительные пути от папки assets/
8. THE ResourceManager SHALL NOT использовать сложные системы предзагрузки

#### Требование 9: Простой UI

**User Story:** Как разработчик, я хочу простой UI без overengineering.

**Acceptance Criteria:**

1. THE MenuUI SHALL быть определён в `internal/ui/menu.go`
2. THE MenuUI SHALL содержать простые кнопки и текст
3. THE MenuUI SHALL NOT использовать сложные UI фреймворки
4. THE MenuUI SHALL иметь методы Update() и Draw()
5. THE UI код SHALL быть минималистичным и понятным

#### Требование 10: Звуковые Эффекты

**User Story:** Как разработчик, я хочу добавить базовые звуковые эффекты.

**Acceptance Criteria:**

1. THE SoundManager SHALL быть определён в `internal/audio/sound.go`
2. THE SoundManager SHALL поддерживать звуки: прыжок, столкновение, фоновая музыка
3. THE SoundManager SHALL иметь простой интерфейс Play(soundName string)
4. THE SoundManager SHALL использовать ResourceManager для загрузки звуков
5. THE SoundManager SHALL NOT использовать сложные аудио-системы

### ЭТАП 4: ЧИСТКА (Приоритет: НИЗКИЙ)

#### Требование 11: Удаление Мусора

**User Story:** Как разработчик, я хочу удалить весь мусорный код, чтобы упростить навигацию.

**Acceptance Criteria:**

1. THE рефакторинг SHALL удалить пустые директории: internal/ui/animations, components, core, styles, widgets
2. THE рефакторинг SHALL удалить весь закомментированный код
3. THE рефакторинг SHALL удалить неиспользуемые геттеры и сеттеры
4. THE рефакторинг SHALL заменить избыточные геттеры прямым доступом к полям (Go идиома)
5. FOR ALL оставшихся файлов, каждый SHALL быть активно используемым

#### Требование 12: Простое Логирование

**User Story:** Как разработчик, я хочу простое логирование для отладки.

**Acceptance Criteria:**

1. THE логирование SHALL использовать стандартный log.Println()
2. THE логирование SHALL NOT использовать сложные системы логирования
3. THE логирование SHALL выводить ключевые события: "Player jumped", "Obstacle spawned", "Collision detected"
4. THE логирование SHALL быть минималистичным
5. THE логирование SHALL легко отключаться через флаг DEBUG

### Требование 13: Обратная Совместимость Функциональности

**User Story:** Как игрок, я хочу, чтобы игра работала так же после рефакторинга.

**Acceptance Criteria:**

1. THE рефакторинг SHALL сохранить механику балансирования игрока
2. THE рефакторинг SHALL сохранить механику прыжков
3. THE рефакторинг SHALL сохранить визуальный стиль (псевдо-3D, текстуры)
4. THE рефакторинг SHALL сохранить систему подсчёта очков
5. THE рефакторинг SHALL сохранить условия game over
6. THE рефакторинг SHALL сохранить управление (клавиши A/D/W)
7. FOR ALL игровых механик, поведение SHALL быть идентичным
