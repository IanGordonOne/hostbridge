package main

import (
	"fmt"
	"io/ioutil"
	"runtime"

	"github.com/progrium/hostbridge/bridge/app"
	"github.com/progrium/hostbridge/bridge/core"
	"github.com/progrium/hostbridge/bridge/menu"
	"github.com/progrium/hostbridge/bridge/screen"
	"github.com/progrium/hostbridge/bridge/shell"
	"github.com/progrium/hostbridge/bridge/window"
)

func init() {
	runtime.LockOSThread()
}

var quitId uint16 = 999
var quitAllId uint16 = 9999

func tick(event core.Event) {
	if event.Type > 0 {
		fmt.Println("[tick] event", event, window.Focused())

		if event.Name == "close" || (event.Name == "menu" && event.MenuID == quitId) {
			w := window.FindByID(event.WindowID)
			if w != nil {
				w.Destroy()
			}

			all := window.All()
			fmt.Println("count of all windows", len(all))
			if len(all) == 0 {
				fmt.Println("  quitting application...")
				core.Quit()
			}
		}

		if event.Name == "menu" && event.MenuID == quitAllId {
			core.Quit()
		}
	}
}

func main() {
	main2()
	core.Run(tick)

	// NOTE(nick): this doesn't appear to be called ever
	fmt.Println("[main] Goodbye.")
}

func main2() {
	menuTemplate := []menu.Item{
		{
			// NOTE(nick): when setting the window menu with wry, the first item title will always be the name of the executable on MacOS
			// so, this property is ignored:
			// @Robustness: maybe we want to make that more visible to the user somehow?
			Title:   "this doesnt matter",
			Enabled: true,
			SubMenu: []menu.Item{
				{
					ID:          121,
					Title:       "About",
					Enabled:     true,
					Accelerator: "Control+I",
				},
				{
					ID:      122,
					Title:   "Disabled",
					Enabled: false,
				},
				{
					ID:          quitId,
					Title:       "Quit",
					Enabled:     true,
					Accelerator: "CommandOrControl+Q",
				},
			},
		},
		{
			ID:      23,
			Title:   "hello world",
			Enabled: true,
			SubMenu: []menu.Item{
				{
					ID:      777,
					Title:   "This is an amazing menu option",
					Enabled: true,
				},
			},
		},
	}

	m := menu.New(menuTemplate)
	app.SetMenu(m)

	trayTemplate := []menu.Item{
		{
			Title:   "Click on this here thing",
			Enabled: true,
		},
		{
			Title:   "Secret stuff",
			Enabled: true,
			SubMenu: []menu.Item{
				{
					ID:      42,
					Title:   "I'm nested!!",
					Enabled: true,
				},
				{
					ID:      101,
					Title:   "Can't touch this",
					Enabled: false,
				},
			},
		},
		{
			ID:          quitAllId,
			Title:       "Quit App",
			Enabled:     true,
			Accelerator: "Command+T",
		},
	}

	iconPath := "assets/icon.png"
	if runtime.GOOS == "windows" {
		iconPath = "assets/icon.ico"
	}

	iconData, err := ioutil.ReadFile(iconPath)
	if err != nil {
		fmt.Println("Error reading icon file:", err)
	}

	app.NewIndicator(iconData, trayTemplate)

	options := window.Options{
		Title: "Demo window",
		// NOTE(nick): resizing a transparent window on MacOS seems really slow?
		Transparent: true,
		Frameless:   false,
		Visible:     true,
		//Position: window.Position{X: 10, Y: 10},
		//Size: window.Size{ Width: 360, Height: 240 },
		Center: true,
		HTML: `
			<!doctype html>
			<html>
				<body style="font-family: -apple-system, BlinkMacSystemFont, avenir next, avenir, segoe ui, helvetica neue, helvetica, Ubuntu, roboto, noto, arial, sans-serif; background-color:rgba(87,87,87,0.8);"></body>
				<script>
					window.onload = function() {
						document.body.innerHTML = '<div style="padding: 30px">Transparency Test!<br><br>${navigator.userAgent}</div>';
					};
				</script>
			</html>
		`,
	}

	w2, _ := window.New(options)
	w2.SetTitle("YO!")
	w2.SetSize(core.Size{ Width: 640, Height: 480 })

	w1, _ := window.New(options)

	fmt.Println("[main] window", w1)

	if w1 == nil {
		return
	}

	w1.SetTitle("Hello, Sailor!")
	fmt.Println("[main] window position", w1.GetOuterPosition())

	shell.ShowNotification(shell.Notification{
		Title:    "Title: Hello, world",
		Subtitle: "Subtitle: MacOS only",
		Body:     "Body: This is the body",
	})

	if false {
		ok := shell.ShowMessage(shell.MessageDialog{
			Title:   "Title: what do you think?",
			Body:    "Body: about this description text",
			Level:   "warning",
			Buttons: "okcancel",
		})

		fmt.Println("ShowMessage ok", ok)
	}

	if false {
		files := shell.ShowFilePicker(shell.FileDialog{
			Title:   "Title: please pick a file...",
			Mode:    "pickfiles",
			Filters: []string{"txt,rs,cpp", "image:png,jpg,jpeg"},
		})

		fmt.Println("ShowFilePicker files", files, len(files))
	}

	success := shell.WriteClipboard("Hello from Go!")
	fmt.Println("Wrote clipboard data:", success)

	fmt.Println("Read clipboard data:", shell.ReadClipboard())

	displays := screen.Displays()
	fmt.Println("Displays:")

	for _, it := range displays {
		fmt.Println("", it.Name)
		fmt.Println("  Size:", it.Size)
		fmt.Println("  Position:", it.Position)
		fmt.Println("  ScaleFactor:", it.ScaleFactor)
	}

	didRegister1 := shell.RegisterShortcut("Control+Shift+R")
	fmt.Println("didRegister", didRegister1)

	didRegister2 := shell.RegisterShortcut("Control+Shift+T")
	fmt.Println("didRegister", didRegister2)

	didUnregister := shell.UnregisterShortcut("Control+Shift+T")
	fmt.Println("didUnregister", didUnregister)

}
