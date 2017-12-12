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

type tipo_cliente struct{
	id int
	biglietto int
}

//costanti
const LOCALI int = 0
const OSPITI int = 1
const CLIENTI int = 50
const OPERATORI int = 10
const MAX_ATTESA int = 10
const MAX_BUFF int = 20

//canali
var done = make(chan bool)
var termina = make(chan bool)
var acquisto = make(chan int, MAX_BUFF)
var varcoLocali = make(chan tipo_cliente, MAX_BUFF)
var varcoOspiti = make(chan tipo_cliente, MAX_BUFF)
var fineControllo = make(chan tipo_cliente, MAX_BUFF)
//ack
var ack_acquisto[CLIENTI] chan int
var ack_varco[CLIENTI] chan int

//funzioni utility
func when(b bool, c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}

func whenC(b bool, c chan tipo_cliente) chan tipo_cliente{
	if(!b){
		return nil
	}
	return c
}



func cliente(id int){
	
	var cli tipo_cliente
	cli.id=id
	
	fmt.Printf("Inizializzazione Cliente %d\n",id)
	var tempo_di_attesa int
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	//faccio una richiesta
	acquisto <- id
	//attendo la risposta
	cli.biglietto= <- ack_acquisto[id]
	if(cli.biglietto==LOCALI){
		fmt.Printf("[Cliente %d] Acquistato biglietto tribuna LOCALI\n",id)
	}else if(cli.biglietto==OSPITI){
		fmt.Printf("[Cliente %d] Acquistato biglietto tribuna OSPITI\n",id)
	}
	
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	//faccio una richiesta
	if(cli.biglietto==LOCALI){
		varcoLocali<-cli
	}else if(cli.biglietto==OSPITI){
		varcoOspiti<-cli
	}
	
	<- ack_varco[id]
	
	if(cli.biglietto==LOCALI){
		fmt.Printf("[Cliente %d] Entrato nella tribuna LOCALI\n",id)
	}else if(cli.biglietto==OSPITI){
		fmt.Printf("[Cliente %d] Entrato nella tribuna OSPITI\n",id)
	}
	
	//ho finito il mio lavoro
	done <- true
}

func varco(){
	var Operatori int = OPERATORI
	var Locali int=0
	var Ospiti int=0
	var biglietto int
	//ciclo di vita del server
	
	for{
		select{
			case x:=<-acquisto:
				biglietto = rand.Intn(2)
				ack_acquisto[x] <- biglietto
				if(biglietto==LOCALI){
					fmt.Printf("[Server] Venduto biglietto tribuna LOCALI al cliente %d\n",x)
				}else if(biglietto==OSPITI){
					fmt.Printf("[Server] Venduto biglietto tribuna OSPITI al cliente %d\n",x)
				}
			case x:=<-whenC(((Locali>=Ospiti && len(varcoOspiti)!=0)||len(varcoOspiti)==0)&& Operatori!=0 , varcoLocali):
				Operatori--
				fmt.Printf("[Server] Inizio controllo Cliente %d con biglietto per tribuna LOCALI [Operatori: %d, Locali: %d, Ospiti: %d]\n",x.id,Operatori,Locali,Ospiti)
				go operatore(x)
			case x:=<-whenC(((Ospiti>=Locali && len(varcoLocali)!=0)||len(varcoLocali)==0)&& Operatori!=0 , varcoOspiti):
				Operatori--
				fmt.Printf("[Server] Inizio controllo Cliente %d con biglietto per tribuna OSPITI [Operatori: %d, Locali: %d, Ospiti: %d]\n",x.id,Operatori,Locali,Ospiti)
				go operatore(x)
			case x:=<-fineControllo:
				Operatori++
				if(x.biglietto==LOCALI){
					Locali++
					fmt.Printf("[Server] Fine controllo Cliente %d con biglietto per tribuna LOCALI [Operatori: %d, Locali: %d, Ospiti: %d]\n",x.id,Operatori,Locali,Ospiti)
				}else if(x.biglietto==OSPITI){
					Ospiti++
					fmt.Printf("[Server] Fine controllo Cliente %d con biglietto per tribuna OSPITI [Operatori: %d, Locali: %d, Ospiti: %d]\n",x.id,Operatori,Locali,Ospiti)
				}
				ack_varco[x.id]<-1
			case <-termina:
				done<-true
				return
		}
	}
}

func operatore(cliente tipo_cliente){
	time.Sleep(time.Duration(rand.Intn(MAX_ATTESA)+1 ) * time.Second)
	fineControllo<-cliente
}

func main(){
	
	fmt.Printf("Programma avviato\n")
	rand.Seed(time.Now().Unix())
	
	//inizializzo canali ack
	for i:=0; i<CLIENTI;i++{
		ack_acquisto[i] = make(chan int, MAX_BUFF)
		ack_varco[i] = make(chan int, MAX_BUFF)
	}
	
	//lancio threads
	for i:=0; i<CLIENTI;i++{
		go cliente(i)
	}
	
	//lancio il server
	go varco()
	
	//attendo la terminazione dei clients
	for i:=0; i<CLIENTI; i++{
		<-done
	}
	
	//avviso il server di terminare
	termina <- true
	
	//attendo la terminazione del server
	<-done
	
	fmt.Printf("Programma terminato\n")
}
