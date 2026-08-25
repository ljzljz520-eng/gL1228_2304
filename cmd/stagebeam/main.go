package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"stagebeam/internal/config"
	"stagebeam/internal/persist"
	"stagebeam/internal/render"
	"stagebeam/internal/service"
	"stagebeam/internal/transport"
)

func main() {
	dataPath := flag.String("data", "stagebeam.db", "embedded stage data file")
	demo := flag.Bool("demo", false, "run the deterministic mouse-free demo")
	flag.Parse()
	store, err := persist.Open(*dataPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer store.Close()
	showService := service.NewShowService(store)
	if *demo || flag.NArg() == 0 {
		runDemo(showService)
		return
	}
	commandData, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	command, err := transport.DecodeCommand(commandData)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	response := transport.Execute(showService, command)
	encoded, encodeErr := transport.EncodeResponse(response)
	if encodeErr != nil {
		fmt.Fprintln(os.Stderr, encodeErr)
		return
	}
	fmt.Println(string(encoded))
}

func runDemo(showService *service.ShowService) {
	settings := config.DefaultSettings()
	show, err := showService.CreateShow("demo-show", "舞台光束手势秀", settings)
	if err != nil {
		fmt.Println(err)
		return
	}
	frames, err := showService.Demonstrate(show.ID, render.DemoSequence())
	if err != nil {
		fmt.Println(err)
		return
	}
	last := ""
	if len(frames) > 0 {
		last = frames[len(frames)-1].Message
	}
	result := map[string]any{"show": show.Name, "frames": len(frames), "last_message": last, "mouse_hint": render.PointerLabel(render.MousePoint{X: 640, Y: 240}, 1280, 720)}
	data, _ := json.Marshal(result)
	fmt.Println(string(data))
}
