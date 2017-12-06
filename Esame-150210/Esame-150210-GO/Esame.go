package main

import (
	"fmt"
	"math/rand"
	"time"
)

const MAXNEVE = 10
const MAXSALE = 10
const N = 5
const K = 2

//definizione canali
var done = make(chan bool)
var terminaC = make(chan bool)
var termina = make(chan bool)
var entrataC = make(chan int, N)
var entrataN = make(chan int, N)
var entrataS = make(chan int, N) 
var uscitaC = make(chan int)
var uscitaN = make(chan int)
var uscitaS = make(chan int)
var ACK_NEVE [MAXNEVE]chan int 
var ACK_SALE [MAXSALE]chan int 
var ACK_CAMION = make(chan int)

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}


func spazzaneve(myid int) {
	var tt int
	tt = rand.Intn(5) + 1
    fmt.Printf("Inizializzazione Spazzaneve %d\n",myid)
	tt = rand.Intn(5)
    entrataN <- myid // send asincrona
    <-ACK_NEVE[myid]
    time.Sleep(time.Duration(tt) * time.Second)
    uscitaN <- myid
    done<-true
}

func spargisale(myid int) {
	var tt int
	tt = rand.Intn(5) + 1
	fmt.Printf("Inizializzazione Spargisale %d\n",myid)
    tt = rand.Intn(10)
    time.Sleep(time.Duration(tt) * time.Second)
    entrataS <- myid // send asincrona
	<-ACK_SALE[myid]      
    time.Sleep(time.Duration(tt) * time.Second)
    uscitaS <- myid
     done<-true
}

func camion() {
	var tt int
	tt = rand.Intn(5) + 1
	fmt.Printf("Inizializzazione Camion\n")
	tt = rand.Intn(5)
	time.Sleep(time.Duration(tt) * time.Second)
	for {
		select {
			case <- terminaC:
				fmt.Println("Termina camion")
				done <- true
			default:
			    entrataC <- 1 // send asincrona
			    <-ACK_CAMION   
			    time.Sleep(time.Duration(tt) * time.Second)
			    uscitaC <- 1
		}
	}
}


func server(){


var contMezzi int = 0
var silos int = K

	for {
		select {
                //thread tipo1
		case x := <-when( contMezzi < N && len(entrataN)==0, entrataC):
			contMezzi++
			silos=K;
			fmt.Printf("Entrato Camion %d [Mezzi: %d, Silos: %d]\n", x,contMezzi,silos)
			ACK_CAMION <- 1 // termine "call"

                //thread tipo2 
        case x := <-when(contMezzi < N , entrataN):
			contMezzi++
			fmt.Printf("Entrato Spazzaneve %d [Mezzi: %d, Silos: %d]\n", x,contMezzi,silos)
			ACK_NEVE[x] <- 1 // termine "call"
			
		case x := <-when( contMezzi < N && len(entrataC)==0 && len(entrataN)==0 && silos> 0, entrataS):
			contMezzi++
			silos--
			fmt.Printf("Entrato Spargisale %d [Mezzi: %d, Silos: %d]\n", x,contMezzi,silos)
			ACK_SALE[x] <- 1 // termine "call"

        
		case <-uscitaC:
			contMezzi--
			fmt.Printf("Uscito Camion\n")

                
        case x := <-uscitaN:
			contMezzi--
			fmt.Printf("Uscito Spazzaneve %d\n", x)
			
		case x := <-uscitaS:
			contMezzi--
			fmt.Printf("Uscito Spargisale %d\n", x)

		case <-termina: // quando tutti i processi hanno finito
			fmt.Println("\n\n\nFINE!!!!!!")
			done <- true
			return
		}//select

	}//for


}//server







func main() {
	rand.Seed(time.Now().Unix())
	var NEVE =rand.Intn(MAXNEVE)+1;
	var SALE=rand.Intn(MAXSALE)+1;
	
	fmt.Printf("Spazzaneve: %d\n", NEVE)
	fmt.Printf("Spargisale: %d\n\n", SALE)
	
	for i := 0; i < NEVE; i++ {
		ACK_NEVE[i] = make(chan int, N)
	}
        for i := 0; i < SALE; i++ {
		ACK_SALE[i] = make(chan int, N)
	}
    
    
	go server()

	for i := 0; i < NEVE; i++ {
		go spazzaneve(i)
	}
	for i := 0; i < SALE; i++ {
		go spargisale(i)
	}
	
	go camion()

	for i := 0; i < NEVE+SALE; i++ {
		<-done
	}
		terminaC <- true
        termina <- true
	<-done
	<-done
	fmt.Printf("\nHO FINITO!!! ^_- \n")
}


