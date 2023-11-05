package main

import (
	"database/sql"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/twsm000/tt/internal/entities/translator"
	"github.com/twsm000/tt/internal/repositories/database/sql/sqlite3"
	"github.com/twsm000/tt/internal/services"
	"github.com/twsm000/tt/internal/utils"
)

//go:embed trans_cli
var transCLI []byte

func main() {
	dbFlag := flag.NewFlagSet("db path", flag.ContinueOnError)
	dbPath := dbFlag.String("dbp", "tt.db", "dbp indica o caminho do banco de dados. (Default ~/tt.db)")
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "Modo de uso: tt [commands] [flags]")
		fmt.Fprintln(w, "Commands:")
		fmt.Fprintln(w, "init    cria o arquivo trans cli para tradução")
		fmt.Fprintln(w, "word    traduz uma palavra")
		fmt.Fprintln(w, "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		InitializeTranslator()
	case "word":
		if err := dbFlag.Parse(os.Args[2:]); err != nil {
			flag.PrintDefaults()
			os.Exit(1)
		}
		Translate(*dbPath, strings.Join(dbFlag.Args(), " "))
	default:
		fmt.Fprintln(os.Stderr, "Opção inválida.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func InitializeTranslator() {
	execPath := utils.MustGet(os.Executable())
	dir := filepath.Dir(execPath)
	transPath := filepath.Join(dir, "trans")
	log.Println("Criando tradutor em:", transPath)
	utils.TryTerminate(os.WriteFile(transPath, transCLI, 0664))
	cmd := exec.Command("chmod", "+x", transPath)
	utils.TryTerminate(cmd.Run())
	return
}

func Translate(dbPath, word string) {
	if word == "" {
		fmt.Fprintln(os.Stderr, "Nenhuma palavra foi fornecida.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Println("Abrindo banco sqlite3:", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Falha ao abrir conexão com banco de dados:", err)
		os.Exit(1)
	}
	if err := db.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, "Falha ao conectar ao banco de dados:", err)
		os.Exit(1)
	}
	defer db.Close()

	translationRepo, err := sqlite3.NewTranslationRepository(db)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Falha ao criar repositório:", err)
		os.Exit(1)
	}

	translator := services.NewTranslationService(translator.TransCLI{}, translationRepo)
	translation, err := translator.Translate(word)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Falha ao traduzir %q: %v", word, err)
		os.Exit(1)
	}
	fmt.Println(translation)
}
