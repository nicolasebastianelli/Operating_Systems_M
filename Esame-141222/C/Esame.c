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

#define N 5

sem_t semDolci[N];
sem_t semConfezionamento;
time_t  t;

void *funcProduzione(void *t) {
    int tipo = (intptr_t)t;
    long result =0;

	while (1) {
		printf("[P%i]: Inizio produzione dolce %d\n",tipo,tipo);
		sleep(rand()%5+1);
		printf("[P%i]: Fine produzione dolce, metto nella cesta\n",tipo);
		sem_post(&semDolci[tipo]);
		sem_wait(&semConfezionamento);
	}

    pthread_exit((void*) result);
}

void *funcConfezionamento() {
    long result =0;
    int i;
	while (1) {
		for(i=0;i<N;i++){
				sem_wait(&semDolci[i]);
				printf("[PC]: Ricevuto dolce %d\n",i);				
		}
		printf("[PC]: Cesta piena, Inizio confezionamento\n");
		sleep(rand()%5+1);
		printf("[PC]: Fine confezionamento, Inizio creazione nuova cesta\n");
		for(i=0;i<N;i++){
			sem_post(&semConfezionamento);			
		}
	}

    pthread_exit((void*) result);
}


int main (int argc, char *argv[]) {
    pthread_t threadProduzione[N],threadConfezionamento;
    srand((unsigned) time(&t));

    int rc,i;

    void *status;

    // inizializzazione semafori
    for(i=0;i<N;i++){
    	sem_init(&semDolci[i], 0, 0);
		
    }
    sem_init(&semConfezionamento, 0, 0);

    //create processi

    for(i=0;i<N;i++){
		rc = pthread_create(&threadProduzione[i], NULL, funcProduzione, (void *)(intptr_t)i);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
    }
    rc = pthread_create(&threadConfezionamento, NULL, funcConfezionamento, NULL);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}

	// join processi

	for(i=0;i<N;i++){
			rc = pthread_join(threadProduzione[i], &status);
				if (rc) {
					printf("ERRORE: %d\n", rc);
					exit(-1);
				}
		}

	rc =  pthread_join(threadConfezionamento, &status);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}

    return 0;
}
