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

const infinito = math.MaxInt64

func bordoDaDirezione(m Mattoncino, dir Direzione, posizioneBordo string) Bordo {
	if dir == Plus {
		if posizioneBordo == "destra" {
			return m.destra
		}
		return m.sinistra
	}
	if posizioneBordo == "destra" {
		return m.sinistra
	}
	return m.destra
}

func stampaMattoncinoInDirezione(g gioco, sigma string, dir Direzione) {
	if mattoncino, ok := g.mattoncini[sigma]; ok {
		if dir == Plus {
			fmt.Printf("%s: %s, %s\n", sigma, mattoncino.sinistra, mattoncino.destra)
		} else {
			fmt.Printf("%s: %s, %s\n", sigma, mattoncino.destra, mattoncino.sinistra)
		}
	}
}

func inserisciMattoncino(g gioco, alpha, beta, sigma string) {
	g.mattoncini[sigma] = Mattoncino{sinistra: Bordo(alpha), destra: Bordo(beta)}
	g.scatola[sigma] = true
}

func stampaMattoncino(g gioco, sigma string) {
	stampaMattoncinoInDirezione(g, sigma, Plus)
}

func disponiFila(g gioco, listaNomi string) {
	mattonciniOrdinati := strings.Split(listaNomi, " ")

	for _, mattoncinoOrdinato := range mattonciniOrdinati {
		nome := mattoncinoOrdinato[1:]
		if presente := g.scatola[nome]; !presente {
			return
		}
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
		stampaMattoncinoInDirezione(g, string(mattoncino.nome), mattoncino.direzione)
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
	var precedenti map[NomeMattoncino]NomeMattoncino
	distanze := make(map[NomeMattoncino]int)

	// 3. definisco la funzione di calcolo
	calcolaFilaMinima := func(mattoncino NomeMattoncino) (nMattoncini int, ultimoMattoncino NomeMattoncino) {
		for len(coda) > 0 {
			// 3.1 prendo il primo mattoncino dalla coda (coda.pop())
			datiMattoncino := coda[0]
			coda = coda[1:]

			mattoncino := datiMattoncino.mattoncino
			bordoDaControllare := g.mattoncini[mattoncino].destra
			if datiMattoncino.da == bordoDaControllare {
				bordoDaControllare = g.mattoncini[mattoncino].sinistra
			}

			// 3.2 controllo
			if bordoDaControllare == Bordo(beta) {
				return distanze[mattoncino], datiMattoncino.mattoncino
			}

			// 3.3 aggiorno distanze e aggiungo alla coda (coda.enqueue())
			for _, next := range archi[bordoDaControllare] {
				if next == mattoncino || visitati[next] {
					continue
				}

				distanze[next] = distanze[mattoncino] + 1
				visitati[next] = true
				precedenti[next] = mattoncino
				coda = append(coda, MattoncinoDaVisitare{next, bordoDaControllare})
			}
		}

		return infinito, ""
	}

	// 4. calcolo la distanza minima
	distanzaMinima := infinito
	precedentiFilaMinima := precedenti
	var ultimoMattoncinoFilaMinima NomeMattoncino

	for _, mattoncino := range archi[Bordo(alpha)] {
		// 4.1 resetto le strutture dati
		precedenti = make(map[NomeMattoncino]NomeMattoncino)
		visitati = make(map[NomeMattoncino]bool)
		visitati[mattoncino] = true
		for m := range g.scatola {
			distanze[m] = infinito
		}
		distanze[mattoncino] = 0
		coda = []MattoncinoDaVisitare{{mattoncino, Bordo(alpha)}}

		// 4.2 cerco la fila
		distanza, ultimoMattoncino := calcolaFilaMinima(mattoncino)

		// 4.3 controllo la fila trovata
		if distanza < distanzaMinima {
			distanzaMinima = distanza
			precedentiFilaMinima = precedenti
			ultimoMattoncinoFilaMinima = ultimoMattoncino
		}
	}

	// 5. dispongo la fila
	if distanzaMinima == infinito {
		fmt.Printf("non esiste fila da %s a %s\n", alpha, beta)
	} else {
		listaNomi := ""
		mattoncino := ultimoMattoncinoFilaMinima
		bordo := Bordo(beta)

		for {
			direzione := Plus
			if bordo == g.mattoncini[mattoncino].sinistra {
				direzione = Minus
			}

			listaNomi = string(direzione) + mattoncino + listaNomi

			var found bool
			bordo = bordoDaDirezione(g.mattoncini[mattoncino], direzione, "sinistra")
			mattoncino, found = precedentiFilaMinima[mattoncino]

			if !found {
				break
			}
			listaNomi = " " + listaNomi
		}

		disponiFila(g, listaNomi)
	}
}

func sottostringaMassima(g gioco, sigma, tao string) string {
	m := len(sigma)
	n := len(tao)

	if m == 0 || n == 0 {
		return ""
	}

	if m < n {
		sigma, tao = tao, sigma
		m, n = n, m
	}

	current := make([]string, n+1)
	previous := make([]string, n+1)

	for i := 1; i <= m; i++ {
		current, previous = previous, current

		for j := 1; j <= n; j++ {
			if sigma[i-1] == tao[j-1] {
				current[j] = previous[j-1] + string(sigma[i-1])
			} else {
				if len(current[j-1]) > len(previous[j]) {
					current[j] = current[j-1]
				} else {
					current[j] = previous[j]
				}
			}
		}
	}

	return current[n]
}

func indiceCacofonia(g gioco, sigma string) {
	filaDaCalcolare := g.mattoncini[sigma].fila
	if filaDaCalcolare == "" {
		return
	}

	somma := 0
	fila := g.file[filaDaCalcolare]

	for i := 0; i < len(fila)-1; i++ {
		mattoncino := fila[i].nome
		proxMattoncino := fila[i+1].nome

		somma += len(sottostringaMassima(g, mattoncino, proxMattoncino))
	}

	fmt.Println(somma)
}

func costo(g gioco, sigma string, listaBordi ...string) {
	filaDaCalcolare := g.mattoncini[sigma].fila
	if filaDaCalcolare == "" {
		return
	}

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

	for _, mattoncino := range g.file[filaDaCalcolare] {
		sinistra := g.mattoncini[mattoncino.nome].sinistra
		destra := g.mattoncini[mattoncino.nome].destra

		if arco, ok := archi[sinistra]; !ok {
			archi[sinistra] = []NomeMattoncino{mattoncino.nome}
		} else {
			archi[sinistra] = append(arco, mattoncino.nome)
		}

		if arco, ok := archi[destra]; !ok {
			archi[destra] = []NomeMattoncino{mattoncino.nome}
		} else {
			archi[destra] = append(arco, mattoncino.nome)
		}
	}

	// 2. creo le strutture dati
	var mattonciniFila []NomeMattoncino
	var possibiliMattonciniListaBordi [][]NomeMattoncino

	// 2.1 valorizzo lista mattoncini da fila
	for _, mattoncino := range g.file[filaDaCalcolare] {
		mattonciniFila = append(mattonciniFila, mattoncino.nome)
	}

	// 2.2 valorizzo liste mattoncini da lista bordi
	for i := 0; i < len(listaBordi)-1; i++ {
		bordo := Bordo(listaBordi[i])
		proxBordo := Bordo(listaBordi[i+1])

		// 2.2.1 trovo mattoncini candidati per il posto
		var possibiliMattoncini []NomeMattoncino
		for _, mattoncino := range archi[bordo] {
			if (g.mattoncini[mattoncino].sinistra == bordo && g.mattoncini[mattoncino].destra == proxBordo) || (g.mattoncini[mattoncino].destra == bordo && g.mattoncini[mattoncino].sinistra == proxBordo) {
				possibiliMattoncini = append(possibiliMattoncini, mattoncino)
			}
		}

		// 2.2.2 controllo se ho possibili mattoncini
		if len(possibiliMattoncini) == 0 {
			fmt.Println("indefinito 1")
			return
		}

		// 2.2.3 aggiorno le possibili liste
		var newPossibili [][]NomeMattoncino
		if len(possibiliMattonciniListaBordi) == 0 {
			for _, mattoncino := range possibiliMattoncini {
				newPossibili = append(newPossibili, []NomeMattoncino{mattoncino})
			}
		}
		for _, mattonciniListaBordi := range possibiliMattonciniListaBordi {
			var temp []NomeMattoncino

			for _, mattoncino := range possibiliMattoncini {
				// 2.2.4 rimuovo liste non possibili perché contenenti elementi duplicati
				valid := true
				for _, m := range mattonciniListaBordi {
					if m == mattoncino {
						valid = false
						break
					}
				}

				// 2.2.5 se possibile aggiungo alle possibili liste
				if valid {
					temp = append(mattonciniListaBordi, mattoncino)
				}
			}

			// 2.2.6 controllo se ho trovato almeno una lista possibile
			if len(temp) == 0 {
				fmt.Println("indefinito 2")
				return
			}

			newPossibili = append(newPossibili, temp)
		}

		possibiliMattonciniListaBordi = newPossibili
	}

	// 3. definisco la funzione ausiliaria al calcolo
	sottoArrayMassimo := func(a, b []string) int {
		m := len(a)
		n := len(b)

		if m == 0 || n == 0 {
			return 0
		}

		if m < n {
			a, b = b, a
			m, n = n, m
		}

		current := make([]int, n+1)
		previous := make([]int, n+1)

		for i := 1; i <= m; i++ {
			current, previous = previous, current

			for j := 1; j <= n; j++ {
				if a[i-1] == b[j-1] {
					current[j] = previous[j-1] + 1
				} else {
					if current[j-1] > previous[j] {
						current[j] = current[j-1]
					} else {
						current[j] = previous[j]
					}
				}
			}
		}

		return current[n]
	}

	// 4. calcolo il costo
	costoMinimo := infinito
	for _, mattonciniListaBordi := range possibiliMattonciniListaBordi {
		max := sottoArrayMassimo(mattonciniFila, mattonciniListaBordi)
		costo := (len(mattonciniFila) - max) + (len(mattonciniListaBordi) - max)
		if costo < costoMinimo {
			costoMinimo = costo
		}
	}

	// 5. stampo il risultato
	if costoMinimo == infinito {
		fmt.Println("indefinito 3")
	} else {
		fmt.Println(costoMinimo)
	}
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
			parti := strings.Split(input, " ")
			fmt.Println(sottostringaMassima(g, parti[1], parti[2]))
			break

		case 'i':
			parti := strings.Split(input, " ")
			indiceCacofonia(g, parti[1])
			break

		case 'c':
			parti := strings.Split(input, " ")
			costo(g, parti[1], parti[2:]...)
			break

		case 'q':
			return

		default:
			fmt.Println("invalid code")
			break
		}
	}
}
