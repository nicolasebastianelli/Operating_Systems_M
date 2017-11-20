// resourcepool_prio project main.go

package main

import (
	"fmt"
	"time"
)

const MAXPROC = 10
const MAXRES = 3
const MAXBUFF = 20

var richiesta = make(chan int, MAXBUFF)
var rilascio = make(chan int, MAXBUFF)
var risorsa [MAXPROC]chan int
var done = make(chan int)
var termina = make(chan int)

func when (b bool,c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}

func client(i int) {

	richiesta <- i
	r := <-risorsa[i]
	fmt.Printf("[client %d] uso della risorsa %d\n", i, r)
	time.Sleep(time.Second * 2)
	rilascio <- r
	done <- i
}

func server() {

	var disponibili int = MAXRES
	var res, p, i int
	var libera [MAXRES]bool

	// inizializzazioni:
	for i := 0; i < MAXRES; i++ {
		libera[i] = true
	}
	
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("nuovo ciclo server\n")
		select {
		case res = <- rilascio:
				disponibili++
				libera[res] = true
				fmt.Printf("[server]  restituita risorsa: %d  \n", res)
		case p = <-when(disponibili>0,richiesta):
				for i = 0; i < MAXRES && !libera[i]; i++ {
				}
				libera[i] = false
				disponibili--
				risorsa[p] <- i
				fmt.Printf("[server]  allocata risorsa %d a cliente %d \n", i, p)		
		case <-termina: // quando tutti i processi hanno finito
			fmt.Println("FINE !!!!!!")
			done <- 1
			return

		}
	}
}

func main() {
	var cli int
	fmt.Printf("\n quanti clienti (max %d)? ", MAXPROC)
	fmt.Scanf("%d", &cli)
	fmt.Println("clienti:", cli)

	//inizializzazione canali
	for i := 0; i < cli; i++ {
		risorsa[i] = make(chan int, MAXBUFF)
	}

	for i := 0; i < cli; i++ {
		go client(i)
	}
	go server()

	for i := 0; i < cli; i++ {
		<-done
	}
	termina <- 1
	<-done //attesa terminazione server

}