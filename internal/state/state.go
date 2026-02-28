package state

import "github.com/hajimehoshi/ebiten/v2"

type State interface {
	Update() error
	Draw(screen *ebiten.Image)
	// Enter вызывается при переходе в это состояние (можно передать параметры)
	Enter(prevState State, data interface{})
	// Exit вызывается перед выходом из состояния
	Exit()
}
