package main

func (m model) View() string {

	return baseStyle.Render(m.table.View())

}
