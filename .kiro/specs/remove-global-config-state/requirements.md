# Requirements Document

## Introduction

Данный документ описывает требования к рефакторингу системы конфигурации для устранения глобальных переменных и зависимости от глобального состояния. Цель рефакторинга - сделать код более тестируемым, модульным и гибким, позволяя использовать несколько конфигураций одновременно и упрощая unit-тестирование компонентов.

## Glossary

- **Config_Package**: Пакет internal/config, содержащий типы и функции конфигурации
- **Camera**: Структура камеры в пакете internal/render
- **Project_Function**: Функция проекции 3D координат на экран
- **Game_State**: Структура состояния игры в пакете internal/state
- **Player**: Структура игрока в пакете internal/entity
- **Global_Config_Variables**: Глобальные переменные config.Camera, config.Game, config.Physics
- **Constructor**: Функция создания экземпляра структуры (например, NewCamera, NewPlayer)
- **Dependency_Injection**: Паттерн передачи зависимостей через параметры конструктора

## Requirements

### Requirement 1: Удаление глобальных переменных конфигурации

**User Story:** Как разработчик, я хочу удалить глобальные переменные конфигурации, чтобы код не зависел от глобального состояния и был более тестируемым.

#### Acceptance Criteria

1. THE Config_Package SHALL NOT contain global variables Camera, Game, or Physics
2. THE Config_Package SHALL retain functions DefaultCameraConfig(), DefaultGameConfig(), and DefaultPhysicsConfig()
3. THE Config_Package SHALL retain type definitions CameraConfig, GameConfig, and PhysicsConfig
4. WHEN the refactoring is complete, THEN no code SHALL reference config.Camera, config.Game, or config.Physics as global variables

### Requirement 2: Явная передача конфигурации через конструкторы

**User Story:** Как разработчик, я хочу передавать конфигурации явно через конструкторы, чтобы каждый компонент получал только необходимые ему настройки.

#### Acceptance Criteria

1. THE Camera Constructor SHALL accept parameters (screenW float64, screenH float64, cfg CameraConfig)
2. THE Player Constructor SHALL accept PlayerConfig containing PhysicsConfig field
3. THE Game_State Constructor SHALL accept parameters (manager, gameCfg GameConfig, cameraCfg CameraConfig, physicsCfg PhysicsConfig)
4. WHEN a Constructor is called, THE component SHALL store the provided configuration for later use
5. THE main.go file SHALL create configuration instances using Default functions and pass them to Constructors

### Requirement 3: Преобразование функции Project в метод Camera

**User Story:** Как разработчик, я хочу сделать функцию Project методом структуры Camera, чтобы устранить зависимость от глобальной конфигурации камеры.

#### Acceptance Criteria

1. THE Camera structure SHALL contain a method Project(point core.Vec3) (float64, float64, float64)
2. THE Camera structure SHALL store CameraConfig or its relevant fields for curve calculations
3. THE Project method SHALL use Camera's stored configuration instead of config.Camera global variable
4. WHEN Project is called, THE method SHALL perform the same calculations as the original function
5. THE render package SHALL NOT contain a standalone Project function after refactoring
6. FOR ALL code calling render.Project, THE calls SHALL be updated to cam.Project()

### Requirement 4: Обновление использования конфигураций в Game_State

**User Story:** Как разработчик, я хочу обновить game_state.go для использования переданной конфигурации, чтобы устранить обращения к глобальным переменным.

#### Acceptance Criteria

1. THE Game_State structure SHALL store GameConfig as a field
2. WHEN Game_State accesses screen dimensions, THE structure SHALL use stored GameConfig.ScreenWidth and GameConfig.ScreenHeight
3. WHEN Game_State checks drift threshold, THE structure SHALL use stored GameConfig.DriftThreshold
4. THE Game_State SHALL NOT reference config.Game global variable

### Requirement 5: Обновление использования конфигураций в Player

**User Story:** Как разработчик, я хочу обновить player.go для использования переданной конфигурации физики, чтобы устранить обращения к глобальной переменной Physics.

#### Acceptance Criteria

1. THE Player structure SHALL store PhysicsConfig as a field
2. WHEN Player calculates jump physics, THE structure SHALL use stored PhysicsConfig.Gravity
3. THE Player SHALL NOT reference config.Physics global variable

### Requirement 6: Обновление всех остальных мест использования глобальных конфигураций

**User Story:** Как разработчик, я хочу обновить все оставшиеся места использования глобальных конфигураций, чтобы полностью устранить зависимость от глобального состояния.

#### Acceptance Criteria

1. THE gameover_state.go SHALL receive GameConfig through constructor or Game_State reference
2. THE camera.go initialization SHALL use passed CameraConfig instead of config.Camera
3. THE world layers (SkyLayer) SHALL receive screen dimensions as parameters instead of accessing config.Game
4. WHEN refactoring is complete, THE codebase SHALL have zero references to config.Camera.*, config.Game.*, or config.Physics.*
5. FOR ALL files in the project, a search for "config\\.(Camera|Game|Physics)\\." SHALL return zero matches

### Requirement 7: Обратная совместимость инициализации

**User Story:** Как разработчик, я хочу сохранить простоту инициализации приложения, чтобы изменения не усложнили код запуска.

#### Acceptance Criteria

1. THE main.go SHALL create default configurations using DefaultCameraConfig(), DefaultGameConfig(), DefaultPhysicsConfig()
2. THE initialization code SHALL remain concise and readable
3. WHEN the application starts with default configurations, THE behavior SHALL be identical to the original implementation
4. THE refactoring SHALL NOT introduce more than 10 additional lines of code in main.go

### Requirement 8: Сохранение функциональности

**User Story:** Как пользователь, я хочу, чтобы игра работала идентично после рефакторинга, чтобы изменения не повлияли на игровой процесс.

#### Acceptance Criteria

1. WHEN the refactored code runs, THE visual output SHALL be identical to the original implementation
2. WHEN the refactored code runs, THE game physics SHALL behave identically to the original implementation
3. WHEN the refactored code runs, THE camera projection SHALL produce identical results to the original implementation
4. FOR ALL game scenarios, THE refactored code SHALL produce the same behavior as the original code
