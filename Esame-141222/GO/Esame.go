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

type Cliente struct{
	id int
	ticket int
}

//costanti
const CLIENTI int = 12
const SFOGLINE int = 5
const MAX int = 10
const N int =5
const MAX_ATTESA int = 10
const MAX_BUFF int = 20

//canali
var done = make(chan bool)
var termina = make(chan bool)
var prenota = make(chan int, MAX_BUFF)
var ritira = make(chan Cliente, MAX_BUFF)
var deposita = make(chan int, MAX_BUFF)

//ack
var ack_prenota[CLIENTI] chan int
var ack_ritira[CLIENTI] chan int
var ack_deposita[SFOGLINE] chan int

var fine bool = false

//funzioni utility
func when(b bool, c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}

func whenR(b bool, c chan Cliente) chan Cliente{
	if(!b){
		return nil
	}
	return c
}


func cliente(id int){
	
	fmt.Printf("Inizializzazione Cliente %d\n",id)
	var tempo_di_attesa int
	var cliente Cliente
	var res int
	cliente.id=id
	//faccio una richiesta
	prenota <- cliente.id
	
	//attendo la risposta
	cliente.ticket=<- ack_prenota[id]

	if(cliente.ticket!=-1){
		tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
		time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
		
		//faccio una richiesta
		ritira <- cliente
		res = <- ack_ritira[id]
		if(res == 1){
			fmt.Printf("Cliente %d con ticket %d tortellini ritirati\n",cliente.id,cliente.ticket)
		}else{
			fmt.Printf("Cliente %d ticket %d già utilizzato\n",cliente.id,cliente.ticket)
		}
	}
	//ho finito il mio lavoro
	done <- true
}

func sfoglina(id int){
	var res int
	fmt.Printf("Inizializzazione Sfoglina %d\n",id)
	var tempo_di_attesa int
	for{
				tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
				time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
				deposita <- id
				
				//attendo la risposta
				res=<- ack_deposita[id]
				if(res==-1){
					done<-true
					return
				}
	}			
}


func server(){
	var ticketUsati[CLIENTI] bool
	var frigorifero int =0
	var ticket int =0
	var prenotazione int = 0
	for i:=0;i<CLIENTI;i++{
		ticketUsati[i]=false;
	} 
	//ciclo di vita del server
	for{
		select{
			case x:=<-prenota:
				if(prenotazione >= MAX){
					fmt.Printf("Prenotazione del Cliente %d NON effettuta, [prenotazioni: %d, MAX %d] \n",x,prenotazione,MAX)
					ack_prenota[x] <- -1
				}else{
					prenotazione++;
					fmt.Printf("Prenotazione del Cliente %d effettuta, [ticket: %d, prenotazioni: %d]  \n",x,ticket,prenotazione)
					ack_prenota[x] <- ticket
					ticket++
				}
			case x:=<-whenR(len(prenota)==0&&frigorifero>0, ritira):
				fmt.Printf("[Server] richiesta ritiro cliente %d con ticket %d\n",x.id,x.ticket)
				if(ticketUsati[x.ticket]==false){
					ticketUsati[x.ticket]=true
					prenotazione--
					frigorifero--
					ack_ritira[x.id] <- 1
				}else{
					ack_ritira[x.id] <- -1
				}
			case x:=<-when(fine==false &&frigorifero<N, deposita):
				frigorifero++
				fmt.Printf("Depisito da Sfoglina %d [Frigorifero: %d]\n",x,frigorifero)
				ack_deposita[x] <- 1
			case x:=<-when(fine==true, deposita):
				fmt.Printf("Termino Sfoglina %d [Frigorifero: %d]\n",x,frigorifero)
				ack_deposita[x] <- -1
			case <-termina:
				done<-true
				fmt.Printf("Sono rimaste %d confezioni in frigo\n",frigorifero)
				return
		}
	}
}

func main(){
	
	fmt.Printf("Programma avviato\n")
	fmt.Printf("Clienti: %d, Sfogline: %d, Max Frigorifero: %d, Max Prenotazioni: %d\n\n",CLIENTI,SFOGLINE,N,MAX)
	rand.Seed(time.Now().Unix())
	
	//inizializzo canali ack
	for i:=0; i<CLIENTI;i++{
		ack_prenota[i] = make(chan int, MAX_BUFF)
		ack_ritira[i] = make(chan int, MAX_BUFF)
	}
	for i:=0; i<SFOGLINE;i++{
		ack_deposita[i] = make(chan int, MAX_BUFF)
	}
	
	//lancio threads
	for i:=0; i<SFOGLINE;i++{
		go sfoglina(i)
	}
	for i:=0; i<CLIENTI;i++{
		go cliente(i)
	}
	
	
	//lancio il server
	go server()
	
	//attendo la terminazione dei clients
	for i:=0; i<CLIENTI; i++{
		<-done
	}
	
	fine = true
	
	for i:=0; i<SFOGLINE;i++{
		<- done
	}
	
	//avviso il server di terminare
	termina <- true
	
	//attendo la terminazione del server
	<-done
	
	fmt.Printf("Programma terminato\n")
}
