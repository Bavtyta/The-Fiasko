package state

import (
	"TheFiaskoTest/internal/config"

	"github.com/hajimehoshi/ebiten/v2"
)

type Manager struct {
	current    State
	next       State
	data       interface{}
	gameConfig config.GameConfig
}

func NewManager(initial State, gameCfg config.GameConfig) *Manager {
	return &Manager{
		current:    initial,
		gameConfig: gameCfg,
	}
}

func (m *Manager) Update() error {
	if m.next != nil {
		if m.current != nil { // проверка, чтобы не вызвать Exit на nil
			m.current.Exit()
		}
		m.current = m.next
		m.current.Enter(nil, m.data)
		m.next = nil
		m.data = nil
	}
	if m.current == nil {
		return nil // или вернуть ошибку, но лучше не допускать такого состояния
	}
	return m.current.Update()
}

func (m *Manager) Draw(screen *ebiten.Image) {
	m.current.Draw(screen)
}

func (m *Manager) ChangeState(state State, data interface{}) {
	m.next = state
	m.data = data
}

// GameConfig возвращает конфигурацию игры
func (m *Manager) GameConfig() config.GameConfig {
	return m.gameConfig
}
