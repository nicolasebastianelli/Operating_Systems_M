// ponte_sensounico_prio project main.go -- priorità di direzione ai veicoli da NORD
// ponte_senso_unico project main.go

package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXBUFF = 100
const MAXPROC = 10
const MAX = 3 // capacità
const N int = 0
const S int = 1

type person struct {
	idTaglia  int
	taglia    string
	idUtente  int
	direzione int
}

var done = make(chan bool)
var termina = make(chan bool)
var entrataNG = make(chan person, MAXBUFF) // canale grassi da nord
var entrataNM = make(chan person, MAXBUFF) //canale magri da nord
var entrataSG = make(chan person, MAXBUFF) //canale grassi da sud
var entrataSM = make(chan person, MAXBUFF) //canale magri da sud
var uscitaNG = make(chan person)
var uscitaNM = make(chan person)
var uscitaSG = make(chan person)
var uscitaSM = make(chan person)
var ACK_N [MAXPROC]chan int //risposte utenti nord
var ACK_S [MAXPROC]chan int //risposte utenti sud
var r int

func when(b bool, c chan person) chan person {
	if !b {
		return nil
	}
	return c
}

func utente(persona person) {
	var tt int
	tt = rand.Intn(5) + 1
	fmt.Printf("inizializzazione utente %d direzione %d taglia %s in secondi %d \n", persona.idUtente, persona.direzione, persona.taglia, tt)
	time.Sleep(time.Duration(tt) * time.Second)
	if persona.direzione == N {
		if persona.idTaglia == 1 {
			entrataNM <- persona       // send asincrona
			<-ACK_N[persona.idUtente] // attesa x sincronizzazione
			fmt.Printf("[Utente %d %s]  sul ponte in direzione  NORD\n", persona.idUtente, persona.taglia)
			tt = rand.Intn(5)
			time.Sleep(time.Duration(tt) * time.Second)
			uscitaNM <- persona
		} else{
			entrataNG <- persona       // send asincrona
			<-ACK_N[persona.idUtente] // attesa x sincronizzazione
			fmt.Printf("[Utente %d %s]  sul ponte in direzione  NORD\n", persona.idUtente, persona.taglia)
			tt = rand.Intn(5)
			time.Sleep(time.Duration(tt) * time.Second)
			uscitaNG <- persona
		}
	} else {
		if persona.idTaglia == 1 {
			entrataSM <- persona       // send asincrona
			<-ACK_S[persona.idUtente] // attesa x sincronizzazione
			fmt.Printf("[Utente %d %s]  sul ponte in direzione  SUD\n", persona.idUtente, persona.taglia)
			tt = rand.Intn(5)
			time.Sleep(time.Duration(tt) * time.Second)
			uscitaSM <- persona
		} else{
			entrataSG <- persona       // send asincrona
			<-ACK_S[persona.idUtente] // attesa x sincronizzazione
			fmt.Printf("[Utente %d %s]  sul ponte in direzione  SUD\n", persona.idUtente, persona.taglia)
			tt = rand.Intn(5)
			time.Sleep(time.Duration(tt) * time.Second)
			uscitaSG <- persona
		}
	}
	done <- true
}

func server() {

	var contNM int = 0
	var contNG int = 0
	var contSM int = 0
	var contSG int = 0

	for {

		select {
		case x := <-when((contNM+contSM+contNG+contSG < MAX) && (contSG ==0), entrataNM):
			contNM++
			fmt.Printf("[Ponte]  entrato utente %d %s in direzione N!  \n", x.idUtente, x.taglia)
			ACK_N[x.idUtente] <- 1 // termine "call"
		case x := <-when((contNM+contSM+contNG+contSG < MAX) && (contNG == 0), entrataSM):
			contSM++
			fmt.Printf("[Ponte]  entrato utente %d %s in direzione S!  \n", x.idUtente, x.taglia)
			ACK_S[x.idUtente] <- 1 // termine "call"

		case x := <-when((len(entrataSM) == 0) &&(len(entrataNM) == 0) && (contNM+contSM+contNG+contSG < MAX) && (contSG == 0)&&(contSM == 0), entrataNG):
			contNG++
			fmt.Printf("[Ponte]  entrato utente %d %s in direzione N!  \n", x.idUtente, x.taglia)
			ACK_N[x.idUtente] <- 1 // termine "call"
		case x := <-when((len(entrataNM) == 0) &&(len(entrataSM) == 0) && (contNM+contSM+contNG+contSG < MAX) && (contNG == 0) && (contNM == 0), entrataSG):
			contSG++
			fmt.Printf("[Ponte]  entrato utente %d %s in direzione S!  \n", x.idUtente, x.taglia)
			ACK_S[x.idUtente] <- 1 // termine "call"

		case x := <-uscitaNM:
			contNM--
			fmt.Printf("[Ponte]  uscito utente %d %s in direzione N!  \n", x.idUtente, x.taglia)
		case x := <-uscitaSM:
			contSM--
			fmt.Printf("[Ponte]  uscito utente %d %s in direzione S!  \n", x.idUtente, x.taglia)
		case x := <-uscitaNG:
			contNG--
			fmt.Printf("[Ponte]  uscito utente %d %s in direzione N!  \n", x.idUtente, x.taglia)
		case x := <-uscitaSG:
			contSG--
			fmt.Printf("[Ponte]  uscito utente %d %s in direzione S!  \n", x.idUtente, x.taglia)
		case <-termina: // quando tutti i processi hanno finito
			fmt.Println("FINE !!!!!!")
			done <- true
			return
		}

	}
}

func main() {
	var UN int
	var US int
	var persona person
	UN = MAXPROC + 1
	US = MAXPROC + 1
	for UN > MAXPROC {
		fmt.Printf("Quanti utenti NORD (max %d)? \n", MAXPROC)
		fmt.Scanf("%d", &UN)
		if UN > MAXPROC {
			fmt.Printf("Numero di utenti inserito maggiore del massimo\n")
		}
	}
	for US > MAXPROC {
		fmt.Printf("Quanti utenti SUD (max %d)? \n", MAXPROC)
		fmt.Scanf("%d", &US)
		if US > MAXPROC {
			fmt.Printf("Numero di utenti inserito maggiore del massimo\n")
		}
	}

	//inizializzazione canali
	for i := 0; i < UN; i++ {
		ACK_N[i] = make(chan int, MAXBUFF)
	}

	//inizializzazione canali
	for i := 0; i < US; i++ {
		ACK_S[i] = make(chan int, MAXBUFF)
	}

	rand.Seed(time.Now().Unix())
	go server()

	for i := 0; i < US; i++ {
		persona.idUtente = i
		persona.direzione = S
		if rand.Intn(2) == 1 {
			persona.idTaglia = 1
			persona.taglia = "Magro"
		} else {
			persona.idTaglia = 0
			persona.taglia = "Grasso"
		}
		go utente(persona)
	}
	for i := 0; i < UN; i++ {
		persona.idUtente = i
		persona.direzione = N
		if rand.Intn(2) == 1 {
			persona.idTaglia = 1
			persona.taglia = "Magro"
		} else {
			persona.idTaglia = 0
			persona.taglia = "Grasso"
		}
		go utente(persona)
	}

	for i := 0; i < UN+US; i++ {
		<-done
	}
	termina <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
