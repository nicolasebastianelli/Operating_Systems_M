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
const BIANCO int = 0
const VERDE int = 1
const GIALLO int = 2
const ROSSO int = 3
const ADULTO int = 0
const BAMBINO int = 1

const NMEDICI int = 2
const NPAZIENTI int = 15
const MAX_ATTESA int = 10
const MAX_BUFF int = 20

//canali
var done = make(chan bool)
var terminaTriage = make(chan bool)
var terminaPediatra = make(chan bool)
var terminaAdulti= make(chan bool)
var triage = make(chan int, MAX_BUFF)
var rossoAdulti = make(chan int, MAX_BUFF)
var gialloAdulti = make(chan int, MAX_BUFF)
var verdeAdulti = make(chan int, MAX_BUFF)
var biancoAdulti = make(chan int, MAX_BUFF)
var rossoBambino = make(chan int, MAX_BUFF)
var gialloBambino = make(chan int, MAX_BUFF)
var verdeBambino = make(chan int, MAX_BUFF)
var biancoBambino = make(chan int, MAX_BUFF)
var visitaAdulti = make(chan int, MAX_BUFF)
var visitaPediatrica = make(chan int, MAX_BUFF)

//ack
var ack_triage[NPAZIENTI] chan int
var ack_paziente[NPAZIENTI] chan int


//funzioni utility
func when(b bool, c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}



func funcPaziente(id int){
	
	var eta int = rand.Intn(2)
	var codice int
	
	if(eta==ADULTO){
		fmt.Printf("Inizializzazione Paziente %d ADULTO\n",id)
	}else if(eta==BAMBINO){
		fmt.Printf("Inizializzazione Paziente %d BAMBINO\n",id)
	}
	var tempo_di_attesa int
	
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	//faccio una richiesta
	triage <- id
	
	//attendo la risposta
	codice = <- ack_triage[id]
	
	if(eta==ADULTO){
		fmt.Printf("[Paziente %d] Mi dirigo verso ambulatorio ADULTI\n",id)
		tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
		time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
		if(codice==ROSSO){
			fmt.Printf("[Paziente %d, Codice ROSSO] Entro in ambulatorio ADULTI\n",id)
			rossoAdulti<-id
		} else if(codice==GIALLO){
			fmt.Printf("[Paziente %d, Codice GIALLO] Entro in ambulatorio ADULTI\n",id)
			gialloAdulti<-id
		} else if(codice==VERDE){
			fmt.Printf("[Paziente %d, Codice VERDE] Entro in ambulatorio ADULTI\n",id)
			verdeAdulti<-id
		} else if(codice==BIANCO){
			fmt.Printf("[Paziente %d, Codice BIANCO] Entro in ambulatorio ADULTI\n",id)
			biancoAdulti<-id
		}
	}else if(eta==BAMBINO){
		fmt.Printf("[Paziente %d] Mi dirigo verso ambulatorio PEDIATRA\n",id)
		tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
		time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
		if(codice==ROSSO){
			fmt.Printf("[Paziente %d, Codice ROSSO] Entro in ambulatorio PEDIATRA\n",id)
			rossoBambino<-id
		} else if(codice==GIALLO){
			fmt.Printf("[Paziente %d, Codice GIALLO] Entro in ambulatorio PEDIATRA\n",id)
			gialloBambino<-id
		} else if(codice==VERDE){
			fmt.Printf("[Paziente %d, Codice VERDE] Entro in ambulatorio PEDIATRA\n",id)
			verdeBambino<-id
		} else if(codice==BIANCO){
			fmt.Printf("[Paziente %d, Codice BIANCO] Entro in ambulatorio PEDIATRA\n",id)
			biancoBambino<-id
		}
	}
	<- ack_paziente[id]
	//ho finito il mio lavoro
	done <- true
	
}

func funcTriage(){
	//ciclo di vita del server
	var codice int
	for{
		select{
			case x:=<-triage:
				codice=rand.Intn(4)
				if(codice==ROSSO){
					fmt.Printf("[Triage] Ricevuto paziente %d, assegnato codice ROSSO\n",x)
				} else if(codice==GIALLO){
					fmt.Printf("[Triage] Ricevuto paziente %d, assegnato codice GIALLO\n",x)
				} else if(codice==VERDE){
					fmt.Printf("[Triage] Ricevuto paziente %d, assegnato codice VERDE\n",x)
				} else if(codice==BIANCO){
					fmt.Printf("[Triage] Ricevuto paziente %d, assegnato codice BIANCO\n",x)
				}
				ack_triage[x] <- codice
			case <-terminaTriage:
				fmt.Printf("[Triage] termino\n")
				done<-true
				return
		}
	}
}

func funcAdulti(){
	//ciclo di vita del server
	var medici int = NMEDICI
	for{
		select{
			case x:=<-when(medici>0,rossoAdulti):
				medici--
				fmt.Printf("[Ambulatorio Adulti] Inizio visita paziente %d con codice ROSSO [Medici liberi:%d]\n",x,medici)
				go funcVisitaAdulti(x)
			case x:=<-when(medici>0&&len(rossoAdulti)==0, gialloAdulti):
				medici--
				fmt.Printf("[Ambulatorio Adulti] Inizio visita paziente %d con codice GIALLO [Medici liberi:%d]\n",x,medici)
				go funcVisitaAdulti(x)
			case x:=<-when(medici>0&&len(rossoAdulti)==0&&len(gialloAdulti)==0, verdeAdulti):
				medici--
				fmt.Printf("[Ambulatorio Adulti] Inizio visita paziente %d con codice VERDE [Medici liberi:%d]\n",x,medici)
				go funcVisitaAdulti(x)
			case x:=<-when(medici>0&&len(rossoAdulti)==0&&len(gialloAdulti)==0&&len(verdeAdulti)==0, biancoAdulti):
				medici--
				fmt.Printf("[Ambulatorio Adulti] Inizio visita paziente %d con codice BIANCO [Medici liberi:%d]\n",x,medici)
				go funcVisitaAdulti(x)
			case x:=<-visitaAdulti:
				medici++
				fmt.Printf("[Ambulatorio Adulti] Fine visita paziente %d [Medici liberi:%d]\n",x,medici)
				ack_paziente[x] <- 1
			case <-terminaAdulti:
				fmt.Printf("[Ambulatorio Adulti] termino\n")
				done<-true
				return
		}
	}
}

func funcPediatra(){
	//ciclo di vita del server
	var medici int =NMEDICI
	for{
		select{
			case x:=<-when(medici>0,rossoBambino):
				medici--
				fmt.Printf("[Ambulatorio Pediatra] Inizio visita paziente %d con codice ROSSO [Medici liberi:%d]\n",x,medici)
				go funcVisitaPediatrica(x)
			case x:=<-when(medici>0&&len(rossoBambino)==0, gialloBambino):
				medici--
				fmt.Printf("[Ambulatorio Pediatra] Inizio visita paziente %d con codice GIALLO [Medici liberi:%d]\n",x,medici)
				go funcVisitaPediatrica(x)
			case x:=<-when(medici>0&&len(rossoBambino)==0&&len(gialloBambino)==0, verdeBambino):
				medici--
				fmt.Printf("[Ambulatorio Pediatra] Inizio visita paziente %d con codice VERDE [Medici liberi:%d]\n",x,medici)
				go funcVisitaPediatrica(x)
			case x:=<-when(medici>0&&len(rossoBambino)==0&&len(gialloBambino)==0&&len(verdeBambino)==0, biancoBambino):
				medici--
				fmt.Printf("[Ambulatorio Pediatra] Inizio visita paziente %d con codice BIANCO [Medici liberi:%d]\n",x,medici)
				go funcVisitaPediatrica(x)
			case x:=<-visitaPediatrica:
				medici++
				fmt.Printf("[Ambulatorio Pediatra] Fine visita paziente %d [Medici liberi:%d]\n",x,medici)
				ack_paziente[x] <- 1
			case <-terminaPediatra:
				fmt.Printf("[Ambulatorio Pediatra] termino\n")
				done<-true
				return
		}
	}
}

func funcVisitaAdulti(paziente int){
	time.Sleep(time.Duration( rand.Intn(MAX_ATTESA)+1) * time.Second)
	visitaAdulti<-paziente
	return
}


func funcVisitaPediatrica(paziente int){
	time.Sleep(time.Duration( rand.Intn(MAX_ATTESA)+1) * time.Second)
	visitaPediatrica<-paziente
	return
}


func main(){
	
	fmt.Printf("Programma avviato\n")
	rand.Seed(time.Now().Unix())
	
	//inizializzo canali ack
	for i:=0; i<NPAZIENTI;i++{
		ack_triage[i] = make(chan int, MAX_BUFF)
		ack_paziente[i] = make(chan int, MAX_BUFF)
	}
	
	//lancio threads
	for i:=0; i<NPAZIENTI;i++{
		go funcPaziente(i)
	}
	
	//lancio il server
	go funcTriage()
	go funcPediatra()
	go funcAdulti()
	
	//attendo la terminazione dei clients
	for i:=0; i<NPAZIENTI; i++{
		<-done
	}
	
	
	//avviso il server di terminare
	terminaTriage <- true
	terminaAdulti <- true
	terminaPediatra <- true
	//attendo la terminazione del server
	<-done
	<-done
	<-done
	
	fmt.Printf("Programma terminato\n")
}
