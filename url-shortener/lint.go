package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	if _, err := exec.LookPath("fieldalignment"); err != nil {
		log.Println("fieldalignment не найден. Устанавливаем...")

		if err := installFieldalignment(); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := exec.LookPath("gci"); err != nil {
		log.Println("gci не найден. Устанавливаем...")

		if err := installGCI(); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := exec.LookPath("golangci-lint"); err != nil {
		log.Println("golangci-lint не найден. Устанавливаем...")

		if err := installGolangciLint(); err != nil {
			log.Fatal(err)
		}
	}

	if err := runGCI(); err != nil {
		log.Fatalf("Ошибка запуска gci: %v", err)
	}

	log.Println("Изменения gci внесены успешно!")

	if err := runGolangciLint(); err != nil {
		log.Fatalf("Ошибка запуска golangci-lint: %v", err)
	}

	log.Println("Проверка golangci-lint завершена успешно!")

	if err := runFieldalignment(); err != nil {
		log.Printf("%v Fieldalignment сейчас всё сам исправит.", err)
	}

	log.Println("Проверка fieldalignment завершена успешно!")
}

func installGCI() error {
	cmd := exec.Command("go", "install", "github.com/daixiang0/gci@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runGCI() error {
	cmd := exec.Command("sh", "-c", "find . -name '*.go' -not -path './vendor/*' -not -path './gen/*' | xargs gci write")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func installGolangciLint() error {
	cmd := exec.Command("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runGolangciLint() error {
	cmd := exec.Command("golangci-lint", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func installFieldalignment() error {
	cmd := exec.Command("go", "install", "golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runFieldalignment() error {
	cmd := exec.Command("fieldalignment", "-fix", "./")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
