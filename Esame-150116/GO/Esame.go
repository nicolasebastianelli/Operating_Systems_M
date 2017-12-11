/*
 * 
 *  Created by Nicola Sebastianelli
 * 
 */

package main

import(
	"fmt"
	"time"
	"math/rand"
)

//costanti
const CLIENTEPISCINA int = 10
const CLIENTESPA int = 10
const CLIENTEPISCINASPA int = 10
const OPERATORIPISCINA int = 2
const OPERATORISPA int = 2
const MAX_ATTESA int = 10
const MAX_BUFF int = 20
const Np int =5
const Ns int =5

//canali
var done = make(chan bool)
var termina = make(chan bool)
var terminaOperatori = make(chan bool)
var entrataClienteSpa = make(chan int, MAX_BUFF)
var entrataClientePiscinaSpa = make(chan int, MAX_BUFF)
var entrataClientePiscina = make(chan int, MAX_BUFF)
var uscitaClienteSpa = make(chan int, MAX_BUFF)
var uscitaClientePiscinaSpa = make(chan int, MAX_BUFF)
var uscitaClientePiscina = make(chan int, MAX_BUFF)
var entrataOperatoreSpa = make(chan int, MAX_BUFF)
var entrataOperatorePiscina = make(chan int, MAX_BUFF)
var uscitaOperatoreSpa = make(chan int, MAX_BUFF)
var uscitaOperatorePiscina = make(chan int, MAX_BUFF)

//ack
var ack_ClienteSpa[CLIENTESPA] chan int
var ack_ClientePiscina[CLIENTEPISCINA] chan int
var ack_ClientePiscinaSpa[CLIENTEPISCINASPA] chan int
var ack_OperatorePiscina[OPERATORIPISCINA] chan int
var ack_OperatoreSpa[OPERATORISPA] chan int
//funzioni utility
func when(b bool, c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}


func clientePiscina(id int){
	
	fmt.Printf("Inizializzazione Cliente Piscina %d\n",id)
	var tempo_di_attesa int
	
	//faccio una richiesta
	entrataClientePiscina <- id
	
	//attendo la risposta
	<- ack_ClientePiscina[id]
	//faccio cose
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	//faccio una richiesta
	uscitaClientePiscina <- id
	<- ack_ClientePiscina[id]
	//ho finito il mio lavoro
	done <- true
}
func clienteSpa(id int){
	
	fmt.Printf("Inizializzazione Cliente Spa %d\n",id)
	var tempo_di_attesa int
	
	//faccio una richiesta
	entrataClienteSpa <- id
	
	//attendo la risposta
	<- ack_ClienteSpa[id]
	//faccio cose
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	//faccio una richiesta
	uscitaClienteSpa <- id
	
	//attendo la risposta
	<- ack_ClienteSpa[id]
	//ho finito il mio lavoro
	done <- true
}
func clientePiscinaSpa(id int){
	
	fmt.Printf("Inizializzazione Cliente Piscina e Spa %d\n",id)
	var tempo_di_attesa int
	
	//faccio una richiesta
	entrataClientePiscinaSpa<- id
	
	//attendo la risposta
	<- ack_ClientePiscinaSpa[id]
	//faccio cose
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	//faccio una richiesta
	uscitaClientePiscinaSpa <- id
	
	//attendo la risposta
	<- ack_ClientePiscinaSpa[id]
	//ho finito il mio lavoro
	done <- true
}

func operatorePiscina(id int){
	
	fmt.Printf("Inizializzazione Operatore Piscina %d\n",id)
	var tempo_di_attesa int
	for{
		
		select {
			case <- terminaOperatori:
				fmt.Printf("Termina operatore Piscina %d\n",id)
				done <- true
				return
			default:
				entrataOperatorePiscina <- id
				
				//attendo la risposta
				<- ack_OperatorePiscina[id]
				//faccio cose
				tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
				time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
				
				//faccio una richiesta
				uscitaOperatorePiscina <- id
				
				//attendo la risposta
				<- ack_OperatorePiscina[id]
				tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
				time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
		}
	}
}

func operatoreSpa(id int){	
	fmt.Printf("Inizializzazione Operatore Spa %d\n",id)
	var tempo_di_attesa int
	for{
		
		select {
			case <- terminaOperatori:
				fmt.Printf("Termina operatore Spa %d\n",id)
				done <- true
				return
			default:
				entrataOperatoreSpa <- id
				
				//attendo la risposta
				<- ack_OperatoreSpa[id]
				//faccio cose
				tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
				time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
				
				//faccio una richiesta
				uscitaOperatoreSpa <- id
				
				//attendo la risposta
				<- ack_OperatoreSpa[id]
				tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
				time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
		}
	}
}

func gestore(){
	var clientiPiscina int =0
	var clientiSpa int =0
	var operatoriPiscina int =0
	var operatoriSpa int =0 
	//ciclo di vita del server
	for{
		select{
			case x:=<-when(clientiPiscina<Np && len(entrataClientePiscinaSpa)==0 && operatoriPiscina>0 && operatoriSpa >0, entrataClientePiscina):
				clientiPiscina++
				fmt.Printf("Entrato Cliente Piscina %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_ClientePiscina[x] <- 1
			case x:=<-when(clientiSpa<Ns && len(entrataClientePiscinaSpa)==0 && len(entrataClientePiscina)==0 && operatoriSpa>0 && operatoriPiscina>0, entrataClienteSpa):
				clientiSpa++
				fmt.Printf("Entrato Cliente Spa %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_ClienteSpa[x] <- 1
			case x:=<-when(clientiPiscina<Np && clientiSpa< Ns && operatoriPiscina>0 &&operatoriSpa>0, entrataClientePiscinaSpa):
				clientiSpa++
				clientiPiscina++
				fmt.Printf("Entrato Cliente Piscina Spa %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_ClientePiscinaSpa[x] <- 1
			case x:=<-entrataOperatorePiscina:
				operatoriPiscina++
				fmt.Printf("Entrato Operatore Piscina %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_OperatorePiscina[x]<-1
			case x:=<-entrataOperatoreSpa:
				operatoriSpa++
				fmt.Printf("Entrato Operatore Spa %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_OperatoreSpa[x]<-1
			case x:=<-uscitaClientePiscina:
				clientiPiscina--
				fmt.Printf("Uscito Cliente Piscina %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_ClientePiscina[x] <- 1
			case x:=<-uscitaClienteSpa:
				clientiSpa--
				fmt.Printf("Uscito Cliente Spa %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_ClienteSpa[x] <- 1
			case x:=<-uscitaClientePiscinaSpa:
				clientiSpa--
				clientiPiscina--
				fmt.Printf("Uscito Cliente Piscina Spa %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_ClientePiscinaSpa[x] <- 1
			case x:=<-when(clientiPiscina==0 && clientiSpa==0 || operatoriPiscina>1 && operatoriSpa>1, uscitaOperatorePiscina):
				operatoriPiscina--
				fmt.Printf("Uscito Operatore Piscina %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_OperatorePiscina[x]<-1
			case x:=<-when(clientiPiscina==0 && clientiSpa==0 || operatoriPiscina>1 && operatoriSpa>1, uscitaOperatoreSpa):
				operatoriSpa--
				fmt.Printf("Uscito Operatore Spa %d [CP: %d CS: %d OP: %d OS: %d]\n",x,clientiPiscina,clientiSpa,operatoriPiscina,operatoriSpa)
				ack_OperatoreSpa[x]<-1
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
		ack_ClientePiscina[i] = make(chan int, MAX_BUFF)
	}
	for i:=0; i<CLIENTEPISCINASPA;i++{
		ack_ClientePiscinaSpa[i] = make(chan int, MAX_BUFF)
	}
	for i:=0; i<CLIENTESPA;i++{
		ack_ClienteSpa[i] = make(chan int, MAX_BUFF)
	}
	for i:=0; i<OPERATORIPISCINA;i++{
		ack_OperatorePiscina[i] = make(chan int, MAX_BUFF)
	}
	for i:=0; i<OPERATORISPA;i++{
		ack_OperatoreSpa[i] = make(chan int, MAX_BUFF)
	}
	
	//lancio threads
	for i:=0; i<CLIENTEPISCINA;i++{
		go clientePiscina(i)
	}
	for i:=0; i<CLIENTEPISCINASPA;i++{
		go clientePiscinaSpa(i)
	}
	for i:=0; i<CLIENTESPA;i++{
		go clienteSpa(i)
	}
	for i:=0; i<OPERATORIPISCINA;i++{
		go operatorePiscina(i)
	}
	for i:=0; i<OPERATORISPA;i++{
		go operatoreSpa(i)
	}
	
	//lancio il server
	go gestore()
	
	//attendo la terminazione dei clients
	for i:=0; i<CLIENTEPISCINA+CLIENTEPISCINASPA+CLIENTESPA; i++{
		<-done
	}
	
	for i:=0; i<OPERATORISPA+OPERATORIPISCINA;i++{
		terminaOperatori <- true
	}
	
	for i:=0; i<OPERATORISPA+OPERATORIPISCINA;i++{
		<- done
	}
	
	//avviso il server di terminare
	termina <- true
	
	//attendo la terminazione del server
	<-done
	
	fmt.Printf("Programma terminato\n")
}
