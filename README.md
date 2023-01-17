# ChordMiniProgect
mini progetto per il corso di Sistemi distribuiti e Cloud Computing della facoltà di Ingegneria Informatica Magistrale all'università degli studi di Roma - Tor Vergata.  
L'obiettivo del progetto é creare un sistema di chord semplicie utilizzando il linguaggio go e e le go-rpc
# Requisiti
il progetto richiede di avere sul proprio pc una versione di docker 3.9 o superiore
# Come avviare il progetto
bisognerà eseguire il comando all'interno della cartella scaricata
```
bash avvio.sh -x [Numero_di_nodi]
```
e verra instanziato un server_registry sulla porta 8000 e tanti nodi quanti specificati dal flag x sulle porte dalla [8000,...,8000+x]<br>
Inseguito si potrà accedere al nodo utilizzando il file eseguibile client in due modalità:
- put: si potrà memorizzare un valore nel sistema che restituirà la chiave dove é stato memorizzato
```
./client 1 string
```
- get: si potrà prendere il valore dal sistema passando come parametro la chiave
```
./client 0 key
```
