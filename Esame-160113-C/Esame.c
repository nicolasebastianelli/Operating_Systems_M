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

#define NPERSESSO 5
#define TEMPOBALLO 5

time_t  t;
typedef struct{
	int coppieFormate;
	sem_t semUomini[NPERSESSO];
	sem_t semFine;
	int uominiScelti[NPERSESSO];
	pthread_mutex_t lock;
}Formazione;

typedef struct{
	sem_t semInizio;
	sem_t semInizioD[NPERSESSO];
	sem_t semInizioU[NPERSESSO];
	sem_t semFineD[NPERSESSO];
	sem_t semFineU[NPERSESSO];
	sem_t semFine;
	int tempoMusica;
}Esibizione;

typedef struct{
	int coppiePremiate;
	sem_t semInizioPremiaU[NPERSESSO];
	sem_t semInizioPremiaD[NPERSESSO];
	sem_t semFinePremiaU[NPERSESSO];
	sem_t semFinePremiaD[NPERSESSO];
	int coppiaPremiata[NPERSESSO];
}Premiazione;

Formazione formazione;
Esibizione esibizione;
Premiazione premiazione;
int coppia[NPERSESSO];


void *uomo(void *t) {
    int tipo = (intptr_t)t;
    long result=0;
	printf("[Uomo %i]: Aspetto di essere scelto...\n",tipo);
	sem_wait(&formazione.semUomini[tipo]);
	printf("[Uomo %i]: Sono stato scelto, attendo che tutte le coppie vengano scelte\n",tipo);
	sem_post(&formazione.semFine);
	sem_wait(&esibizione.semInizioU[tipo]);
	printf("[Uomo %i]: Inizio a ballare\n",tipo);
	sem_wait(&esibizione.semFineU[tipo]);
	printf("[Uomo %i]: Fine esibizione\n",tipo);
	sem_post(&esibizione.semFine);
	sem_wait(&premiazione.semInizioPremiaU[tipo]);
	printf("[Uomo %i]: Ritiro premio\n",tipo);
	sem_post(&premiazione.semFinePremiaU[tipo]);
    pthread_exit((void*) result);
}

void *donna(void *t) {
    int tipo = (intptr_t)t;
    int uomo;
    long result=0;
    printf("[Donna %i]: Inizializzazione \n",tipo);
    sleep(rand()%5+1);
	while (1) {
		pthread_mutex_lock(&formazione.lock);
		uomo = rand()%NPERSESSO;
		if(formazione.uominiScelti[uomo]==0){
			printf("[Donna %i]: Scelgo uomo %d, attendo che tutte le coppie vengano scelte\n",tipo,uomo);
			formazione.uominiScelti[uomo]=1;
			sem_post(&formazione.semUomini[uomo]);
			formazione.coppieFormate++;
			coppia[tipo]=uomo;
			pthread_mutex_unlock(&formazione.lock);
			if(formazione.coppieFormate==NPERSESSO){
				for(int i =0;i<NPERSESSO;i++){
					sem_wait(&formazione.semFine);
				}
				printf("[Donna %i]: Sono stata l'ultima a scegliere\n",tipo);
				sem_post(&esibizione.semInizio);
			}
			break;
		}
		pthread_mutex_unlock(&formazione.lock);
	}
	
	sem_wait(&esibizione.semInizioD[tipo]);
	printf("[Donna %i]: Inizio a ballare con Uomo %d\n",tipo,uomo);
	sem_wait(&esibizione.semFineD[tipo]);
	printf("[Donna %i]: fine esibizione\n",tipo);
	sem_post(&esibizione.semFine);
	sem_wait(&premiazione.semInizioPremiaD[tipo]);
	printf("[Donna %i]: Ritiro premio\n",tipo);
	sem_post(&premiazione.semFinePremiaD[tipo]);
    pthread_exit((void*) result);
}

void *presidente() {
    long result=0;
    int coppiaP;
    sem_wait(&esibizione.semInizio);
    printf("\n[Presidente]: Tutte le coppie sono state scelte, inizia l'esibizione\n\n");
    for(int i =0;i<NPERSESSO;i++){
    	sem_post(&esibizione.semInizioD[i]);
    	sem_post(&esibizione.semInizioU[i]);
    }
    sleep(esibizione.tempoMusica);
    for(int i =0;i<NPERSESSO;i++){
    	sem_post(&esibizione.semFineD[i]);
		sem_post(&esibizione.semFineU[i]);
    }
    for(int i =0;i<2*NPERSESSO;i++){
        	sem_wait(&esibizione.semFine);
	}
    printf("\n[Presidente]: Fine dell'Esibizione, inizio Premiazione\n\n");
    
    while(1){
    	coppiaP=rand()%NPERSESSO;
    	if(premiazione.coppiaPremiata[coppiaP]==0){
    		premiazione.coppiePremiate++;
    		premiazione.coppiaPremiata[coppiaP]=1;
    		printf("[Presidente]: Premio numero %d la coppia Uomo %d e Donna %d\n",premiazione.coppiePremiate,coppia[coppiaP],coppiaP);
    		sem_post(&premiazione.semInizioPremiaD[coppiaP]);
			sem_post(&premiazione.semInizioPremiaU[coppia[coppiaP]]);
			sem_wait(&premiazione.semFinePremiaD[coppiaP]);
			sem_wait(&premiazione.semFinePremiaU[coppia[coppiaP]]);
    		if(premiazione.coppiePremiate==NPERSESSO){
				printf("\n[Presidente]: Tutte le coppie sono state premiate\n");
				break;
			}
    	}
    }
    pthread_exit((void*) result);
}


int main (int argc, char *argv[]) {
    pthread_t threadUomini[NPERSESSO],threadDonne[NPERSESSO],threadPresidente;
    srand((unsigned) time(&t));
    int rc,i;
    void *status;

    
    // inizializzazione semafori
    formazione.coppieFormate=0;
    sem_init(&formazione.semFine, 0, 0);
    sem_init(&esibizione.semInizio, 0, 0);
    sem_init(&esibizione.semFine, 0, 0);
    esibizione.tempoMusica=TEMPOBALLO;
    premiazione.coppiePremiate=0;
    
    for(i=0;i<NPERSESSO;i++){
       	sem_init(&formazione.semUomini[i], 0, 0);
       	formazione.uominiScelti[i]=0;
       	sem_init(&esibizione.semInizioD[i], 0, 0);
       	sem_init(&esibizione.semInizioU[i], 0, 0);
       	sem_init(&esibizione.semFineD[i], 0, 0);
       	sem_init(&esibizione.semFineU[i], 0, 0);
       	sem_init(&premiazione.semInizioPremiaD[i], 0, 0);
       	sem_init(&premiazione.semInizioPremiaU[i], 0, 0);
       	sem_init(&premiazione.semFinePremiaD[i], 0, 0);
		sem_init(&premiazione.semFinePremiaU[i], 0, 0);
       	premiazione.coppiaPremiata[i]=0;
    }

    //create processi

    for(i=0;i<NPERSESSO;i++){
		rc = pthread_create(&threadUomini[i], NULL, uomo, (void *)(intptr_t)i);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
    }
    for(i=0;i<NPERSESSO;i++){
		rc = pthread_create(&threadDonne[i], NULL, donna, (void *)(intptr_t)i);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
	}
    rc = pthread_create(&threadPresidente, NULL, presidente, NULL);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}

	// join processi

		for(i=0;i<NPERSESSO;i++){
			rc = pthread_join(threadUomini[i], &status);
				if (rc) {
					printf("ERRORE: %d\n", rc);
					exit(-1);
				}
		}
		for(i=0;i<NPERSESSO;i++){
					rc = pthread_join(threadDonne[i], &status);
						if (rc) {
							printf("ERRORE: %d\n", rc);
							exit(-1);
						}
				}

	rc =  pthread_join(threadPresidente, &status);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}

    return 0;
}
