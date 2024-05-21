package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"win-files-to-print/pkg/print"
	"win-files-to-print/pkg/ui"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

type winservice struct{}

var isFirstRun bool = true
var elog *eventlog.Log
var eventid uint32 = 10002

const svcDescription = "Win files to print"

func main() {
	const svcName = "WinFilesToPrint"

	elog, _ = eventlog.Open(svcName)
	defer elog.Close()

	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		elog.Error(eventid, fmt.Sprint("Failed to determine if we are running in an interactive session."))
	}

	if !isIntSess {
		runService(svcName)
		return
	}

	if len(os.Args) < 2 {
		usage("No command specified!!!")
	}

	cmd := strings.ToLower(os.Args[1])
	switch cmd {
	case "list":
		if print.CheckLibPdfExist() == false {
			log.Fatalf("Lib pdf does not exist: %s", print.DefaultSumatraPath)
		}
		printList, _ := print.NewScanner().PrintList()
		config, _ := print.NewConfigPrinter()
		cli := ui.NewCLI(config)
		cli.ChoosePrinter(printList)
		cli.ChooseFolder()
		folderScan := print.NewFolderScan(config.GetFolder())
		err := folderScan.ScanPDFFiles()
		if err != nil {
			log.Fatal(err)
		}
		print.NewPrinter(config.PrintName).Print("h:\\test.pdf")
	case "test_print":
		print.NewPrinter("Microsoft Print to PDF").Print("h:\\test.pdf")
	case "configure":
		if print.CheckLibPdfExist() == false {
			log.Fatalf("Lib pdf does not exist: %s", print.DefaultSumatraPath)
		}
		printList, _ := print.NewScanner().PrintList()
		config, _ := print.NewConfigPrinter()
		cli := ui.NewCLI(config)
		cli.ChoosePrinter(printList)
		cli.ChooseFolder()
		err = config.SaveFile()
		if err != nil {
			log.Fatal(err)
		}
	case "install":
		err = installService(svcName, svcDescription)
	case "remove":
		err = removeService(svcName)
	case "start":
		err = startService(svcName)
	case "run":
		runApp()
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped)
	default:
		usage(fmt.Sprintf("Invalid command %s", cmd))
	}
	if err != nil {
		elog.Error(eventid, fmt.Sprintf("Failed to %s %s: %v", cmd, svcName, err))
		fmt.Printf("Failed to %s %s: %v", cmd, svcName, err)
	}
	return
}

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr, "%s\n\n"+
		"Usage: %s <command>\n"+
		"        where <command> is one of\n"+
		"        install, remove, start or stop.\n",
		errmsg, os.Args[0])

	os.Exit(2)
}

func runService(name string) {
	elog.Info(eventid, fmt.Sprintf("Starting %s service.", name))

	run := svc.Run
	if err := run(name, &winservice{}); err != nil {
		elog.Error(eventid, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}

	elog.Info(eventid, fmt.Sprintf("%s service stopped.", name))
}

func (m *winservice) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	fasttick := time.Tick(500 * time.Millisecond)
	slowtick := time.Tick(2 * time.Second)
	tick := fasttick
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case <-tick:
			runApp()
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				elog.Info(eventid, "Shutdown Win files to print")
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				tick = slowtick
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				tick = fasttick
			default:
				elog.Error(eventid, fmt.Sprintf("Unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}

func installService(name, desc string) error {
	elog.Info(eventid, fmt.Sprintf("Install %s service.", name))

	exepath, err := exePath()
	if err != nil {
		return err
	}
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", name)
	}
	s, err = m.CreateService(name, exepath, mgr.Config{DisplayName: desc}, "is", "auto-started")
	if err != nil {
		return err
	}
	defer s.Close()
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("SetupelogSource() failed: %s", err)
	}
	return nil
}

func removeService(name string) error {
	elog.Info(eventid, fmt.Sprintf("Remove %s service.", name))

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer s.Close()
	err = s.Delete()
	if err != nil {
		return err
	}
	err = eventlog.Remove(name)
	if err != nil {
		return fmt.Errorf("RemoveelogSource() failed: %s", err)
	}
	return nil
}

func startService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

func controlService(name string, c svc.Cmd, to svc.State) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()
	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access service: %v", err)
	}
	defer s.Close()
	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control=%d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	return nil
}

func runApp() {
	if isFirstRun {
		isFirstRun = false
		elog.Info(eventid, fmt.Sprintf("Start %s", svcDescription))

		if print.CheckLibPdfExist() == false {
			log.Fatalf("Lib pdf does not exist: %s", print.DefaultSumatraPath)
		}

		config, err := print.NewConfigPrinter()
		err = config.LoadConfig()
		if err != nil {
			elog.Error(eventid, err.Error())
			elog.Info(eventid, err.Error())
			log.Fatal(err)
		}

		// Canal pour recevoir les signaux d'interruption
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

		go func() {

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			for {
				elog.Info(eventid, fmt.Sprintf("Itération toutes les 5 secondes %s %s %d", time.Now().String(), config.Folder, config.PrintName))
				fmt.Println("Itération toutes les 5 secondes:", time.Now())
				time.Sleep(5 * time.Second)
			}
		}()
		// Attendre le signal d'arrêt
		sig := <-sigChan
		elog.Info(eventid, fmt.Sprintf("Signal reçu, arrêt:", sig))
		fmt.Println("Signal reçu, arrêt:", sig)
	}
}
