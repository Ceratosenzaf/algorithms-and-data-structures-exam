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
	fila     NomeFila // x = nella fila x, "" = nella scatola
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

type NomeFila string // +a -b ... -z
type Fila []MattoncinoOrdinato

type gioco struct {
	mattoncini map[NomeMattoncino]Mattoncino
	file       map[NomeFila]Fila
	scatola    map[Bordo]map[Bordo]int
}

const INFINITO = math.MaxInt64

// funzioni per uso interno
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

func riduciORimuovi(target map[Bordo]map[Bordo]int, a, b Bordo) {
	if target[a][b] > 1 {
		target[a][b] = target[a][b] - 1
	} else {
		delete(target[a], b)
	}
	if target[b][a] > 1 {
		target[b][a] = target[b][a] - 1
	} else {
		delete(target[b], a)
	}

	if len(target[a]) == 0 {
		delete(target, a)
	}
	if len(target[b]) == 0 {
		delete(target, b)
	}
}

func aggiungiMattoncinoAScatola(g gioco, m Mattoncino) {
	if _, ok := g.scatola[m.sinistra]; !ok {
		g.scatola[m.sinistra] = make(map[Bordo]int)
	}
	if _, ok := g.scatola[m.destra]; !ok {
		g.scatola[m.destra] = make(map[Bordo]int)
	}
	g.scatola[m.sinistra][m.destra] = g.scatola[m.sinistra][m.destra] + 1
	g.scatola[m.destra][m.sinistra] = g.scatola[m.destra][m.sinistra] + 1
}

func rimuoviMattoncinoDaScatola(g gioco, m NomeMattoncino) {
	riduciORimuovi(g.scatola, g.mattoncini[m].sinistra, g.mattoncini[m].destra)
}

func copiaScatola(g gioco, target map[Bordo]map[Bordo]int) {
	for b, bordiAdiacenti := range g.scatola {
		if _, ok := target[b]; !ok {
			target[b] = make(map[Bordo]int)
		}
		for v, d := range bordiAdiacenti {
			target[b][v] = d
		}
	}
}

// funzioni del progetto
func listaNomiDaListaBordi(g gioco, listaBordi []Bordo) (listaNomi string) {
	archi := make(map[Bordo]map[Bordo]map[NomeMattoncino]bool)
	for nomeMattoncino, mattoncino := range g.mattoncini {
		if mattoncino.fila == "" {
			if _, ok := archi[mattoncino.sinistra]; !ok {
				archi[mattoncino.sinistra] = make(map[Bordo]map[NomeMattoncino]bool)
			}
			if _, ok := archi[mattoncino.sinistra][mattoncino.destra]; !ok {
				archi[mattoncino.sinistra][mattoncino.destra] = make(map[NomeMattoncino]bool)
			}
			archi[mattoncino.sinistra][mattoncino.destra][nomeMattoncino] = true
		}
	}

	for i := 0; i < len(listaBordi)-1; i++ {
		bordo := listaBordi[i]
		proxBordo := listaBordi[i+1]

		var dir Direzione
		var mattoncino NomeMattoncino

		if mattonciniDritti := archi[bordo][proxBordo]; len(mattonciniDritti) > 0 {
			dir = Plus
			for m := range mattonciniDritti {
				mattoncino = m
				break
			}
			delete(mattonciniDritti, mattoncino)
		} else if mattonciniRovesci := archi[proxBordo][bordo]; len(mattonciniRovesci) > 0 {
			dir = Minus
			for m := range mattonciniRovesci {
				mattoncino = m
				break
			}
			delete(mattonciniRovesci, mattoncino)
		}

		if mattoncino == "" {
			return ""
		}

		listaNomi += " " + string(dir) + mattoncino
	}

	return listaNomi[1:]
}

func aggiungiAListaDiAdiacenzaDaNome(g gioco, archi map[Bordo][]NomeMattoncino, sigma string) {
	sinistra := g.mattoncini[sigma].sinistra
	destra := g.mattoncini[sigma].destra

	if arco, ok := archi[sinistra]; !ok {
		archi[sinistra] = []NomeMattoncino{sigma}
	} else {
		archi[sinistra] = append(arco, sigma)
	}

	if arco, ok := archi[destra]; !ok {
		archi[destra] = []NomeMattoncino{sigma}
	} else {
		archi[destra] = append(arco, sigma)
	}
}

func inserisciMattoncino(g gioco, alpha, beta, sigma string) {
	g.mattoncini[sigma] = Mattoncino{sinistra: Bordo(alpha), destra: Bordo(beta)}
	aggiungiMattoncinoAScatola(g, g.mattoncini[sigma])
}

func stampaMattoncino(g gioco, sigma string) {
	stampaMattoncinoInDirezione(g, sigma, Plus)
}

func disponiFila(g gioco, listaNomi string) {
	mattonciniOrdinati := strings.Split(listaNomi, " ")

	visti := make(map[NomeMattoncino]bool)
	for _, mattoncinoOrdinato := range mattonciniOrdinati {
		nome := mattoncinoOrdinato[1:]
		if m := g.mattoncini[nome]; m.fila != "" {
			return
		}
		if duplicato := visti[nome]; duplicato {
			return
		}
		visti[nome] = true
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

		fila[i] = MattoncinoOrdinato{nome, direzione}
		g.mattoncini[nome] = Mattoncino{sinistra: g.mattoncini[nome].sinistra, destra: g.mattoncini[nome].destra, fila: nomeFila}
		rimuoviMattoncinoDaScatola(g, nome)
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
		g.mattoncini[elementoFila.nome] = Mattoncino{
			sinistra: g.mattoncini[elementoFila.nome].sinistra,
			destra:   g.mattoncini[elementoFila.nome].destra,
		}
		aggiungiMattoncinoAScatola(g, g.mattoncini[elementoFila.nome])
	}
	delete(g.file, filaDaEliminare)
}

func disponiFilaMinima(g gioco, alpha, beta string) {
	precedenti := make(map[Bordo]Bordo)

	// 1. funzione di calcolo distanza minima
	trovaMinimo := func() int {
		// bordi non presenti
		if _, ok := g.scatola[Bordo(alpha)]; !ok {
			return INFINITO
		}
		if _, ok := g.scatola[Bordo(beta)]; !ok {
			return INFINITO
		}

		// mattoncino con bordi alpha e beta
		if g.scatola[Bordo(alpha)][Bordo(beta)] > 0 {
			precedenti[Bordo(beta)] = Bordo(alpha)
			return 1
		}

		visitati := make(map[Bordo]bool)
		distanze := make(map[Bordo]int)
		archiDisponibili := make(map[Bordo]map[Bordo]int, len(g.scatola))
		copiaScatola(g, archiDisponibili)

		bfs := func(partenza, arrivo Bordo) (costo int) {
			coda := []Bordo{partenza}

			for len(coda) > 0 {
				bordo := coda[0]
				coda = coda[1:]

				for proxBordo, d := range archiDisponibili[bordo] {
					if !visitati[proxBordo] && d > 0 {
						riduciORimuovi(archiDisponibili, bordo, proxBordo)
						visitati[proxBordo] = true

						distanze[proxBordo] = distanze[bordo] + 1
						precedenti[proxBordo] = bordo

						if proxBordo == arrivo {
							return distanze[proxBordo]
						}

						coda = append(coda, proxBordo)
					}
				}
			}

			return INFINITO
		}

		// mattoncini diversi con alpha e beta diversi
		if alpha != beta {
			for b := range g.scatola {
				distanze[b] = INFINITO
			}
			distanze[Bordo(alpha)] = 0
			return bfs(Bordo(alpha), Bordo(beta))
		}

		// mattoncini diversi con alpha e beta uguali
		distanzaMinima := INFINITO
		precedentiDistanzaMinima := make(map[Bordo]Bordo)
		for bordoAdiacenteAlpha := range g.scatola[Bordo(alpha)] {
			// reset strutture dati
			for b := range g.scatola {
				visitati[b] = false
				distanze[b] = INFINITO
				precedenti[b] = ""
			}
			visitati[bordoAdiacenteAlpha] = true
			distanze[bordoAdiacenteAlpha] = 1
			precedenti[bordoAdiacenteAlpha] = Bordo(alpha)
			copiaScatola(g, archiDisponibili)
			riduciORimuovi(archiDisponibili, bordoAdiacenteAlpha, Bordo(alpha))

			// calcolo
			distanza := bfs(bordoAdiacenteAlpha, Bordo(beta))

			// confronto
			if distanza < distanzaMinima {
				distanzaMinima = distanza
				for k, v := range precedenti {
					precedentiDistanzaMinima[k] = v
				}
				precedentiDistanzaMinima[bordoAdiacenteAlpha] = Bordo(alpha)
			}
		}
		precedenti = precedentiDistanzaMinima
		return distanzaMinima
	}

	// 2. calcolo la distanza minima
	if distanzaMinima := trovaMinimo(); distanzaMinima == INFINITO {
		fmt.Printf("non esiste fila da %s a %s\n", alpha, beta)
	} else {
		// 3. dispongo la fila
		listaBordi := make([]Bordo, distanzaMinima+1)

		bordo := Bordo(beta)
		for i := distanzaMinima; i >= 0; i-- {
			listaBordi[i] = bordo
			bordo = precedenti[bordo]
		}

		listaNomi := listaNomiDaListaBordi(g, listaBordi)
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
	for i := 0; i < len(g.file[filaDaCalcolare])-1; i++ {
		mattoncino := g.file[filaDaCalcolare][i].nome
		proxMattoncino := g.file[filaDaCalcolare][i+1].nome

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
	for nomeMattoncino, mattoncino := range g.mattoncini {
		if mattoncino.fila == "" {
			aggiungiAListaDiAdiacenzaDaNome(g, archi, nomeMattoncino)
		}
	}
	for _, mattoncino := range g.file[filaDaCalcolare] {
		aggiungiAListaDiAdiacenzaDaNome(g, archi, mattoncino.nome)
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
			fmt.Println("indefinito")
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
				fmt.Println("indefinito")
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
	costoMinimo := INFINITO
	for _, mattonciniListaBordi := range possibiliMattonciniListaBordi {
		max := sottoArrayMassimo(mattonciniFila, mattonciniListaBordi)
		costo := (len(mattonciniFila) - max) + (len(mattonciniListaBordi) - max)
		if costo < costoMinimo {
			costoMinimo = costo
		}
	}

	// 5. stampo il risultato
	if costoMinimo == INFINITO {
		fmt.Println("indefinito")
	} else {
		fmt.Println(costoMinimo)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	g := gioco{mattoncini: make(map[NomeMattoncino]Mattoncino), file: make(map[NomeFila]Fila), scatola: make(map[Bordo]map[Bordo]int)}

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
