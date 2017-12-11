/*
 * 
 *  Created by Nicola Sebastianelli
 *	s0000850827
 * 
 */

package main

import(
	"fmt"
	"time"
	"math/rand"
)

type nuovo_tipo struct{
	id int
	dato int
}

//costanti
const COSTANTE1 int = 10
const COSTANTE2 int = 10
const ENTITA1 int = 2
const MAX_ATTESA int = 10
const MAX_BUFF int = 20

//canali
var done = make(chan bool)
var termina = make(chan bool)
var terminaEntita = make(chan bool)
var canale1 = make(chan int, MAX_BUFF)
var canale2 = make(chan int, MAX_BUFF)
var canaleEntita1 = make(chan int, MAX_BUFF)
var canaleEntita2 = make(chan int, MAX_BUFF)

//ack
var ack_canale1[COSTANTE1] chan int
var ack_canale2[COSTANTE2] chan int
var ack_canaleEntita1[COSTANTE1] chan int
var ack_canaleEntita2[COSTANTE2] chan int

//funzioni utility
func when(b bool, c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}


func thread1(id int){
	
	fmt.Printf("Inizializzazione Thread %d\n",id)
	var tempo_di_attesa int
	
	//faccio una richiesta
	canale1 <- id
	
	//attendo la risposta
	<- ack_canale1[id]


	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	//faccio una richiesta
	canale2 <- id
	<- ack_canale2[id]
	//ho finito il mio lavoro
	done <- true
}

func entita1(id int){
	
	fmt.Printf("Inizializzazione Entita %d\n",id)
	var tempo_di_attesa int
	for{
		
		select {
			case <- terminaEntita:
				fmt.Printf("Termina entita %d\n",id)
				done <- true
				return
			default:
				canaleEntita1 <- id
				
				//attendo la risposta
				<- ack_canaleEntita1[id]
				//faccio cose
				tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
				time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
				
				//faccio una richiesta
				canaleEntita2 <- id
				
				//attendo la risposta
				<- ack_canaleEntita2[id]
				tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
				time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
		}
	}
}


func server(){
	//ciclo di vita del server
	for{
		select{
			case x:=<-when(condizione, canale1):
				fmt.Printf("Stampa \n")
				ack_canale1[x] <- 1
			case x:=<-when(condizione, canale2):
				fmt.Printf("Stampa\n")
				ack_canale2[x] <- 1
			case <-termina:
				done<-true
				return
		}
	}
}

func main(){
	
	fmt.Printf("Programma avviato\n")
	rand.Seed(time.Now().Unix())
	
	//inizializzo canali ack
	for i:=0; i<CLIENTEPISCINA;i++{
		ack_canale1[i] = make(chan int, MAX_BUFF)
	}
	
	//lancio threads
	for i:=0; i<COSTANTE1;i++{
		go thread1(i)
	}
	for i:=0; i<ENTITA1;i++{
		go entita1(i)
	}
	
	//lancio il server
	go server()
	
	//attendo la terminazione dei clients
	for i:=0; i<COSTANTE1+COSTANTE2; i++{
		<-done
	}
	
	for i:=0; i<ENTITA1;i++{
		terminaEntita1 <- true
	}
	
	for i:=0; i<ENTITA1;i++{
		<- done
	}
	
	//avviso il server di terminare
	termina <- true
	
	//attendo la terminazione del server
	<-done
	
	fmt.Printf("Programma terminato\n")
}
