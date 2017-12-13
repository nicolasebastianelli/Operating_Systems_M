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

type Magazzino struct{
	pieni int
	vuoti int
}

type Fondo struct{
	cassa float32
	cc float32
}

//costanti
const PIENI int = 25
const VUOTI int = 25
const CASSA float32 = 500
const CC float32 = 500
const Z int = 15
const Y int = 10
const X int = 10
const MAX_P int = 50
const MAX_V int = 50
const PV float32 = 1.5
const K int = 10
const PA float32 = 1.5
const ACQUIRENTI int = 10
const FORNITORI int = 3
const MAX_ATTESA int = 6
const MAX_BUFF int = 20

//canali
var done = make(chan bool)
var termina = make(chan bool)
var terminaFornitori = make(chan bool)
var acquistoCassa = make(chan int, MAX_BUFF)
var acquistoCC = make(chan int, MAX_BUFF)
var consegnaCassa = make(chan int, MAX_BUFF)
var consegnaCC = make(chan int, MAX_BUFF)
//ack
var ack_acquisto[ACQUIRENTI] chan int
var ack_consegna[FORNITORI] chan bool
//funzioni utility
func when(b bool, c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}


func acquirenti(id int){
	
	fmt.Printf("Inizializzazione Acquirente %d\n",id)
	var tempo_di_attesa int
	var tipoPagamento int
	
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	tipoPagamento=rand.Intn(2)
	//faccio una richiesta
	if(tipoPagamento==0){
		fmt.Printf("[Acquirente %d] Richiesta acquisto %d bottiglie di acqua e consegno %d vuoti con pagamento in CONTANTI di %.2f euro\n",id,Y,Z,float32(Y)*PA)
		acquistoCassa <- id
	}else if(tipoPagamento==1){
		fmt.Printf("[Acquirente %d] Richiesta acquisto %d bottiglie di acqua e consegno %d vuoti con pagamento in BANCOMAT di %.2f euro\n",id,Y,Z,float32(Y)*PA)
		acquistoCC<- id
	}
	
	//attendo la risposta
	<- ack_acquisto[id]
	
	done <- true
}

func fornitori(id int){
	
	fmt.Printf("Inizializzazione Fornitore %d\n",id)
	var tempo_di_attesa int
	var tipoPagamento int
	var res bool
	for{
			tempo_di_attesa = rand.Intn(MAX_ATTESA)+3 //+1 perchè randomizza da 0 a MAX_ATTESA
			time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
			tipoPagamento=rand.Intn(2)
			if(tipoPagamento==0){
				fmt.Printf("[Fornitore %d] Richiesta consegna %d bottiglie di acqua e ritiro tutti i vuoti con richiesta di pagamento in CONTANTI di %.2f euro\n",id,X,float32(X)*PV)
				consegnaCassa <- id
			}else if(tipoPagamento==1){
				fmt.Printf("[Fornitore %d] Richiesta consegna %d bottiglie di acqua e ritiro tutti i vuoti con richiesta di pagamento in BONIFICO di %.2f euro\n",id,X,float32(X)*PV)
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


func ditta(){
	//ciclo di vita del server
	var magazzino Magazzino
	var fondo Fondo
	var fine bool = false
	fondo.cassa=CASSA
	fondo.cc=CC
	magazzino.pieni=PIENI
	magazzino.vuoti=VUOTI
	for{
		select{
			case x:=<-when((magazzino.vuoti+Z<=MAX_V && magazzino.pieni-Y>=0) && (fondo.cassa<=float32(K) || fondo.cassa>float32(K) && len(acquistoCC)==0), acquistoCassa):
				magazzino.pieni-=Y
				magazzino.vuoti+=Z
				fondo.cassa+=float32(Y)*PA
				fmt.Printf("[Ditta] Acquirente %d, acquisto di %d bottiglie di acqua e consegna di %d vuoti con pagamento in CONTANTI effettuato [Cassa: %.2f, CC: %.2f, Bottiglie: %d, Vuoti: %d] \n",x,Y,Z,fondo.cassa,fondo.cc,magazzino.pieni,magazzino.vuoti)
				ack_acquisto[x] <- 1
			case x:=<-when((magazzino.vuoti+Z<=MAX_V && magazzino.pieni-Y>=0) && (fondo.cassa>=float32(K) || fondo.cassa<float32(K) && len(acquistoCassa)==0), acquistoCC):
				magazzino.pieni-=Y
				magazzino.vuoti+=Z
				fondo.cc+=float32(Y)*PA
				fmt.Printf("[Ditta] Acquirente %d, acquisto di %d bottiglie di acqua e consegna di %d vuoti con pagamento in BANCOMAT effettuato [Cassa: %.2f, CC: %.2f, Bottiglie: %d, Vuoti: %d] \n",x,Y,Z,fondo.cassa,fondo.cc,magazzino.pieni,magazzino.vuoti)
				ack_acquisto[x] <- 1
			case x:=<-when(magazzino.pieni+X<=MAX_P && (fondo.cassa>=float32(K) || fondo.cassa<float32(K) && len(consegnaCC)==0) && fine==false, consegnaCassa):
				magazzino.pieni+=X
				magazzino.vuoti=0
				fondo.cassa-=float32(X)*PV
				fmt.Printf("[Ditta] Fornitore %d, vendita di %d bottiglie di acqua e consegna di tutti i vuoti con richiesta di pagamento in CONTANTI effettuato [Cassa: %.2f, CC: %.2f, Bottiglie: %d, Vuoti: %d] \n",x,X,fondo.cassa,fondo.cc,magazzino.pieni,magazzino.vuoti)
				ack_consegna[x] <- false
			case x:=<-when(magazzino.pieni+X<=MAX_P && (fondo.cassa<=float32(K) || fondo.cassa>float32(K) && len(consegnaCassa)==0) &&fine==false, consegnaCC):
				magazzino.pieni+=X
				magazzino.vuoti=0
				fondo.cc-=float32(X)*PV
				fmt.Printf("[Ditta] Fornitore %d, vendita di %d bottiglie di acqua e consegna di tutti i vuoti con richiesta di pagamento in BONIFICO effettuato [Cassa: %.2f, CC: %.2f, Bottiglie: %d, Vuoti: %d] \n",x,X,fondo.cassa,fondo.cc,magazzino.pieni,magazzino.vuoti)
				ack_consegna[x] <- false
			case x:=<-when(fine==true, consegnaCassa):
				ack_consegna[x] <- true
			case x:=<-when(fine==true, consegnaCC):
				ack_consegna[x] <- true
			case <-terminaFornitori:
				fine=true
			case <-termina:
				done<-true
				return
		}
	}
}

func main(){
	
	fmt.Printf("Programma avviato\n")
	fmt.Printf("[ACQUIRENTI: %d, FORNITORI: %d, PIENI: %d, VUOTI: %d, CASSA: %.2f, CC: %.2f, X: %d, Y: %d, Z: %d, MAX_P: %d, MAX_V: %d, PV: %.2f, PA: %.2f, K: %d]\n\n",ACQUIRENTI,FORNITORI,PIENI,VUOTI,CASSA,CC,X,Y,Z,MAX_P,MAX_V,PV,PA,K)
	rand.Seed(time.Now().Unix())
	
	//inizializzo canali ack
	for i:=0; i<ACQUIRENTI;i++{
		ack_acquisto[i] = make(chan int, MAX_BUFF)
	}
	for i:=0; i<FORNITORI;i++{
		ack_consegna[i] = make(chan bool, MAX_BUFF)
	}
	
	//lancio threads
	for i:=0; i<ACQUIRENTI;i++{
		go acquirenti(i)
	}
	for i:=0; i<FORNITORI;i++{
		go fornitori(i)
	}
	
	//lancio il server
	go ditta()
	
	//attendo la terminazione dei clients
	for i:=0; i<ACQUIRENTI; i++{
		<-done
	}
	terminaFornitori <- true
	for i:=0; i<FORNITORI; i++{
		<-done
	}
	//avviso il server di terminare
	termina <- true
	
	//attendo la terminazione del server
	<-done
	
	fmt.Printf("Programma terminato\n")
}
