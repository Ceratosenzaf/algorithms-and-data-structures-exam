// Davide Cerato (matricola 984481)
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type NomeMattoncino = string
type Bordo string

type Mattoncino struct {
	sinistra Bordo    // α
	destra   Bordo    // β
	fila     NomeFila // x = nella fila x, "" = nessuna fila
}

type Direzione string

const (
	Plus  Direzione = "+"
	Minus Direzione = "-"
)

type MattoncinoOrdinato struct {
	nome      NomeMattoncino
	direzione Direzione
}

type NomeFila string           // +a -b ... -z
type Fila []MattoncinoOrdinato // l mattoncini

type gioco struct {
	mattoncini map[NomeMattoncino]Mattoncino // n mattoncini
	scatola    map[NomeMattoncino]bool       // m mattoncini
	file       map[NomeFila]Fila             // k file
}

func inserisciMattoncino(g gioco, alpha, beta, sigma string) {
	g.mattoncini[sigma] = Mattoncino{sinistra: Bordo(alpha), destra: Bordo(beta)}
	g.scatola[sigma] = true
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
		if presente := g.scatola[nome]; !presente {
			return
		}
	}

	bordoDaDirezione := func(b Mattoncino, dir Direzione, posizioneBordo string) Bordo {
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
		dir := Direzione(mattonciniOrdinati[i][0])
		nome := mattonciniOrdinati[i][1:]
		proxDir := Direzione(mattonciniOrdinati[i+1][0])
		proxNome := mattonciniOrdinati[i+1][1:]

		if bordoDaDirezione(g.mattoncini[nome], dir, "destra") != bordoDaDirezione(g.mattoncini[proxNome], proxDir, "sinistra") {
			return
		}
	}

	nomeFila := NomeFila(listaNomi)
	fila := make(Fila, len(mattonciniOrdinati))
	for i, mattoncinoOrdinato := range mattonciniOrdinati {
		direzione := Direzione(mattoncinoOrdinato[0])
		nome := mattoncinoOrdinato[1:]

		delete(g.scatola, nome)
		fila[i] = MattoncinoOrdinato{nome, direzione}
		g.mattoncini[nome] = Mattoncino{sinistra: g.mattoncini[nome].sinistra, destra: g.mattoncini[nome].destra, fila: nomeFila}
	}

	g.file[nomeFila] = fila
}

func stampaFila(g gioco, sigma string) {
	if m, ok := g.mattoncini[sigma]; !ok || m.fila == "" {
		return
	}

	fmt.Println("(")
	for _, mattoncino := range g.file[g.mattoncini[sigma].fila] {
		stampaMattoncino(g, string(mattoncino.nome))
	}
	fmt.Println(")")
}

func eliminaFila(g gioco, sigma string) {
	if m, ok := g.mattoncini[sigma]; !ok || m.fila == "" {
		return
	}

	filaDaEliminare := g.mattoncini[sigma].fila
	for _, elementoFila := range g.file[filaDaEliminare] {
		g.scatola[elementoFila.nome] = true
		g.mattoncini[elementoFila.nome] = Mattoncino{sinistra: g.mattoncini[elementoFila.nome].sinistra, destra: g.mattoncini[elementoFila.nome].destra}
	}
	delete(g.file, filaDaEliminare)
}

func disponiFilaMinima(g gioco, alpha, beta string) {
	const infinito = math.MaxInt64

	// 1. creo lista di adiacenza
	archi := make(map[Bordo][]NomeMattoncino)

	for mattoncino := range g.scatola {
		sinistra := g.mattoncini[mattoncino].sinistra
		destra := g.mattoncini[mattoncino].destra

		if arco, ok := archi[sinistra]; !ok {
			archi[sinistra] = []NomeMattoncino{mattoncino}
		} else {
			archi[sinistra] = append(arco, mattoncino)
		}

		if arco, ok := archi[destra]; !ok {
			archi[destra] = []NomeMattoncino{mattoncino}
		} else {
			archi[destra] = append(arco, mattoncino)
		}
	}

	// 2. inizializzo le strutture dati
	type MattoncinoDaVisitare struct {
		mattoncino NomeMattoncino
		da         Bordo
	}
	var coda []MattoncinoDaVisitare
	var visitati map[NomeMattoncino]bool
	distanze := make(map[NomeMattoncino]int)

	// 3. definisco la funzione di calcolo
	calcolaFilaMinima := func(mattoncino NomeMattoncino) (nMattoncini int) {
		for len(coda) > 0 {
			// 3.1 prendo il primo mattoncino dalla coda (coda.pop())
			datiMattoncino := coda[0]
			coda = coda[1:]

			mattoncino := datiMattoncino.mattoncino
			bordoDaControllare := g.mattoncini[mattoncino].destra
			if datiMattoncino.da == bordoDaControllare {
				bordoDaControllare = g.mattoncini[mattoncino].sinistra
			}

			fmt.Printf("Controllando %s da %s a %s\n", mattoncino, datiMattoncino.da, bordoDaControllare)

			// 3.2 controllo
			if bordoDaControllare == Bordo(beta) {
				return distanze[mattoncino]
			}

			// 3.3 aggiorno distanze e aggiungo alla coda (coda.enqueue())
			for _, next := range archi[bordoDaControllare] {
				if next == mattoncino || visitati[next] {
					continue
				}

				distanze[next] = distanze[mattoncino] + 1
				visitati[next] = true
				coda = append(coda, MattoncinoDaVisitare{next, bordoDaControllare})
			}
		}

		return infinito
	}

	// 4. calcolo la distanza minima
	distanzaMinima := infinito

	for _, mattoncino := range archi[Bordo(alpha)] {
		// 4.1 resetto le strutture dati
		visitati = make(map[NomeMattoncino]bool)
		visitati[mattoncino] = true
		for m := range g.scatola {
			distanze[m] = infinito
		}
		distanze[mattoncino] = 0
		coda = []MattoncinoDaVisitare{{mattoncino, Bordo(alpha)}}

		// 4.2 cerco la fila
		distanza := calcolaFilaMinima(mattoncino)

		// 4.3 controllo la fila trovata
		if distanza < distanzaMinima {
			distanzaMinima = distanza
		}
	}

	// 5. dispongo la fila
	if distanzaMinima == infinito {
		fmt.Printf("non esiste fila da %s a %s\n", alpha, beta)
	} else {
		fmt.Println("distanza minima", distanzaMinima+1)
	}

}

func sottostringaMassima(g gioco, sigma, tao Fila) string {
	return ""
}

func indiceCacofonia(g gioco, sigma string) {

}

func costo(g gioco, sigma Fila, listaNomi string) {

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	g := gioco{mattoncini: make(map[NomeMattoncino]Mattoncino), scatola: make(map[NomeMattoncino]bool), file: make(map[NomeFila]Fila)}

	for scanner.Scan() {
		input := scanner.Text()

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
			parti := strings.Split(input, " ")
			disponiFilaMinima(g, parti[1], parti[2])
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

		fmt.Println(g.file)
		fmt.Println(g.scatola)
		fmt.Println(g.mattoncini)
	}
}
