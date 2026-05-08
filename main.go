package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mainmodel "github.com/Keivan-sf/Bushuray-tui/components/MainModel"
	appconfig "github.com/Keivan-sf/Bushuray-tui/lib/AppConfig"
	"github.com/Keivan-sf/Bushuray-tui/global"
	connection "github.com/Keivan-sf/Bushuray-tui/lib/Connection"
	notif_publisher "github.com/Keivan-sf/Bushuray-tui/lib/NotifPublisher"
	servercmds "github.com/Keivan-sf/Bushuray-tui/lib/ServerCommands"
	"github.com/Keivan-sf/Bushuray-tui/utils"

	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   os.TempDir() + "/bushuray-debug.log",
		MaxSize:    20,
		MaxBackups: 1,
		MaxAge:     0,
		Compress:   false,
	})
	log.SetPrefix("debug: ")
	log.SetFlags(log.LstdFlags | log.Lmsgprefix)

	appconfig.LoadConfig()
	global.LoadTheme(appconfig.GetConfig().Theme)
	port := appconfig.GetConfig().CoreTCPPort

	C := connection.ConnectionHandler{}
	C.Init("127.0.0.1", port)

	err := C.GetConnection()
	if err != nil {
		fmt.Println("core was not found at", port, "trying to spawn")
		err := utils.SpawnBushurayCore()
		if err != nil {
			fmt.Println("failed to spawn core:", err)
			return
		}
		time.Sleep(1000 * time.Millisecond)
		err = C.GetConnection()
		if err != nil {
			fmt.Println("failed to connect to core:", err)
			return
		}
	} else {
		fmt.Println("connection established")
	}

	zone.NewGlobal()
	p := tea.NewProgram(mainmodel.InitModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())

	go C.HandleConnection(p)
	servercmds.Init(&C)
	servercmds.GetApplicationState()
	notif_publisher.Init(p)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error %v\n", err)
		os.Exit(1)
	}
}
