package util

import "fmt"

func Colorize(color string, msg string) string {
	return fmt.Sprintf("\033[%sm%s\033[0m", color, msg)
}

const (
	// Terminal color codes used to print things in color
	Red          = "31"
	Green        = "32"
	Yellow       = "33"
	Blue         = "34"
	BrightYellow = "33;1"
)

func PrintLogo() {
	// This will look like this, but colored:
	//     _     _          .          _      _     _     _      _
	// 	  | |   | |        / \        | |    | |  / /    | |    | |  (R)
	// 	  | |..:::::'     /...\       | |    | | / /     | | :  | |
	// ...:::::::::   ..:::::::::.    | |    | |/ /      | |  ::::|
	// 	  | ''''| |     / ::::: \ '.  | |    | | \ \     [ [___:::::
	// 	  |_|   |_|    /_/     \_\    |_|    |_|  \_\     \______':::.
	//
	//                D O N A T I O N     T R A C K E R

	fmt.Println("         _     _          .          _      _     _     _      _")
	fmt.Println("        | |   | |        / \\        | |    | |  / /    | |    | |  (R)")
	fmt.Printf("        | |%s     /%s\\       | |    | | / /     | | %s  | |\n",
		Colorize(Green, "..:::::'"), Colorize(Yellow, "..."), Colorize(BrightYellow, ":"))
	fmt.Printf("     %s   %s    | |    | |/ /      | |  %s|\n",
		Colorize(Green, "...:::::::::"), Colorize(Yellow, "..:::::::::."), Colorize(BrightYellow, "::::"))
	fmt.Printf("        | %s| |     / %s \\ %s  | |    | | \\ \\     [ [___%s\n",
		Colorize(Green, "''''"), Colorize(Yellow, ":::::"), Colorize(Yellow, "'."),
		Colorize(BrightYellow, ":::::"))
	fmt.Printf("        |_|   |_|    /_/     \\_\\    |_|    |_|  \\_\\     \\______%s\n",
		Colorize(BrightYellow, "':::."))
	fmt.Println("")
	fmt.Println("                    D O N A T I O N     T R A C K E R")
	fmt.Println("")
}
