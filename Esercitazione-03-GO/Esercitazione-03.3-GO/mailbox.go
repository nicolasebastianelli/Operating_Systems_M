package main

import (
	"fmt"
	"time"
	"math/rand"
)


const DIM = 3
const MAXCONS = 10
const MAXPROD = 10
const MAXNMESS = 3
var r = rand.New(rand.NewSource(time.Now().UnixNano()))
var ncons = r.Intn(MAXCONS)+1
var nprod = r.Intn(MAXPROD)+1
var nmess = r.Intn(MAXNMESS)+1
var DATIPROD = make(chan Message)
var DATICONS = make(chan Message)
var buf_in[DIM-1] chan Message
var readymsg =make([]chan int,ncons)
var done = make(chan int)
var termina [DIM] chan int

type Message struct {
  from int
  to int
  content string
}

func produttore(i int){
	var mess Message
	fmt.Printf("[Produttore %d] Avvio\n", i+1)
	mess.from=i;
	for j:= 0;j<nmess;j++{
		mess.to = r.Intn(ncons)
		DATIPROD<- mess
		fmt.Printf("[Produttore %d] Invio messaggio numero %d per destinatario: %d\n", i+1,j+1,mess.to+1)
	}
	done<-1
	fmt.Printf("[Produttore %d] Termino\n", i+1)
}


func consumatore(i int){
	fmt.Printf("[Consumatore %d] Avvio\n", i+1)
	for{
			_,ok:= <-readymsg[i]
			if !ok{
				fmt.Printf("[Consumatore %d] Termino\n", i+1)
				done <- 1
				return
			}
			mess:= <-DATICONS
			fmt.Printf("[Consumatore %d] ricevuto messaggio da produttore %d\n", i+1, mess.from+1)
	}
}

func server(stage int, ch_in chan Message, ch_out chan Message) { 
	fmt.Printf("[Stage %d] Avvio\n", stage+1)
	var mess Message
	if stage!=DIM-1{
	for{
		select {
			case mess=<- ch_in:
				fmt.Printf("[Stage %d] ricevuto messaggio da produttore %d per destinatario %d \n", stage+1, mess.from+1, mess.to+1)
					ch_out<-mess
			case <-termina[stage]: 
				fmt.Printf("[Stage %d] Termino\n", stage+1)
				termina[stage+1] <- 1
				done <- 1
				return
		}
	}
	
	} else {
		for{
		select {
			case mess=<- ch_in:
				fmt.Printf("[Stage %d] ricevuto messaggio da produttore %d per destinatario %d \n", stage+1, mess.from+1, mess.to+1)
					readymsg[mess.to]<-1
					ch_out<-mess
			case <-termina[stage]: 
				fmt.Printf("[Stage %d] Chiusura canali\n", stage+1)
				for i:=0;i< ncons;i++{
					close(readymsg[i])
				}
				done <- 1
				fmt.Printf("[Stage %d] Termino\n", stage+1)
				return
		}
	}
	}
	
	
}

func main() {
	fmt.Printf("Numero consumatori: %d\n", ncons)
	fmt.Printf("Numero produttori: %d\n", nprod)
	fmt.Printf("Numero messaggi per produttore: %d\n", nmess)
	fmt.Printf("Numero stage server: %d\n\n", DIM)
	
	//inizializzazione canale in/out stage
	for i := 0; i < DIM-1; i++ {
		buf_in[i] = make(chan Message)
	}
	//inizializzazione canale termina stage
	for i := 0; i < DIM; i++ {
		termina[i] = make(chan int)
	}
	//inizializzazione canale ready stage
	for i := 0; i < ncons; i++ {
		readymsg[i] = make(chan int)
	}
	
	//creazione goroutine server
	for i := 0; i < DIM; i++ {
		if(i==0){
			go server(i,DATIPROD,buf_in[i])
		} else if(i==DIM-1){
			go server(i,buf_in[i-1],DATICONS)	
		} else{
		go server(i,buf_in[i-1],buf_in[i])
		}
	}
	
	//creazione goroutine produttore
	for i := 0; i < nprod; i++ {
		go produttore(i)
	}
	
	//creazione goroutine consumatore
	for i := 0; i < ncons; i++ {
		go consumatore(i)
	}
	
	//attendo terminazione produttori
	for i := 0; i < nprod; i++ {
		<-done
	}
	
	// invio segnale di terminazione server
	termina[0] <- 1
	
	//attendo terminazione server
	for i := 0; i < DIM; i++ {
		<-done
	}

	//termina consumatore
	for i := 0; i < ncons; i++ {
		<-done
	}
}