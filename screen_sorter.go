package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func get_month_year(processed_screenshot_name string) (time.Month, int, error) {
	// split_name[0] = name, [1] = date, [2] = time
	split_name := strings.Split(processed_screenshot_name, " ")
	date, err := time.Parse("2006-01-02", split_name[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Skipping %s: failed to get date", processed_screenshot_name)
	}

	present_year := time.Now().Year()
	if date.Year() > present_year || date.Year() < 2000 {
		return 0, 0, fmt.Errorf("Skipping %s: failed to get date", processed_screenshot_name)
	}

	return date.Month(), date.Year(), nil
}

func create_folder_by_year_month(screenshot_path string, month time.Month, year int) {
	year_path := filepath.Join(screenshot_path, strconv.Itoa(year))
	err := os.Mkdir(year_path, 0o700)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal("Error creating year folder")
		}
	}

	month_path := filepath.Join(year_path, month.String())
	err = os.Mkdir(month_path, 0o700)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal("Error creating month folder")
		}
	}
}

func move_screenshot(screenshot_path string, month time.Month, year int, output_name string) error {
	output_path_current := filepath.Join(screenshot_path, output_name)
	output_path_moved := filepath.Join(screenshot_path, strconv.Itoa(year), month.String(), output_name)

	_, err := os.Stat(output_path_moved)
	if err == nil {
		return fmt.Errorf("File %s already exists, skipping it", output_name)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		log.Fatal(err)
	}

	err = os.Rename(output_path_current, output_path_moved)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func validate_name(output os.DirEntry) error {
	if output.IsDir() {
		return fmt.Errorf("skipping %s since its a folder", output.Name())
	}
	if !strings.HasSuffix(output.Name(), ".png") {
		return fmt.Errorf("Skipping %s since its not a png", output.Name())
	} // there should be 2 spcaes in the screenshots made by snipping tool
	if strings.Count(output.Name(), " ") != 2 {
		return fmt.Errorf("Skipping %s since its not the right prefix", output.Name())
	}
	return nil
}

func main() {
	if runtime.GOOS != "windows" {
		log.Fatal("This tool is only supported for Windows.")
	}
	user_dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	screenshot_path := filepath.Join(user_dir, "Pictures", "Screenshots")
	file_container, err := os.ReadDir(screenshot_path)
	if err != nil {
		log.Fatal(err)
	}

	var counter_moved int = 0
	for _, output := range file_container {
		err := validate_name(output)
		if err != nil {
			fmt.Println(err)
			continue
		}
		processed_screenshot_name, _ := strings.CutSuffix(output.Name(), ".png")
		month, year, err := get_month_year(processed_screenshot_name)
		if err != nil {
			fmt.Println(err)
			continue
		}
		create_folder_by_year_month(screenshot_path, month, year)
		err = move_screenshot(screenshot_path, month, year, output.Name())
		if err != nil {
			fmt.Println(err)
			continue
		}
		counter_moved++
	}

	fmt.Printf("moved %d screenshots\n", counter_moved)
}
