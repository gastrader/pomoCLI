package main

import (
	"fmt"
	"os"

	"time" // We'll need the time package for the timer functionality

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	padding  = 2
	maxWidth = 80
)

type model struct {
	totalTime       time.Duration
	studyTime       time.Duration // Pomodoro study time duration
	breakTime       time.Duration
	bigBreakTime    time.Duration
	running         bool          // Whether the timer is running
	remaining       time.Duration // Remaining time
	isBreak         bool
	elapsedTime     time.Duration
	cycle           int
	percent         float64
	studyPercent    float64
	breakPercent    float64
	bigBreakPercent float64
	progress        progress.Model
}

func (*model) Init() tea.Cmd {
	return nil
}

func initialModel(totalTime, studyTime, breakTime, bigBreakTime time.Duration, progress progress.Model) *model {
	return &model{
		totalTime:       totalTime,
		studyTime:       studyTime,
		breakTime:       breakTime,
		bigBreakTime:    bigBreakTime,
		running:         false,
		remaining:       studyTime,
		isBreak:         false,
		elapsedTime:     0,
		studyPercent:    0.0,
		breakPercent:    0.0,
		bigBreakPercent: 0.0,
		progress:        progress,
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
				m.elapsedTime += m.studyTime
				m.cycle += 1
				return &model{running: true, remaining: m.studyTime, isBreak: false, elapsedTime: m.elapsedTime, studyTime: m.studyTime, breakTime: m.breakTime, cycle: m.cycle, bigBreakTime: m.bigBreakTime, totalTime: m.totalTime}, tickCmd()
			}
		}
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil
	case tickMsg:

		if m.elapsedTime <= m.totalTime {
			if m.remaining <= 0 {
				if m.isBreak {
					m.cycle += 1
					m.elapsedTime += m.studyTime
					
					return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: true, remaining: m.studyTime, isBreak: false, elapsedTime: m.elapsedTime, cycle: m.cycle, percent: m.percent}, tickCmd()

				} else {
					if (m.cycle % 4) == 0 {
						m.elapsedTime += m.bigBreakTime
						
						return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: true, remaining: m.bigBreakTime, isBreak: true, elapsedTime: m.elapsedTime, cycle: m.cycle, percent: m.percent}, tickCmd()
					} else {
						m.elapsedTime += m.breakTime
						
						return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: true, remaining: m.breakTime, isBreak: true, elapsedTime: m.elapsedTime, cycle: m.cycle, percent: m.percent}, tickCmd()
					}
				}
			}
		} else {
			return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: false, remaining: m.remaining - tickInterval, isBreak: m.isBreak, elapsedTime: m.elapsedTime, cycle: m.cycle, percent: m.percent}, tea.Quit
		}

		return &model{totalTime: m.totalTime, studyTime: m.studyTime, breakTime: m.breakTime, bigBreakTime: m.bigBreakTime, running: true, remaining: m.remaining - tickInterval, isBreak: m.isBreak, elapsedTime: m.elapsedTime, cycle: m.cycle, studyPercent: (1- ((float64(m.remaining-time.Second) / 1000000000) / float64((m.studyTime)/1000000000))), breakPercent: (1 - ((float64(m.remaining-time.Second) / 1000000000) / float64((m.breakTime)/1000000000))), bigBreakPercent: (1 - ((float64(m.remaining-time.Second) / 1000000000) / float64((m.bigBreakTime)/1000000000)))}, tickCmd()
	}
	return m, tickCmd()
}

type tickMsg time.Time

const tickInterval = time.Second

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *model) View() string {
	var status string

	if m.running {
        
		if m.isBreak {
			if (m.cycle % 4) == 0 {
				status = fmt.Sprintf("*********************\nBig Break time! Time remaining: %s, Elapsed: %v, percent %.1f%%", m.remaining.String(), m.elapsedTime.Minutes(), m.bigBreakPercent*100)

			} else {
				status = fmt.Sprintf("*********************\nBreak time! Time remaining: %s, Elapsed: %v, percent %.1f%% ", m.remaining.String(), m.elapsedTime.Minutes(), m.breakPercent*100)

			}
		} else {
			status = fmt.Sprintf("*********************\nStudy time! Time remaining: %s, Elapsed: %v, percent %.1f%%", m.remaining.String(), m.elapsedTime.Minutes(), m.studyPercent*100)

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
	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	totalTime := time.Minute * time.Duration(studyMinutes)

	bt := time.Minute * 5
	st := time.Minute * 25
	bbt := time.Minute * 15
	p := tea.NewProgram(initialModel(totalTime, st, bt, bbt, prog))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
