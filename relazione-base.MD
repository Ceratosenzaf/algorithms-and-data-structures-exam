// Davide Cerato (matricola 984481)

# bordoDaDirezione
- tempo: O(1)
- spazio: O(1)

---

# stampaMattoncinoInDirezione
- tempo: O(1)
- spazio: O(1)

---

# inserisciMattoncino
- tempo: O(1)
- spazio: O(1)

---

# stampaMattoncino
- tempo: O(1)
- spazio: O(1)

---

# disponiFila
- tempo: O(n+m) = O(m)
  - len(fila) = n
  - len(listaNomi) = m >= 3n
- spazio: O(n)
  - len(fila) = n

---

# stampaFila
- tempo: O(n)
  - len(fila) = n
- spazio: O(1)

---

# eliminaFila
- tempo: O(n)
  - len(fila) = n
- spazio: O(1)

---

# disponiFilaMinima
## nessun mattoncino con bordo alpha o beta
- tempo: O(1)
- spazio: O(1)
## singolo mattoncino con bordi alpha e beta
- tempo: O(1)
- spazio: O(1)
## alpha e beta diversi
- tempo: O(n+m+k) = O(n+m)
  - len(scatola) = n
  - len(bordi distinti) = m
  - len(fila) = k
- spazio: O(m+k) = O(m)
  - len(bordi distinti) = m
  - len(fila) = k
## alpha e beta uguali
- tempo: O(n+p*(n+m)+k) = O(n+p*(n+m))
  - len(scatola) = n
  - len(bordi distinti) = m
  - len(mattoncini con bordo alpha) = p
  - len(fila) = k
- spazio: O(m+p*m+k) = O(m+p*m)
  - len(bordi distinti) = m
  - len(mattoncini con bordo alpha) = p
  - len(fila) = k

---

# sottostringaMassima
- tempo: O(m*n)
  - len(sigma) = m
  - len(tao) = n
- spazio: O(m+n)
  - len(sigma) = m
  - len(tao) = n

---

# indiceCacofonia
- tempo: O(m*n\*k) = O(k) se le stringhe sono di lunghezza contenuta
  - len(fila) = k
  - maxLen(fila) = m
  - secondMaxLen(fila) = n
- spazio: O(k*(m+n)) = O(k) se le stringhe sono di lunghezza contenuta
  - len(fila) = k
  - maxLen(fila) = m
  - secondMaxLen(fila) = n

---

# costo
- tempo: O(k+m*n\*(n+k))
  - len(fila) = k
  - len(lista bordi) = n
  - len(set mattoncini lista bordi per bordi uguali) = m <= n
  - len(scatola) = p
- spazio: O(k+n*m)
  - len(fila) = k
  - len(lista bordi) = n
  - len(set mattoncini lista bordi per bordi uguali) = m <= n
## lista di adiacenza
- tempo: O(n+k)
  - len(scatola) = n
  - len(fila) = k
- spazio: O(k+n+m)
  - len(scatola) = n
  - len(fila) = k
  - len(bordi distinti) = m
## strutture dati
- tempo: O(k+n*m)
  - len(fila) = k
  - len(lista bordi) = n
  - len(bordi distinti) = m
- spazio: O(k+n*m)
  - len(fila) = k
  - len(lista bordi) = n
  - len(set mattoncini lista bordi per bordi uguali) = m
## funzione di calcolo
- tempo: O(m*n)
  - len(a) = m
  - len(b) = n
- spazio: O(m+n)
  - len(a) = m
  - len(b) = n
## calcolo effettivo
- tempo: O(k*m\*(n+m))
  - len(fila) = n
  - len(mattoncini lista bordi) = m = len(lista bordi) - 1
  - len(set mattoncini lista bordi per bordi uguali) = k
- spazio: O(m+n)
  - len(fila) = n
  - len(mattoncini lista bordi) = m = len(lista bordi) - 1

---

# main
- tempo: 
- spazio: 
