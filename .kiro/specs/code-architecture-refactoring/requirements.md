# Документ требований: Рефакторинг архитектуры Go-проекта игры

## Введение

Данный документ описывает требования к рефакторингу архитектуры игрового проекта на Go. Цель рефакторинга - улучшить качество кодовой базы путём унификации интерфейсов, инкапсуляции данных, устранения магических чисел и правильного разделения ответственности между компонентами. Рефакторинг должен сохранить всю существующую функциональность при улучшении поддерживаемости и читаемости кода.

## Глоссарий

- **System**: Кодовая база игрового проекта на Go
- **WorldContext**: Интерфейс, предоставляющий информацию о состоянии игрового мира
- **Refactoring_Tool**: Инструмент или процесс, выполняющий рефакторинг кода
- **Structure**: Структура данных Go (Camera, Segment, World, Layer)
- **Magic_Number**: Числовой литерал в коде без явного семантического значения
- **UI_Component**: Компонент пользовательского интерфейса (например, BalanceBarLayer)
- **Game_World**: Компоненты, представляющие игровой мир (слои, сегменты)
- **Config_Package**: Пакет internal/config для хранения конфигурационных констант
- **Getter**: Метод для чтения значения приватного поля
- **Setter**: Метод для изменения значения приватного поля

## Требования

### Требование 1: Унификация интерфейса WorldContext

**Пользовательская история**: Как разработчик, я хочу иметь единое определение интерфейса WorldContext, чтобы избежать дублирования кода и упростить поддержку.

#### Критерии приёмки

1. THE System SHALL contain exactly one definition of WorldContext interface in package internal/common
2. WHEN the System is analyzed, THEN all packages except internal/common SHALL NOT define WorldContext interface
3. WHEN any file uses WorldContext, THEN that file SHALL import it from internal/common package
4. THE System SHALL compile without errors after WorldContext unification
5. WHEN WorldContext unification is complete, THEN all existing functionality SHALL remain unchanged

### Требование 2: Инкапсуляция полей структуры Camera

**Пользовательская история**: Как разработчик, я хочу инкапсулировать поля структуры Camera, чтобы контролировать доступ к внутреннему состоянию и предотвратить некорректные изменения.

#### Критерии приёмки

1. THE Camera structure SHALL have all fields as private (lowercase first letter)
2. FOR each field in Camera structure, THE System SHALL provide a getter method
3. FOR each mutable field in Camera structure, THE System SHALL provide a setter method
4. WHEN external code accesses Camera fields, THEN it SHALL use getter methods
5. WHEN external code modifies Camera fields, THEN it SHALL use setter methods
6. THE System SHALL compile without errors after Camera encapsulation

### Требование 3: Инкапсуляция полей структуры Segment

**Пользовательская история**: Как разработчик, я хочу инкапсулировать поля структуры Segment, чтобы обеспечить целостность данных сегментов игрового мира.

#### Критерии приёмки

1. THE Segment structure SHALL have all fields as private
2. FOR each field in Segment structure, THE System SHALL provide a getter method
3. FOR each mutable field in Segment structure, THE System SHALL provide a setter method
4. WHEN external code accesses Segment fields, THEN it SHALL use getter methods
5. THE System SHALL compile without errors after Segment encapsulation

### Требование 4: Инкапсуляция полей структуры World

**Пользовательская история**: Как разработчик, я хочу инкапсулировать поля структуры World, чтобы защитить состояние игрового мира от прямого внешнего доступа.

#### Критерии приёмки

1. THE World structure SHALL have all fields as private
2. FOR each field in World structure, THE System SHALL provide a getter method
3. FOR each mutable field in World structure, THE System SHALL provide a setter method
4. WHEN external code accesses World fields, THEN it SHALL use getter methods
5. THE System SHALL compile without errors after World encapsulation

### Требование 5: Инкапсуляция полей слоёв (Layers)

**Пользовательская история**: Как разработчик, я хочу инкапсулировать поля всех структур слоёв (SkyLayer, FarBankLayer, SegmentLayer), чтобы обеспечить единообразный подход к управлению данными.

#### Критерии приёмки

1. THE SkyLayer structure SHALL have all fields as private
2. THE FarBankLayer structure SHALL have all fields as private
3. THE SegmentLayer structure SHALL have all fields as private
4. FOR each layer structure, THE System SHALL provide getter methods for all fields
5. FOR each layer structure, THE System SHALL provide setter methods for mutable fields
6. THE System SHALL compile without errors after layer encapsulation

### Требование 6: Создание пакета конфигурации

**Пользовательская история**: Как разработчик, я хочу централизованное хранение конфигурационных констант, чтобы легко находить и изменять параметры игры.

#### Критерии приёмки

1. THE System SHALL contain package internal/config
2. THE Config_Package SHALL define CameraConfig structure with camera-related constants
3. THE Config_Package SHALL define GameConfig structure with game-related constants
4. THE Config_Package SHALL define PhysicsConfig structure with physics-related constants
5. THE Config_Package SHALL provide global instances of configuration structures
6. THE System SHALL compile without errors after config package creation

### Требование 7: Устранение магических чисел

**Пользовательская история**: Как разработчик, я хочу заменить все магические числа именованными константами, чтобы код был более понятным и легче поддерживался.

#### Критерии приёмки

1. WHEN the System is analyzed, THEN all numeric literals in camera.go SHALL be replaced with named constants from Config_Package
2. WHEN the System is analyzed, THEN all numeric literals in game_state.go SHALL be replaced with named constants from Config_Package
3. WHEN the System is analyzed, THEN all numeric literals in player.go SHALL be replaced with named constants from Config_Package
4. WHEN a file uses configuration constants, THEN that file SHALL import internal/config package
5. THE System SHALL compile without errors after magic number extraction
6. WHEN magic numbers are replaced, THEN the System behavior SHALL remain identical

### Требование 8: Разделение UI и игрового мира

**Пользовательская история**: Как разработчик, я хочу отделить UI-компоненты от игрового мира, чтобы обеспечить правильное разделение ответственности и упростить поддержку обеих подсистем.

#### Критерии приёмки

1. THE BalanceBarLayer SHALL be moved from internal/world package to internal/ui package
2. WHEN BalanceBarLayer is in ui package, THEN it SHALL NOT implement world.Layer interface
3. THE GameState SHALL render UI_Component separately from Game_World
4. WHEN GameState draws the screen, THEN it SHALL draw Game_World first and UI_Component after
5. THE System SHALL NOT have circular dependencies between world and ui packages
6. THE System SHALL compile without errors after UI reorganization
7. WHEN UI is reorganized, THEN the BalanceBar SHALL display correctly on screen

### Требование 9: Сохранение функциональности

**Пользовательская история**: Как пользователь игры, я хочу, чтобы после рефакторинга все функции работали так же, как и раньше, чтобы мой игровой опыт не изменился.

#### Критерии приёмки

1. WHEN refactoring is complete, THEN all existing game features SHALL work identically to before refactoring
2. WHEN the player moves, THEN movement SHALL behave the same as before refactoring
3. WHEN the player jumps, THEN jump mechanics SHALL behave the same as before refactoring
4. WHEN the balance bar is displayed, THEN it SHALL show the same information as before refactoring
5. WHEN the camera renders the scene, THEN the visual output SHALL be identical to before refactoring

### Требование 10: Компилируемость на каждом шаге

**Пользовательская история**: Как разработчик, я хочу, чтобы проект компилировался после каждого шага рефакторинга, чтобы можно было легко обнаружить и исправить ошибки.

#### Критерии приёмки

1. WHEN WorldContext unification is complete, THEN the System SHALL compile without errors
2. WHEN any Structure encapsulation is complete, THEN the System SHALL compile without errors
3. WHEN magic numbers are extracted, THEN the System SHALL compile without errors
4. WHEN UI reorganization is complete, THEN the System SHALL compile without errors
5. WHEN any refactoring step fails compilation, THEN the Refactoring_Tool SHALL report the error and stop

### Требование 11: Согласованность именования

**Пользовательская история**: Как разработчик, я хочу единообразного именования геттеров и сеттеров, чтобы код был предсказуемым и легко читаемым.

#### Критерии приёмки

1. FOR each private field, THE Getter method name SHALL match the field name with capitalized first letter
2. FOR each private field, THE Getter SHALL return the field type
3. FOR each mutable private field, THE Setter method name SHALL be "Set" + field name with capitalized first letter
4. FOR each Setter, THE method SHALL accept exactly one parameter of the field type
5. WHEN a getter or setter name conflicts with existing method, THEN the Refactoring_Tool SHALL use "Get" prefix for getter

### Требование 12: Отсутствие циклических зависимостей

**Пользовательская история**: Как разработчик, я хочу избежать циклических зависимостей между пакетами, чтобы архитектура оставалась чистой и поддерживаемой.

#### Критерии приёмки

1. THE System SHALL NOT have circular dependencies between any two packages
2. WHEN ui package is created, THEN it SHALL NOT depend on world package
3. WHEN world package is analyzed, THEN it SHALL NOT depend on ui package
4. THE Config_Package SHALL NOT depend on any other internal packages
5. THE internal/common package SHALL NOT depend on any other internal packages except core
