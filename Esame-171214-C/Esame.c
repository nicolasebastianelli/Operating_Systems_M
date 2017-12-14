//  Created by Nicola Sebastianelli
//	s0000850827
//
//  Compilare con
//  gcc -D_REENTRANT -o Esame Esame.c -lpthread
//

#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <math.h>
#include <unistd.h>

#define NMEDICI 5
#define NPAZIENTI 15
#define ADULTO 0
#define BAMBINO 1
#define ROSSO 0
#define GIALLO 1
#define VERDE 2
#define BIANCO 3

typedef struct{
	int eta;
	int codice;
}Paziente;

typedef struct{
	sem_t blocco;
	sem_t generaCodice[NPAZIENTI];
	sem_t aspettaCodice[NPAZIENTI];
	int pazienteServito;
	pthread_mutex_t lock;
}Triage;

typedef struct{
	sem_t rosso;
	sem_t giallo;
	sem_t verde;
	sem_t bianco;
	sem_t fineVisita;
	int rossoInCoda;
	int gialloInCoda;
	int verdeInCoda;
	int biancoInCoda;
	int medici;
}Ambulatorio;


Triage triage;
Ambulatorio ambulatorioBambini;
Ambulatorio ambulatorioAdulti;
Paziente paziente[NPAZIENTI];
time_t  t;
int fine=0;

void *funcPaziente(void *t) {
    int tipo = (intptr_t)t;
    long result=0;
    paziente[tipo].eta=rand()%2;
    
    if(paziente[tipo].eta==ADULTO){
    	printf("[Paziente %i ADULTO]: Inizializzazione \n",tipo);
    } else if(paziente[tipo].eta==BAMBINO){
    	printf("[Paziente %i BAMBINO]: Inizializzazione \n",tipo);
    }
	sleep(rand()%5+1);
	printf("[Paziente %i]: Richiedo codice\n",tipo);
	
	pthread_mutex_lock(&triage.lock);
	triage.pazienteServito=tipo;
	sem_post(&triage.blocco);	
	sem_post(&triage.generaCodice[tipo]);	
	sem_wait(&triage.aspettaCodice[tipo]);
	pthread_mutex_unlock(&triage.lock);
	
	
	if(paziente[tipo].eta==ADULTO){
		if(paziente[tipo].codice==ROSSO){
				printf("[Paziente %i ADULTO]: codice assegnato ROSSO,vado all'ambulatorio\n",tipo);
				sleep(1);
				ambulatorioAdulti.rossoInCoda++;
				sem_wait(&ambulatorioAdulti.rosso);
				printf("[Paziente %i ADULTO]: inizio visita con codice ROSSO, medici disponibili: %d\n",tipo,ambulatorioAdulti.medici);
				sleep(rand()%3+1);
				ambulatorioAdulti.medici++;
				printf("[Paziente %i ADULTO]: fine visita con codice ROSSO, medici disponibili: %d\n",tipo,ambulatorioAdulti.medici);
				sem_post(&ambulatorioAdulti.fineVisita);
			}else if(paziente[tipo].codice==GIALLO){
				printf("[Paziente %i ADULTO]: codice assegnato GIALLO,vado all'ambulatorio\n",tipo);
				sleep(1);
				ambulatorioAdulti.gialloInCoda++;
				sem_wait(&ambulatorioAdulti.giallo);
				printf("[Paziente %i ADULTO]: inizio visita con codice GIALLO, medici disponibili: %d\n",tipo,ambulatorioAdulti.medici);
				sleep(rand()%3+1);
				ambulatorioAdulti.medici++;
				printf("[Paziente %i ADULTO]: fine visita con codice GIALLO, medici disponibili: %d\n",tipo,ambulatorioAdulti.medici);	
				sem_post(&ambulatorioAdulti.fineVisita);
			}else if(paziente[tipo].codice==VERDE){
				printf("[Paziente %i ADULTO]: codice assegnato VERDE,vado all'ambulatorio\n",tipo);
				sleep(1);
				ambulatorioAdulti.verdeInCoda++;
				sem_wait(&ambulatorioAdulti.verde);
				printf("[Paziente %i ADULTO]: inizio visita con codice VERDE, medici disponibili: %d\n",tipo,ambulatorioAdulti.medici);
				sleep(rand()%3+1);
				ambulatorioAdulti.medici++;
				printf("[Paziente %i ADULTO]: fine visita con codice VERDE, medici disponibili: %d\n",tipo,ambulatorioAdulti.medici);
				sem_post(&ambulatorioAdulti.fineVisita);
			}else if(paziente[tipo].codice==BIANCO){
				printf("[Paziente %i ADULTO]: codice assegnato BIANCO,vado all'ambulatorio\n",tipo);
				sleep(1);
				ambulatorioAdulti.biancoInCoda++;
				sem_wait(&ambulatorioAdulti.bianco);
				printf("[Paziente %i ADULTO]: inizio visita con codice BIANCO, medici disponibili: %d\n",tipo,ambulatorioAdulti.medici);
				sleep(rand()%3+1);
				ambulatorioAdulti.medici++;
				printf("[Paziente %i ADULTO]: fine visita con codice BIANCO, medici disponibili: %d\n",tipo,ambulatorioAdulti.medici);
				sem_post(&ambulatorioAdulti.fineVisita);
			}
	} else if(paziente[tipo].eta==BAMBINO){
		if(paziente[tipo].codice==ROSSO){
				printf("[Paziente %i BAMBINO]: codice assegnato ROSSO,vado all'ambulatorio\n",tipo);
				sleep(1);
				ambulatorioBambini.rossoInCoda++;
				sem_wait(&ambulatorioBambini.rosso);
				printf("[Paziente %i BAMBINO]: inizio visita con codice ROSSO, medici disponibili: %d\n",tipo,ambulatorioBambini.medici);
				sleep(rand()%3+1);
				ambulatorioBambini.medici++;
				printf("[Paziente %i BAMBINO]: fine visita con codice ROSSO, medici disponibili: %d\n",tipo,ambulatorioBambini.medici);
				sem_post(&ambulatorioBambini.fineVisita);
			}else if(paziente[tipo].codice==GIALLO){
				printf("[Paziente %i BAMBINO]: codice assegnato GIALLO,vado all'ambulatorio\n",tipo);
				sleep(1);
				ambulatorioBambini.gialloInCoda++;
				sem_wait(&ambulatorioBambini.giallo);
				printf("[Paziente %i BAMBINO]: inizio visita con codice GIALLO, medici disponibili: %d\n",tipo,ambulatorioBambini.medici);
				sleep(rand()%3+1);
				ambulatorioBambini.medici++;
				printf("[Paziente %i BAMBINO]: fine visita con codice GIALLO, medici disponibili: %d\n",tipo,ambulatorioBambini.medici);
				sem_post(&ambulatorioBambini.fineVisita);
			}else if(paziente[tipo].codice==VERDE){
				printf("[Paziente %i BAMBINO]: codice assegnato VERDE,vado all'ambulatorio\n",tipo);
				sleep(1);
				ambulatorioBambini.verdeInCoda++;
				sem_wait(&ambulatorioBambini.verde);
				printf("[Paziente %i BAMBINO]: inizio visita con codice VERDE, medici disponibili: %d\n",tipo,ambulatorioBambini.medici);
				sleep(rand()%3+1);
				ambulatorioBambini.medici++;
				printf("[Paziente %i BAMBINO]: fine visita con codice VERDE, medici disponibili: %d\n",tipo,ambulatorioBambini.medici);
				sem_post(&ambulatorioBambini.fineVisita);
			}else if(paziente[tipo].codice==BIANCO){
				printf("[Paziente %i BAMBINO]: codice assegnato BIANCO,vado all'ambulatorio\n",tipo);
				sleep(1);
				ambulatorioBambini.biancoInCoda++;
				sem_wait(&ambulatorioBambini.bianco);
				printf("[Paziente %i BAMBINO]: inizio visita con codice BIANCO, medici disponibili: %d\n",tipo,ambulatorioBambini.medici);
				sleep(rand()%3+1);
				ambulatorioBambini.medici++;
				printf("[Paziente %i BAMBINO]: fine visita con codice BIANCO, medici disponibili: %d\n",tipo,ambulatorioBambini.medici);
				sem_post(&ambulatorioBambini.fineVisita);
			}
	}
    pthread_exit((void*) result);
}

void *funcTriage() {
    long result=0;

	for(int i =0;i<NPAZIENTI;i++){
		sem_wait(&triage.blocco);
		sem_wait(&triage.generaCodice[triage.pazienteServito]);	
		sleep(1);
		paziente[triage.pazienteServito].codice = rand()%4;
		sem_post(&triage.aspettaCodice[triage.pazienteServito]);
	}
	printf("[Triage]: tutti i pazienti sono stati serviti, TERMINO\n");
    pthread_exit((void*) result);
}

void *funcPediatra() {
    long result=0;

	while (1) {
		if(fine!=1){
			if(ambulatorioBambini.medici!=0){
				if(ambulatorioBambini.rossoInCoda!=0){
					ambulatorioBambini.rossoInCoda--;
					ambulatorioBambini.medici--;
					sem_post(&ambulatorioBambini.rosso);
				}else if(ambulatorioBambini.gialloInCoda!=0){
					ambulatorioBambini.gialloInCoda--;
					ambulatorioBambini.medici--;
					sem_post(&ambulatorioBambini.giallo);
				}else if(ambulatorioBambini.verdeInCoda!=0){
					ambulatorioBambini.verdeInCoda--;
					ambulatorioBambini.medici--;
					sem_post(&ambulatorioBambini.verde);
				}else if(ambulatorioBambini.biancoInCoda!=0){
					ambulatorioBambini.biancoInCoda--;
					ambulatorioBambini.medici--;
					sem_post(&ambulatorioBambini.bianco);
				}
			}
			else if(ambulatorioBambini.medici==0){
				sem_wait(&ambulatorioBambini.fineVisita);
			}	
			
		} else{
			printf("[Pediatra]: tutti i pazienti sono stati serviti, TERMINO\n");
			break;
		}
	}

    pthread_exit((void*) result);
}

void *funcAdulti() {
    long result=0;

	while (1) {
		if(fine!=1){
			if(ambulatorioAdulti.medici!=0){
				if(ambulatorioAdulti.rossoInCoda!=0){
					ambulatorioAdulti.rossoInCoda--;
					ambulatorioAdulti.medici--;
					sem_post(&ambulatorioAdulti.rosso);
				}else if(ambulatorioAdulti.gialloInCoda!=0){
					ambulatorioAdulti.gialloInCoda--;
					ambulatorioAdulti.medici--;
					sem_post(&ambulatorioAdulti.giallo);
				}else if(ambulatorioAdulti.verdeInCoda!=0){
					ambulatorioAdulti.verdeInCoda--;
					ambulatorioAdulti.medici--;
					sem_post(&ambulatorioAdulti.verde);
				}else if(ambulatorioAdulti.biancoInCoda!=0){
					ambulatorioAdulti.biancoInCoda--;
					ambulatorioAdulti.medici--;
					sem_post(&ambulatorioAdulti.bianco);
				}
			}
			else if(ambulatorioAdulti.medici==0){
				sem_wait(&ambulatorioBambini.fineVisita);
			}	
		}else{
			printf("[Adulti]: tutti i pazienti sono stati serviti, TERMINO\n");
			break;
		}
	}

    pthread_exit((void*) result);
}


int main (int argc, char *argv[]) {
    pthread_t threadPazienti[NPAZIENTI],threadTriage,threadAdulti,threadPediatra;
    srand((unsigned) time(&t));

    int rc,i;

    void *status;
    pthread_mutex_init(&triage.lock,NULL);
    sem_init(&triage.blocco, 0, 0);
    // inizializzazione semafori
    for(i=0;i<NPAZIENTI;i++){
    	sem_init(&triage.generaCodice[i], 0, 0);
		sem_init(&triage.aspettaCodice[i], 0, 0);
    }
    sem_init(&ambulatorioBambini.rosso, 0, 0);
    sem_init(&ambulatorioBambini.giallo, 0, 0);
    sem_init(&ambulatorioBambini.verde, 0, 0);
    sem_init(&ambulatorioBambini.bianco, 0, 0);
    sem_init(&ambulatorioAdulti.rosso, 0, 0);
    sem_init(&ambulatorioAdulti.giallo, 0, 0);
    sem_init(&ambulatorioAdulti.verde, 0, 0);
    sem_init(&ambulatorioAdulti.bianco, 0, 0);
    ambulatorioBambini.rossoInCoda=0;
    ambulatorioBambini.gialloInCoda=0;
    ambulatorioBambini.verdeInCoda=0;
    ambulatorioBambini.biancoInCoda=0;
    ambulatorioBambini.medici=NMEDICI;
    ambulatorioAdulti.rossoInCoda=0;
    ambulatorioAdulti.gialloInCoda=0;
    ambulatorioAdulti.verdeInCoda=0;
    ambulatorioAdulti.biancoInCoda=0;
    ambulatorioAdulti.medici=NMEDICI;
        
    //create processi

    for(i=0;i<NPAZIENTI;i++){
		rc = pthread_create(&threadPazienti[i], NULL, funcPaziente, (void *)(intptr_t)i);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
    }
    rc = pthread_create(&threadTriage, NULL, funcTriage, NULL);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}
	rc = pthread_create(&threadAdulti, NULL, funcAdulti, NULL);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}
	rc = pthread_create(&threadPediatra, NULL, funcPediatra, NULL);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}

	// join processi

	for(i=0;i<NPAZIENTI;i++){
			rc = pthread_join(threadPazienti[i], &status);
				if (rc) {
					printf("ERRORE: %d\n", rc);
					exit(-1);
				}
		}
	fine=1;
	rc =  pthread_join(threadTriage, &status);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}
	rc =  pthread_join(threadAdulti, &status);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
	rc =  pthread_join(threadPediatra, &status);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
    return 0;
}
