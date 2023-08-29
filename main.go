package main

import (
	"fmt"
	"os"
	"time" // We'll need the time package for the timer functionality

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	totalTime    time.Duration
	studyTime    time.Duration // Pomodoro study time duration
	breakTime    time.Duration
	bigBreakTime time.Duration
	running      bool          // Whether the timer is running
	remaining    time.Duration // Remaining time
	isBreak      bool
	elapsedTime  time.Duration
	cycle        int
}

// Init implements tea.Model.
func (model) Init() tea.Cmd {
	return nil
}

func initialModel(totalTime time.Duration) *model {
	studyTime := 3 * time.Second
	breakTime := 2 * time.Second
	bigBreakTime := 10 * time.Second
	return &model{
		totalTime:    totalTime,
		studyTime:    studyTime,
		breakTime:    breakTime,
		bigBreakTime: bigBreakTime,
		running:      false,
		remaining:    studyTime,
		isBreak:      false,
		elapsedTime:  0,
	}
}
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "s":
			if !m.running {

				m.cycle += 1
				return &model{running: true, remaining: m.studyTime, isBreak: false, elapsedTime: m.elapsedTime, studyTime: m.studyTime, breakTime: m.breakTime, cycle: m.cycle}, tickCmd
			}
		}
	case tickMsg:
		if m.remaining <= 0 {
			if m.isBreak {
				m.cycle += 1

				return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: true, remaining: m.studyTime, isBreak: false, elapsedTime: m.elapsedTime, cycle: m.cycle}, tickCmd

			} else if !m.isBreak {
				if (m.cycle % 4) == 0 {
					return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: true, remaining: (10 * time.Second), isBreak: true, elapsedTime: m.elapsedTime, cycle: m.cycle}, tickCmd
				} else {
					remainingBreakTime := m.breakTime - (m.elapsedTime % m.breakTime)
					return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: true, remaining: remainingBreakTime, isBreak: true, elapsedTime: m.elapsedTime, cycle: m.cycle}, tickCmd
				}
			}
		}
		return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: true, remaining: m.remaining - tickInterval, isBreak: m.isBreak, elapsedTime: m.elapsedTime, cycle: m.cycle}, tickCmd
	}
	return m, nil
}

type tickMsg struct{}

const tickInterval = time.Second

var tickCmd = tea.Tick(tickInterval, func(time.Time) tea.Msg {
	return tickMsg{}
})

func (m *model) View() string {
	var status string
	if m.running {
		if m.isBreak {
			if (m.cycle % 4) == 0 {
				status = fmt.Sprintf("*********************\nBig Break time! Time remaining: %s, Elapsed: %v cycles", m.remaining.String(), m.cycle)
			} else {
				status = fmt.Sprintf("*********************\nBreak time! Time remaining: %s, Elapsed: %v cycles", m.remaining.String(), m.cycle)
			}
		} else {
			status = fmt.Sprintf("*********************\nStudy time! Time remaining: %s, Elapsed: %v cycles", m.remaining.String(), m.cycle)
		}
	} else {
		status = "Press 's' to start the timer."
	}
	return fmt.Sprintf("%s\nPress q to quit.", status)
}

func main() {
	var studyMinutes int
	fmt.Print("Enter total study time (minutes): ")

	_, err := fmt.Scan(&studyMinutes)
	if err != nil {
		fmt.Printf("Error reading input: %v", err)
		return
	}

	totalTime := time.Duration(studyMinutes) * time.Minute
	p := tea.NewProgram(initialModel(totalTime))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
