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

type Fondo struct{
	cassa float32
	cc float32
}

//costanti
const CASSA float32 = 500
const CC float32 = 500
const KGACQUISTO int = 1
const KGVENDITA int = 10
const X int = 500
const N int = 50
const PVp float32 = 1.5
const PAp float32 = 1.5
const CLIENTI int = 20
const AGRICOLTORI int = 3
const MAX_ATTESA int = 6
const MAX_BUFF int = 20

//canali
var done = make(chan bool)
var termina = make(chan bool)
var terminaAgricoltori = make(chan bool)
var acquistoCassa = make(chan int, MAX_BUFF)
var acquistoCC = make(chan int, MAX_BUFF)
var consegnaCassa = make(chan int, MAX_BUFF)
var consegnaCC = make(chan int, MAX_BUFF)
//ack
var ack_acquisto[CLIENTI] chan int
var ack_consegna[AGRICOLTORI] chan bool
//funzioni utility
func when(b bool, c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}


func clienti(id int){
	
	fmt.Printf("Inizializzazione Cliente %d\n",id)
	var tempo_di_attesa int
	var tipoPagamento int
	
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	tipoPagamento=rand.Intn(2)
	//faccio una richiesta
	if(tipoPagamento==0){
		fmt.Printf("[Acquirente %d] Richiesta acquisto %d kg di patate con pagamento in CONTANTI di %.2f euro\n",id,KGACQUISTO,float32(KGACQUISTO)*PAp)
		acquistoCassa <- id
	}else if(tipoPagamento==1){
		fmt.Printf("[Acquirente %d] Richiesta acquisto %d kg di patate con pagamento in BANCOMAT di %.2f euro\n",id,KGACQUISTO,float32(KGACQUISTO)*PAp)
		acquistoCC<- id
	}
	
	//attendo la risposta
	<- ack_acquisto[id]
	
	done <- true
}

func agricoltori(id int){
	fmt.Printf("Inizializzazione Fornitore %d\n",id)
	var tempo_di_attesa int
	var tipoPagamento int
	var res bool
	for{
			tempo_di_attesa = rand.Intn(MAX_ATTESA)+3 //+1 perchè randomizza da 0 a MAX_ATTESA
			time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
			tipoPagamento=rand.Intn(2)
			if(tipoPagamento==0){
				fmt.Printf("[Agricoltore %d] Richiesta consegna %d kg di patate con richiesta di pagamento in CONTANTI di %.2f euro\n",id,KGVENDITA,float32(KGVENDITA)*PVp)
				consegnaCassa <- id
			}else if(tipoPagamento==1){
				fmt.Printf("[Agricoltore %d] Richiesta consegna %d kg di patate con richiesta di pagamento in BONIFICO di %.2f euro\n",id,KGVENDITA,float32(KGVENDITA)*PVp)
				consegnaCC<- id
			}
			//attendo la risposta
			res=<- ack_consegna[id]
			if (res==true){
				fmt.Printf("Termina Fornitore %d\n",id)
				done <- true
				return
			}
		}
}


func mercato(){
	//ciclo di vita del server
	var fondo Fondo
	var magazzino int =0;
	var fine bool = false
	fondo.cassa=CASSA
	fondo.cc=CC
	for{
		select{
			case x:=<-when(magazzino>0&&(fondo.cassa<float32(X) || fondo.cassa>=float32(X) && len(acquistoCC)==0), acquistoCassa):
				magazzino-=KGACQUISTO
				fondo.cassa+=float32(KGACQUISTO)*PAp
				fmt.Printf("[Mercato] Cliente %d, acquisto di %d kg di patate con pagamento in CONTANTI effettuato [Cassa: %.2f, CC: %.2f, Magazzino: %d kg] \n",x,KGACQUISTO,fondo.cassa,fondo.cc,magazzino)
				ack_acquisto[x] <- 1
			case x:=<-when(magazzino>0&&(fondo.cassa>=float32(X) || fondo.cassa<float32(X) && len(acquistoCassa)==0), acquistoCC):
				magazzino-=KGACQUISTO
				fondo.cc+=float32(KGACQUISTO)*PAp
				fmt.Printf("[Mercato] Cliente %d, acquisto di %d kg di patate con pagamento in BANCOMAT effettuato [Cassa: %.2f, CC: %.2f, Magazzino: %d kg] \n",x,KGACQUISTO,fondo.cassa,fondo.cc,magazzino)
				ack_acquisto[x] <- 1
			case x:=<-when(magazzino+KGVENDITA<=N && (fondo.cassa>=float32(X) || fondo.cassa<float32(X) && len(consegnaCC)==0) && fine==false, consegnaCassa):
				magazzino+=KGVENDITA
				fondo.cassa-=float32(KGVENDITA)*PVp
				fmt.Printf("[Mercato] Agricoltore %d, vendita di %d kg di patate con richiesta di pagamento in CONTANTI effettuato [Cassa: %.2f, CC: %.2f, Magazzino: %d kg] \n",x,KGVENDITA,fondo.cassa,fondo.cc,magazzino)
				ack_consegna[x] <- false
			case x:=<-when(magazzino+KGVENDITA<=N && (fondo.cassa<=float32(X) || fondo.cassa>float32(X) && len(consegnaCassa)==0) &&fine==false, consegnaCC):
				magazzino+=KGVENDITA
				fondo.cc-=float32(KGVENDITA)*PVp
				fmt.Printf("[Mercato] Agricoltore %d, vendita di %d kg di patate con richiesta di pagamento in BONIFICO effettuato [Cassa: %.2f, CC: %.2f, Magazzino: %d kg] \n",x,KGVENDITA,fondo.cassa,fondo.cc,magazzino)
				ack_consegna[x] <- false
			case x:=<-when(fine==true, consegnaCassa):
				ack_consegna[x] <- true
			case x:=<-when(fine==true, consegnaCC):
				ack_consegna[x] <- true
			case <-terminaAgricoltori:
				fine=true
			case <-termina:
				done<-true
				return
		}
	}
}

func main(){
	
	fmt.Printf("Programma avviato\n")
	fmt.Printf("[CLIENTI: %d, AGRICOLTORI: %d, CASSA: %.2f, CC: %.2f, X: %d, PVp: %.2f, PAp: %.2f, N: %d]\n\n",CLIENTI,AGRICOLTORI,CASSA,CC,X,PVp,PAp,N)
	rand.Seed(time.Now().Unix())
	
	//inizializzo canali ack
	for i:=0; i<CLIENTI;i++{
		ack_acquisto[i] = make(chan int, MAX_BUFF)
	}
	for i:=0; i<AGRICOLTORI;i++{
		ack_consegna[i] = make(chan bool, MAX_BUFF)
	}
	
	//lancio threads
	for i:=0; i<CLIENTI;i++{
		go clienti(i)
	}
	for i:=0; i<AGRICOLTORI;i++{
		go agricoltori(i)
	}
	
	//lancio il server
	go mercato()
	
	//attendo la terminazione dei clients
	for i:=0; i<CLIENTI; i++{
		<-done
	}
	terminaAgricoltori <- true
	for i:=0; i<AGRICOLTORI; i++{
		<-done
	}
	//avviso il server di terminare
	termina <- true
	
	//attendo la terminazione del server
	<-done
	
	fmt.Printf("Programma terminato\n")
}
