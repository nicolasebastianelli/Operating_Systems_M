package main

import (
	"fmt"
	"math/rand"
	"time"
)


const MAXCLI = 50
const N = 5
const K = 5

//definizione canali
var done = make(chan bool)
var termina = make(chan bool)
var entrataTUR = make(chan int, N)
var uscitaTUR = make(chan int, N)
var entrataEVE = make(chan int, N)
var uscitaEVE = make(chan int, N)
var ACK_TUR [MAXCLI]chan int 
var ACK_EVE [MAXCLI]chan int 

func when(b bool, c chan int) chan int {
	if !b {
		return nil
	}
	return c
}


func goroutine(myid int, tipo int) {
	var tt int
	tt = rand.Intn(5) + 1


     if tipo == 0{
            entrataTUR <- myid // send asincrona
	    <-ACK_TUR[myid] // attesa x sincronizzazione resto fermo fino a he il server non mi da conferma sul mio canale
            //fmt.Printf("")
	    tt = rand.Intn(5)+1
	    time.Sleep(time.Duration(tt) * time.Second)
	    uscitaTUR <- myid
     } else {
            entrataEVE <- myid // send asincrona
	    <-ACK_EVE[myid]    // attesa x sincronizzazione resto fermo fino a he il server non mi da conferma sul mio canale
            //fmt.Printf("")
	    tt = rand.Intn(5)+1
	    time.Sleep(time.Duration(tt) * time.Second)
	    uscitaEVE <- myid

     }

     done<-true

}


func server(){


var cont int = 0 //capacita
var contT int =0
var contE int =0
var serv int = 0 //se uguale 0 serve TUR altrimenti EVE

	for {

		select {
                //thread tipo1
		case x := <-when( ((serv==0 || len(entrataEVE)==0) && cont<N), entrataTUR):
			cont++
			contT++
			fmt.Printf("Entrato Cliente %d TUR [cont: %d contT: %d]\n", x,cont , contT)
			if(contT ==K){
				serv=1
				contT=0
			}
			ACK_TUR[x] <- 1 // termine "call"

                //thread tipo2 
        case x := <-when( ((serv==1 || len(entrataTUR)==0)&& cont<N), entrataEVE):
			cont++
			contE++
			fmt.Printf("Entrato Cliente %d EVE [cont: %d contE: %d]\n", x,cont , contE)
			if(contE ==K){
				serv=0
				contE=0
			}
			ACK_EVE[x] <- 1 // termine "call"

        
		case x := <-uscitaTUR:
			cont--
			fmt.Printf("Uscito Cliente %d TUR\n", x)

                
        case x := <-uscitaEVE:
			cont--
			fmt.Printf("Uscito Cliente %d EVE\n", x)

		case <-termina: // quando tutti i processi hanno finito
			fmt.Println("\n\n\nFINE!!!!!!")
			done <- true
			return
		}//select

	}//for


}//server







func main() {
	var TUR int
	var EVE int
    rand.Seed(time.Now().Unix())   

	TUR = rand.Intn(MAXCLI/2) + 1
	EVE = rand.Intn(MAXCLI/2) + 1
	
	fmt.Printf("Ci sono %d clienti TUR \n", TUR)
	fmt.Printf("Ci sono %d clienti EVE \n\n", EVE)
        
	for i := 0; i < TUR; i++ {
		ACK_TUR[i] = make(chan int)
	}
    for i := 0; i < EVE; i++ {
		ACK_EVE[i] = make(chan int)
	}
    
	go server()

	for i := 0; i < TUR; i++ {
		go goroutine(i,0)
	}
        for i := 0; i < EVE; i++ {
		go goroutine(i,1)
	}

	for i := 0; i < TUR+EVE; i++ {
		<-done
	}
	
        termina <- true
	<-done
	fmt.Printf("\nHO FINITO!!! ^_- \n")
}


