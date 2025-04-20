package sugar

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func Update[T tea.Model](m *T, msg tea.Msg, batch *[]tea.Cmd) {
	m2, cmd := (*m).Update(msg)
	m3, ok := m2.(T)
	if !ok {
		panic(fmt.Sprintf("Update returns not self: was %T, returns %T", m, m2))
	}

	*m = m3
	*batch = append(*batch, cmd)
}

type ModelStruct[T any] interface {
	Update(tea.Msg) (T, tea.Cmd)
}

func Update2[T ModelStruct[T]](m *T, msg tea.Msg, batch *[]tea.Cmd) {
	var cmd tea.Cmd
	*m, cmd = (*m).Update(msg)
	*batch = append(*batch, cmd)
}
