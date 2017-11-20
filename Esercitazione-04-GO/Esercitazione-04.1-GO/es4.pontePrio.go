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

var done = make(chan bool)
var termina = make(chan bool)
var entrataN = make(chan int, MAXBUFF) // necessità di accodamento per priorità
var entrataS = make(chan int, MAXBUFF) // necessità di accodamento per priorità
var uscitaN = make(chan int)
var uscitaS = make(chan int)
var ACK_N [MAXPROC]chan int //risposte client nord
var ACK_S [MAXPROC]chan int //risposte client sud
var r int

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}

func veicolo(myid int, dir int) {
	var tt int
	tt = rand.Intn(5) + 1
	fmt.Printf("inizializzazione veicolo  %d direzione %d in secondi %d \n", myid, dir, tt)
	time.Sleep(time.Duration(tt) * time.Second)
	if dir == N {

		entrataN <- myid // send asincrona
		<-ACK_N[myid]    // attesa x sincronizzazione
		fmt.Printf("[veicolo %d]  sul ponte in direzione  NORD\n", myid)
		tt = rand.Intn(5)
		time.Sleep(time.Duration(tt) * time.Second)
		uscitaN <- myid
	} else {
		entrataS <- myid
		<-ACK_S[myid] // attesa x sincronizzazione
		fmt.Printf("[veicolo %d]  sul ponte in direzione  SUD\n", myid)
		tt = rand.Intn(5)
		time.Sleep(time.Duration(tt) * time.Second)
		uscitaS <- myid
	}
	done <- true
}

func server() {

	var contN int = 0
	var contS int = 0

	for {

		select {
		case x := <-when((contN < MAX) && (contS == 0), entrataN):
			contN++
			fmt.Printf("[ponte]  entrato veicolo %d in direzione N!  \n", x)
			ACK_N[x] <- 1 // termine "call"

		case x := <-when((contS < MAX) && (contN == 0) && (len(entrataN) == 0), entrataS):
			contS++
			fmt.Printf("[ponte]  entrato veicolo %d in direzione S!  \n", x)
			ACK_S[x] <- 1 // termine "call"
		case x := <-uscitaN:
			contN--
			fmt.Printf("[ponte]  uscito veicolo %d in direzione N!  \n", x)
		case x := <-uscitaS:
			contS--
			fmt.Printf("[ponte]  uscito veicolo %d in direzione S!  \n", x)
		case <-termina: // quando tutti i processi hanno finito
			fmt.Println("FINE !!!!!!")
			done <- true
			return
		}

	}
}

func main() {
	var VN int
	var VS int

	fmt.Printf("\n quanti veicoli NORD (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &VN)
	fmt.Printf("\n quanti veicoli SUD (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &VS)

	//inizializzazione canali
	for i := 0; i < VN; i++ {
		ACK_N[i] = make(chan int, MAXBUFF)
	}

	//inizializzazione canali
	for i := 0; i < VS; i++ {
		ACK_S[i] = make(chan int, MAXBUFF)
	}

	rand.Seed(time.Now().Unix())
	go server()

	for i := 0; i < VS; i++ {
		go veicolo(i, S)
	}
	for i := 0; i < VN; i++ {
		go veicolo(i, N)
	}

	for i := 0; i < VN+VS; i++ {
		<-done
	}
	termina <- true
	<-done
	fmt.Printf("\n HO FINITO ")
}
