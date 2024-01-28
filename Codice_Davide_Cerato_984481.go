// Davide Cerato (matricola 984481)
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type nomeMattoncino = string
type bordo string
type bordi struct {
	sinistra bordo // α
	destra   bordo // β
}

type direzione string

const (
	Plus  direzione = "+"
	Minus direzione = "-"
)

type elementoFila struct { // mattoncino ordinato
	nome      nomeMattoncino
	direzione direzione
}

type nomeFila string       // "first.dir first.nome ... last.dir last.nome"
type fila = []elementoFila // k mattoncini

type gioco struct {
	mattoncini map[nomeMattoncino]bordi    // n  mattoncini
	scatola    map[nomeMattoncino]nomeFila // l mattoncini, "": nella scatola, x != "": nella fila x
	file       map[nomeFila]fila           // m file
}

func inserisciMattoncino(g gioco, alpha, beta, sigma string) {
	g.mattoncini[sigma] = bordi{sinistra: bordo(alpha), destra: bordo(beta)}
	g.scatola[sigma] = ""
}

func stampaMattoncino(g gioco, sigma string) {
	if m, ok := g.mattoncini[sigma]; ok {
		fmt.Printf("%s: %s, %s\n", sigma, m.sinistra, m.destra)
	}
}

func disponiFila(g gioco, listaNomi string) {
	mattonciniOrdinati := strings.Split(listaNomi, " ")

	for _, mattoncinoOrdinato := range mattonciniOrdinati {
		nome := mattoncinoOrdinato[1:]
		if fila, ok := g.scatola[nome]; !ok || fila != "" {
			return
		}
	}

	getBordoDaDirezione := func(b bordi, dir direzione, posizioneBordo string) bordo {
		if dir == Plus {
			if posizioneBordo == "destra" {
				return b.destra
			}
			return b.sinistra
		}
		if posizioneBordo == "destra" {
			return b.sinistra
		}
		return b.destra
	}

	for i := 0; i < len(mattonciniOrdinati)-1; i++ {
		dir := direzione(mattonciniOrdinati[i][0])
		nomeMattoncino := mattonciniOrdinati[i][1:]
		proxDir := direzione(mattonciniOrdinati[i+1][0])
		proxNomeMattoncino := mattonciniOrdinati[i+1][1:]

		if getBordoDaDirezione(g.mattoncini[nomeMattoncino], dir, "destra") != getBordoDaDirezione(g.mattoncini[proxNomeMattoncino], proxDir, "sinistra") {
			return
		}
	}

	f := make([]elementoFila, len(mattonciniOrdinati))
	for i, mattoncinoOrdinato := range mattonciniOrdinati {
		direzione := direzione(mattoncinoOrdinato[0])
		nome := mattoncinoOrdinato[1:]
		f[i] = elementoFila{nome, direzione}
		g.scatola[nome] = nomeFila(listaNomi)
	}

	g.file[nomeFila(listaNomi)] = f
}

func stampaFila(g gioco, sigma string) {
	if fila, ok := g.scatola[sigma]; !ok || fila == "" {
		return
	}

	fmt.Println("(")
	for _, elementoFila := range g.file[g.scatola[sigma]] {
		stampaMattoncino(g, elementoFila.nome)
	}
	fmt.Println(")")
}

func eliminaFila(g gioco, sigma string) {
	if fila, ok := g.scatola[sigma]; !ok || fila == "" {
		return
	}

	filaDaEliminare := g.scatola[sigma]
	for _, elementoFila := range g.file[filaDaEliminare] {
		g.scatola[elementoFila.nome] = ""
	}
	delete(g.file, filaDaEliminare)
}

func disponiFilaMinima(g gioco, alpha, beta string) {

}

func sottostringaMassima(g gioco, sigma, tao fila) {

}

func indiceCacofonia(g gioco, sigma string) {

}

func costo(g gioco, sigma fila, listaNomi string) {

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	g := gioco{mattoncini: make(map[string]bordi), scatola: make(map[string]nomeFila), file: make(map[nomeFila][]elementoFila)}

	for scanner.Scan() {
		input := scanner.Text()
		fmt.Println("pre", g)

		switch input[0] {
		case 'm':
			parti := strings.Split(input, " ")
			inserisciMattoncino(g, parti[1], parti[2], parti[3])
			break

		case 's':
			parti := strings.Split(input, " ")
			stampaMattoncino(g, parti[1])
			break

		case 'd':
			parti := strings.SplitN(input, " ", 2)
			disponiFila(g, parti[1])
			break

		case 'S':
			parti := strings.Split(input, " ")
			stampaFila(g, parti[1])
			break

		case 'e':
			parti := strings.Split(input, " ")
			eliminaFila(g, parti[1])
			break

		case 'f':
			fmt.Println("disponiFilaMinima()")
			break

		case 'M':
			fmt.Println("sottostringaMassima()")
			break

		case 'i':
			fmt.Println("indiceCacofonia()")
			break

		case 'c':
			fmt.Println("costo()")
			break

		case 'q':
			return

		default:
			fmt.Println("invalid code")
			break
		}
		fmt.Println("post", g)
	}
}
