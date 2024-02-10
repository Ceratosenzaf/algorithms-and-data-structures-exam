// Davide Cerato (matricola 984481)

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

type NomeMattoncino = string // σ
type Bordo string

type mattoncino struct {
	sinistra Bordo    // α
	destra   Bordo    // β
	fila     NomeFila // x = nella fila x, "" = nella scatola
}

type Mattoncini map[NomeMattoncino]mattoncino

type Direzione string

const (
	Plus  Direzione = "+"
	Minus Direzione = "-"
)

type MattoncinoOrdinato struct {
	nome      NomeMattoncino
	direzione Direzione
}

type NomeFila string // +σ1 -σ2 ... -σn
type fila []MattoncinoOrdinato

type File map[NomeFila]fila

type SetMattoncini map[NomeMattoncino]bool
type Scatola map[Bordo]map[Bordo]SetMattoncini
type ScatolaSenzaNomi map[Bordo]map[Bordo]int

type gioco struct {
	mattoncini Mattoncini
	file       File
	scatola    Scatola
}

const INFINITO = math.MaxInt64

// funzioni per uso interno
func bordoDaDirezione(m mattoncino, dir Direzione, posizioneBordo string) Bordo {
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
	if m, ok := g.mattoncini[sigma]; ok {
		if dir == Plus {
			fmt.Printf("%s: %s, %s\n", sigma, m.sinistra, m.destra)
		} else {
			fmt.Printf("%s: %s, %s\n", sigma, m.destra, m.sinistra)
		}
	}
}

func riduciORimuovi(target ScatolaSenzaNomi, a, b Bordo) {
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

func aggiungiAListaAdiacenza(g gioco, target Scatola, m NomeMattoncino) {
	a := g.mattoncini[m].sinistra
	b := g.mattoncini[m].destra

	if _, ok := target[a]; !ok {
		target[a] = make(map[Bordo]SetMattoncini)
	}
	if _, ok := target[b]; !ok {
		target[b] = make(map[Bordo]SetMattoncini)
	}

	if _, ok := target[a][b]; !ok {
		target[a][b] = make(SetMattoncini)
	}
	if _, ok := target[b][a]; !ok {
		target[b][a] = make(SetMattoncini)
	}

	target[a][b][m] = true
	target[b][a][m] = true
}

func aggiungiMattoncinoAScatola(g gioco, m NomeMattoncino) {
	aggiungiAListaAdiacenza(g, g.scatola, m)
}

func rimuoviMattoncinoDaScatola(g gioco, m NomeMattoncino) {
	a := g.mattoncini[m].sinistra
	b := g.mattoncini[m].destra

	delete(g.scatola[a][b], m)
	delete(g.scatola[b][a], m)

	if len(g.scatola[a][b]) == 0 {
		delete(g.scatola[a], b)
	}
	if len(g.scatola[b][a]) == 0 {
		delete(g.scatola[b], a)
	}

	if len(g.scatola[a]) == 0 {
		delete(g.scatola, a)
	}
	if len(g.scatola[b]) == 0 {
		delete(g.scatola, b)
	}
}

func trascriviScatolaSenzaNomi(g gioco, target map[Bordo]map[Bordo]int) {
	for b, bordiAdiacenti := range g.scatola {
		if _, ok := target[b]; !ok {
			target[b] = make(map[Bordo]int)
		}
		for v, d := range bordiAdiacenti {
			target[b][v] = len(d)
		}
	}
}

// funzioni del progetto
func listaNomiDaListaBordi(g gioco, listaBordi []Bordo) (listaNomi string) {
	usati := make(SetMattoncini)

	for i := 0; i < len(listaBordi)-1; i++ {
		bordo := listaBordi[i]
		proxBordo := listaBordi[i+1]

		var nomeMattoncino NomeMattoncino
		dir := Plus

		if mattoncini := g.scatola[bordo][proxBordo]; len(mattoncini) > 0 {
			for m := range mattoncini {
				if !usati[m] {
					nomeMattoncino = m
					break
				}
			}
			usati[nomeMattoncino] = true
			if g.mattoncini[nomeMattoncino].sinistra != bordo {
				dir = Minus
			}
		}

		if nomeMattoncino == "" {
			return ""
		}

		listaNomi += " " + string(dir) + nomeMattoncino
	}

	return listaNomi[1:]
}

func inserisciMattoncino(g gioco, alpha, beta, sigma string) {
	g.mattoncini[sigma] = mattoncino{sinistra: Bordo(alpha), destra: Bordo(beta)}
	aggiungiMattoncinoAScatola(g, sigma)
}

func stampaMattoncino(g gioco, sigma string) {
	stampaMattoncinoInDirezione(g, sigma, Plus)
}

func disponiFila(g gioco, listaNomi string) {
	mattonciniOrdinati := strings.Split(listaNomi, " ")

	visti := make(SetMattoncini)
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
	nuovaFila := make(fila, len(mattonciniOrdinati))
	for i, mattoncinoOrdinato := range mattonciniOrdinati {
		direzione := Direzione(mattoncinoOrdinato[0])
		nome := mattoncinoOrdinato[1:]

		nuovaFila[i] = MattoncinoOrdinato{nome, direzione}
		g.mattoncini[nome] = mattoncino{sinistra: g.mattoncini[nome].sinistra, destra: g.mattoncini[nome].destra, fila: nomeFila}
		rimuoviMattoncinoDaScatola(g, nome)
	}

	g.file[nomeFila] = nuovaFila
}

func stampaFila(g gioco, sigma string) {
	if m, ok := g.mattoncini[sigma]; !ok || m.fila == "" {
		return
	}

	fmt.Println("(")
	for _, m := range g.file[g.mattoncini[sigma].fila] {
		stampaMattoncinoInDirezione(g, string(m.nome), m.direzione)
	}
	fmt.Println(")")
}

func eliminaFila(g gioco, sigma string) {
	if m, ok := g.mattoncini[sigma]; !ok || m.fila == "" {
		return
	}

	filaDaEliminare := g.mattoncini[sigma].fila
	for _, elementoFila := range g.file[filaDaEliminare] {
		g.mattoncini[elementoFila.nome] = mattoncino{
			sinistra: g.mattoncini[elementoFila.nome].sinistra,
			destra:   g.mattoncini[elementoFila.nome].destra,
		}
		aggiungiMattoncinoAScatola(g, elementoFila.nome)
	}
	delete(g.file, filaDaEliminare)
}

func disponiFilaMinima(g gioco, alpha, beta string) {
	var precedenti map[Bordo]Bordo

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
		if len(g.scatola[Bordo(alpha)][Bordo(beta)]) > 0 {
			precedenti = make(map[Bordo]Bordo, 1)
			precedenti[Bordo(beta)] = Bordo(alpha)
			return 1
		}

		precedenti = make(map[Bordo]Bordo)
		visitati := make(map[Bordo]bool)
		distanze := make(map[Bordo]int)
		archiDisponibili := make(ScatolaSenzaNomi, len(g.scatola))
		trascriviScatolaSenzaNomi(g, archiDisponibili)

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
			trascriviScatolaSenzaNomi(g, archiDisponibili)
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

	stringheComuni := make([]string, n+1)
	prevStringheComuni := make([]string, n+1)

	for i := 1; i <= m; i++ {
		stringheComuni, prevStringheComuni = prevStringheComuni, stringheComuni
		for j := 1; j <= n; j++ {
			if sigma[i-1] == tao[j-1] {
				stringheComuni[j] = prevStringheComuni[j-1] + string(sigma[i-1])
			} else {
				if len(stringheComuni[j-1]) > len(prevStringheComuni[j]) {
					stringheComuni[j] = stringheComuni[j-1]
				} else {
					stringheComuni[j] = prevStringheComuni[j]
				}
			}
		}
	}

	return stringheComuni[n]
}

func indiceCacofonia(g gioco, sigma string) {
	filaDaCalcolare := g.mattoncini[sigma].fila
	if filaDaCalcolare == "" {
		return
	}

	somma := 0
	for i := 0; i < len(g.file[filaDaCalcolare])-1; i++ {
		m := g.file[filaDaCalcolare][i].nome
		prox := g.file[filaDaCalcolare][i+1].nome

		somma += len(sottostringaMassima(g, m, prox))
	}

	fmt.Println(somma)
}

func costo(g gioco, sigma string, listaBordi ...string) {
	filaDaCalcolare := g.mattoncini[sigma].fila
	if filaDaCalcolare == "" {
		return
	}

	// 1. creo lista di adiacenza della fila
	archiFila := make(Scatola)
	for _, m := range g.file[filaDaCalcolare] {
		aggiungiAListaAdiacenza(g, archiFila, m.nome)
	}

	// 2. creo le strutture dati
	var mattonciniFila []NomeMattoncino
	var possibiliMattonciniPerPosizione [][]NomeMattoncino

	// 2.1 valorizzo lista mattoncini da fila
	for _, m := range g.file[filaDaCalcolare] {
		mattonciniFila = append(mattonciniFila, m.nome)
	}

	// 2.2 valorizzo liste mattoncini da lista bordi
	for i := 0; i < len(listaBordi)-1; i++ {
		bordo := Bordo(listaBordi[i])
		proxBordo := Bordo(listaBordi[i+1])

		// 2.2.1 trovo mattoncini candidati per il posto
		var possibiliMattoncini []NomeMattoncino
		for m := range g.scatola[bordo][proxBordo] {
			possibiliMattoncini = append(possibiliMattoncini, m)
		}
		for m := range archiFila[bordo][proxBordo] {
			possibiliMattoncini = append(possibiliMattoncini, m)
		}

		// 2.2.2 controllo se ho possibili mattoncini e li setto
		if len(possibiliMattoncini) == 0 {
			fmt.Println("indefinito")
			return
		}

		possibiliMattonciniPerPosizione = append(possibiliMattonciniPerPosizione, possibiliMattoncini)
	}

	// 2.3 creo le file possibili
	var possibiliFile [][]NomeMattoncino
	for _, possibiliMattoncini := range possibiliMattonciniPerPosizione {
		if len(possibiliFile) == 0 {
			// 2.3.1 se abbiamo tutte le liste vuote allora ogni possibile mattoncino è valido
			for _, m := range possibiliMattoncini {
				possibiliFile = append(possibiliFile, []NomeMattoncino{m})
			}
		} else {
			// 2.3.2 aggiungo solo le liste possibili

			// 2.3.2.1 segno i mattoncini usati da ogni fila
			mattonciniUsatiPossibiliFile := make([]SetMattoncini, len(possibiliFile))
			for i, possibileFila := range possibiliFile {
				mattonciniUsatiPossibiliFile[i] = make(SetMattoncini)
				for _, m := range possibileFila {
					mattonciniUsatiPossibiliFile[i][m] = true
				}
			}

			// 2.3.2.2 creo le nuove possibili file aggiungendo i possibili mattoncini validi
			var nuovePossibiliFile [][]NomeMattoncino
			for _, m := range possibiliMattoncini {
				for i, possibileFila := range possibiliFile {
					if !mattonciniUsatiPossibiliFile[i][m] {
						nuovePossibiliFile = append(nuovePossibiliFile, append(possibileFila, m))
					}
				}
			}

			// 2.3.2.3 controllo le nuove possibili file
			if len(nuovePossibiliFile) == 0 {
				fmt.Println("indefinito")
				return
			}

			possibiliFile = nuovePossibiliFile
		}
	}

	// controllo liste ottenute
	if len(possibiliFile) == 0 {
		fmt.Println("indefinito")
		return
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

		ordineRelativo := make([]int, n+1)
		ordineRelativoPrecedente := make([]int, n+1)

		for i := 1; i <= m; i++ {
			ordineRelativo, ordineRelativoPrecedente = ordineRelativoPrecedente, ordineRelativo

			for j := 1; j <= n; j++ {
				if a[i-1] == b[j-1] {
					ordineRelativo[j] = ordineRelativoPrecedente[j-1] + 1
				} else {
					if ordineRelativo[j-1] > ordineRelativoPrecedente[j] {
						ordineRelativo[j] = ordineRelativo[j-1]
					} else {
						ordineRelativo[j] = ordineRelativoPrecedente[j]
					}
				}
			}
		}

		return ordineRelativo[n]
	}

	// 4. calcolo il costo
	costoMinimo := INFINITO
	for _, mattonciniListaBordi := range possibiliFile {
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
	g := gioco{mattoncini: make(Mattoncini), file: make(File), scatola: make(Scatola)}

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

		fmt.Println("----------")
	}
}
