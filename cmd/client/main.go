package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gophkeeper/gophkeeper/internal/client/tui"
)

var (
	version   = "dev"
	buildDate = "unknown"
	server    = flag.String("server", "localhost:50051", "Server address")
)

func main() {
	// Показываем версию при запросе (до Parse, чтобы -v/--version не считались неизвестными флагами)
	for _, arg := range os.Args[1:] {
		if arg == "-v" || arg == "--version" || arg == "version" {
			fmt.Printf("GophKeeper Client\nVersion: %s\nBuild Date: %s\n", version, buildDate)
			os.Exit(0)
		}
	}
	flag.Parse()

	// Создаём модель приложения
	app, err := tui.NewAppModel(*server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating app: %v\n", err)
		os.Exit(1)
	}
	defer app.Close()

	// Запускаем приложение
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}
