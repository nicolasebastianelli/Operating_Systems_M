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

//costanti
const NUTENTI int = 20
const MAXPOSTI int = 5
const ORDINARIO int = 0
const STUDENTE int = 1
const NONABBONATO int = 0
const ABBONATO int = 1
const MAX_ATTESA int = 10
const MAX_BUFF int = 20

//canali
var done = make(chan bool)
var terminaPiscina = make(chan bool)
var terminaBiglietteria = make(chan bool)
var biglietto = make(chan int, MAX_BUFF)
var ordinarioAbbonato = make(chan int, MAX_BUFF)
var ordinarioNonAbbonato = make(chan int, MAX_BUFF)
var studenteAbbonato = make(chan int, MAX_BUFF)
var studenteNonAbbonato = make(chan int, MAX_BUFF)
var uscitaStudenti = make(chan int, MAX_BUFF)
var uscitaOrdinari = make(chan int, MAX_BUFF)
var restituzioneChiave = make(chan int, MAX_BUFF)

//ack
var ack_biglietto[NUTENTI] chan int
var ack_entrataPiscina[NUTENTI] chan int
var ack_uscitaPiscina[NUTENTI] chan int
var ack_restituzione[NUTENTI] chan int

//funzioni utility
func when(b bool, c chan int) chan int{
	if(!b){
		return nil
	}
	return c
}


func utente(id int){
	
	//Inizializzazione Utente
	fmt.Printf("Inizializzazione Utente %d\n",id)
	var tempo_di_attesa int
	var abbonamento int =rand.Intn(2)
	var tipoUtente int =rand.Intn(2)
	
	//1) Acquisto Biglietto
	if(abbonamento==ABBONATO){
		if(tipoUtente==STUDENTE){
			fmt.Printf("[Utente %d STUDENTE ABBONATO] Acquisto biglietto\n",id)
		}else if(tipoUtente==ORDINARIO){
			fmt.Printf("[Utente %d ORDINARIO ABBONATO] Acquisto biglietto\n",id)
		}
	}else if (abbonamento==NONABBONATO){
		if(tipoUtente==STUDENTE){
			fmt.Printf("[Utente %d STUDENTE NON ABBONATO] Acquisto biglietto\n",id)
		}else if(tipoUtente==ORDINARIO){
			fmt.Printf("[Utente %d ORDINARIO NON ABBONATO] Acquisto biglietto\n",id)
		}
	}
	biglietto <- id
	<- ack_biglietto[id]
	fmt.Printf("[Utente %d] Biglietto acquistato, ricevuta chiave\n",id)
	
	
	//2) Entrata in piscina
	fmt.Printf("[Utente %d] Richiesta entrata in piscina\n",id)
	if(abbonamento==ABBONATO){
		if(tipoUtente==STUDENTE){
			studenteAbbonato<-id
		}else if(tipoUtente==ORDINARIO){
			ordinarioAbbonato<-id
		}
	}else if (abbonamento==NONABBONATO){
		if(tipoUtente==STUDENTE){
			studenteNonAbbonato<-id
		}else if(tipoUtente==ORDINARIO){
			ordinarioNonAbbonato<-id
		}
	}
	<-ack_entrataPiscina[id]
	
	//3) Trascorro un tempo arbitrario in piscina
	tempo_di_attesa = rand.Intn(MAX_ATTESA)+1 //+1 perchè randomizza da 0 a MAX_ATTESA
	time.Sleep(time.Duration(tempo_di_attesa) * time.Second)
	
	//4)Uscita Piscina
	fmt.Printf("[Utente %d] Esco dalla piscina\n",id)
	if(tipoUtente==STUDENTE){
			uscitaStudenti<-id
	}else if(tipoUtente==ORDINARIO){
		uscitaOrdinari<-id
	}
	<-ack_uscitaPiscina[id]
	
	//5)Restituisco chiave
	fmt.Printf("[Utente %d] Restituisco chiave alla biglietteria\n",id)
	restituzioneChiave<-id
	<-ack_restituzione[id]
	
	
	// Termino Utente
	fmt.Printf("Termino Utente %d\n",id)
	done <- true
}

func biglietteria(){
	var nChiavi int =0
	for{
		select{
			case x:=<-biglietto:
				nChiavi ++
				fmt.Printf("[Biglietteria] Vendo biglietto all' utente %d e consegno la chiave, [Chiavi consegnate: %d]\n",x,nChiavi)
				ack_biglietto[x] <- 1
			case x:=<-restituzioneChiave:
				fmt.Printf("[Biglietteria] Restituita chiave dall'utente %d\n",x)
				ack_restituzione[x] <- 1
			case <-terminaBiglietteria:
				fmt.Printf("Termino Biglietteria\n")
				done<-true
				return
		}
	}
}


func piscina(){
	var postiOccupati int=0
	var totEntrati int=0
	var nStudenti int=0
	var nOrdinari int=0
	for{
		select{
			case x:=<-when((nStudenti<=nOrdinari  || len(ordinarioAbbonato)==0 && len(ordinarioNonAbbonato)==0) && postiOccupati<MAXPOSTI, studenteAbbonato):
				nStudenti++
				postiOccupati++
				totEntrati++
				fmt.Printf("[Piscina] Entrato Utente STUDENTE ABBONATO %d [Strudenti: %d, Ordinari: %d, TotaleDentro:%d, TotaleEntrati %d]\n",x,nStudenti,nOrdinari,postiOccupati,totEntrati)
				ack_entrataPiscina[x] <- 1
			case x:=<-when((nStudenti<=nOrdinari || len(ordinarioAbbonato)==0 && len(ordinarioNonAbbonato)==0) && len(studenteAbbonato)==0 && postiOccupati<MAXPOSTI, studenteNonAbbonato):
				nStudenti++
				postiOccupati++
				totEntrati++
				fmt.Printf("[Piscina] Entrato Utente STUDENTE NON ABBONATO %d [Strudenti: %d, Ordinari: %d, TotaleDentro:%d, TotaleEntrati %d]\n",x,nStudenti,nOrdinari,postiOccupati,totEntrati)
				ack_entrataPiscina[x] <- 1
			case x:=<-when((nOrdinari<nStudenti || len(studenteAbbonato)==0 && len(studenteNonAbbonato)==0) && postiOccupati<MAXPOSTI, ordinarioAbbonato):
				nOrdinari++
				postiOccupati++
				totEntrati++
				fmt.Printf("[Piscina] Entrato Utente ORDINARIO ABBONATO %d [Strudenti: %d, Ordinari: %d, TotaleDentro:%d, TotaleEntrati %d]\n",x,nStudenti,nOrdinari,postiOccupati,totEntrati)
				ack_entrataPiscina[x] <- 1
			case x:=<-when((nOrdinari<nStudenti || len(studenteAbbonato)==0 && len(studenteNonAbbonato)==0) && len(ordinarioAbbonato)==0 && postiOccupati<MAXPOSTI, ordinarioNonAbbonato):
				nOrdinari++
				postiOccupati++
				totEntrati++
				fmt.Printf("[Piscina] Entrato Utente ORDINARIO NON ABBONATO %d [Strudenti: %d, Ordinari: %d, TotaleDentro:%d, TotaleEntrati %d]\n",x,nStudenti,nOrdinari,postiOccupati,totEntrati)
				ack_entrataPiscina[x] <- 1
			case x:=<-uscitaOrdinari:
				nOrdinari--
				postiOccupati--
				fmt.Printf("[Piscina] Uscito Utente ORDINARIO %d [Strudenti: %d, Ordinari: %d, TotaleDentro:%d, TotaleEntrati %d]\n",x,nStudenti,nOrdinari,postiOccupati,totEntrati)
				ack_uscitaPiscina[x] <- 1
			case x:=<-uscitaStudenti:
				nStudenti--
				postiOccupati--
				fmt.Printf("[Piscina] Uscito Utente STUDENTE %d [Strudenti: %d, Ordinari: %d, TotaleDentro:%d, TotaleEntrati %d]\n",x,nStudenti,nOrdinari,postiOccupati,totEntrati)
				ack_uscitaPiscina[x] <- 1
			case <-terminaPiscina:
				fmt.Printf("Termino Piscina\n")
				done<-true
				return
		}
	}
}

func main(){
	
	fmt.Printf("Programma avviato\n")
	rand.Seed(time.Now().Unix())
	fmt.Printf("NUTENTI: %d, MAXPOSTI: %d, MAX_ATTESA: %d,MAX_BUFF: %d\n\n",NUTENTI,MAXPOSTI,MAX_ATTESA,MAX_BUFF)
	
	//inizializzo canali ack
	for i:=0; i<NUTENTI;i++{
		ack_biglietto[i] = make(chan int, MAX_BUFF)
		ack_entrataPiscina[i] = make(chan int, MAX_BUFF)
		ack_uscitaPiscina[i] = make(chan int, MAX_BUFF)
		ack_restituzione[i] = make(chan int, MAX_BUFF)
	}
	
	//lancio goroutine
	for i:=0; i<NUTENTI;i++{
		go utente(i)
	}
	
	go biglietteria()
	go piscina()
	
	//attendo la terminazione degli utenti
	for i:=0; i<NUTENTI; i++{
		<-done
	}
	
	//termino biglietteria
	terminaBiglietteria <-true
	<- done
	
	//termino piscina
	terminaPiscina <- true
	<-done
	
	fmt.Printf("Programma terminato\n")
}
