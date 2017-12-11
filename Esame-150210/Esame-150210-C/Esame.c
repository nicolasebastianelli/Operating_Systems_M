//  Created by Nicola Sebastianelli
//
//  Compilare con
//  gcc -D_REENTRANT -std=c99 -o Esame Esame.c -lpthread
//

#include <pthread.h>
#include <semaphore.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <time.h>
#include <math.h>
#include <unistd.h>

#define MAXCONTENITORE 5
#define K 2
#define N 5
#define CLIENTI 2
sem_t stampaMaglie;
sem_t stampaAdesivi;
sem_t stampaBorsa;
sem_t confeziona;
sem_t ordini;

void *funcOrdina(void *t) {
    long result = 0;
    int tipo = (intptr_t)t;

	while (1) {
		printf("[P%i]: Arrivato nuovo ordine, inizio produzione\n",tipo);
		sem_post(&stampaMaglie);
		sem_post(&stampaAdesivi);
		sem_post(&stampaBorsa);

		sem_wait(&ordini);
		printf("[P%i]: Ordine consegnato al cliente\n",tipo);
	}

    pthread_exit((void*) result);
}

void *funcMaglie(void *t) {
    long result = 0;
	int tipo = (intptr_t)t;
	int i=0;
	while(1){
		sem_wait(&stampaMaglie);
		printf("[P%i]: Inizio produzione %d magliette\n", tipo,K);
		for(i=0;i<K;i++){
			sleep(3);
		}
		printf("[P%i]: Fine produzione di tutte le %d maglie richieste\n", tipo,K);
		i=0;
		sem_post(&confeziona);
	}

    pthread_exit((void*) result);
}

void *funcAdesivi(void *t) {
    long result = 0;
	int tipo = (intptr_t)t;
	int i=0;
	while(1){
			sem_wait(&stampaAdesivi);
			printf("[P%i]: Inizio produzione %d adesivi\n", tipo,N);
			for(i=0;i<N;i++){
				sleep(1);
			}
			printf("[P%i]: Fine produzione di tutti i %d adesivi richiesti\n", tipo,N);
			i=0;
			sem_post(&confeziona);
		}
    pthread_exit((void*) result);
}

void *funcBorsa(void *t) {
    long result = 0;
    	int tipo = (intptr_t)t;
    	int i=0;
    	while(1){
			sem_wait(&stampaBorsa);

			printf("[P%i]: Inizio produzione borsa\n", tipo);
			sleep(5);
			printf("[P%i]: Fine produzione borsa\n", tipo);
			sem_post(&confeziona);
    		}
    pthread_exit((void*) result);
}

void *funcConfeziona(void *t) {
    long result = 0;
    int tipo = (intptr_t)t;

	while (1) {
		sem_wait(&confeziona);

		sem_wait(&confeziona);

		sem_wait(&confeziona);
		printf("[P%i]: Inizio confezionamento\n",tipo);
		sleep(5);
		printf("[P%i]: Kit realizzato\n",tipo);

		sem_post(&ordini);
	}

    pthread_exit((void*) result);
}


int main (int argc, char *argv[]) {
    pthread_t threadMaglie,threadAdesivi,threadBorsa,threadConfezionamento,threadOrdina;
    int rc;

    void *status;

    sem_init(&stampaMaglie, 0, 0);
    	sem_init(&stampaAdesivi, 0, 0);
    	sem_init(&stampaBorsa, 0, 0);
    	sem_init(&confeziona, 0, 0);
    	sem_init(&ordini, 0, 0);

	rc = pthread_create(&threadOrdina, NULL, funcOrdina, (void *)(intptr_t)0);
	if (rc) {
		printf("ERRORE: %d\n", rc);
		exit(-1);
	}
	rc = pthread_create(&threadMaglie, NULL, funcMaglie, (void *)(intptr_t)1);
	if (rc) {
		printf("ERRORE: %d\n", rc);
		exit(-1);
	}
	rc = pthread_create(&threadAdesivi, NULL, funcAdesivi, (void *)(intptr_t)2);
	if (rc) {
		printf("ERRORE: %d\n", rc);
		exit(-1);
	}
	rc = pthread_create(&threadBorsa, NULL, funcBorsa, (void *)(intptr_t)3);
	if (rc) {
		printf("ERRORE: %d\n", rc);
		exit(-1);
	}
	rc = pthread_create(&threadConfezionamento, NULL, funcConfeziona, (void *)(intptr_t)4);
	if (rc) {
		printf("ERRORE: %d\n", rc);
		exit(-1);
	}



	rc = pthread_join(threadOrdina, &status);
	if (rc) {
		printf("ERRORE join thread ordini codice %d\n", rc);
	}
	rc = pthread_join(threadMaglie, &status);
	if (rc) {
		printf("ERRORE join thread maglie codice %d\n", rc);
	}
	rc = pthread_join(threadAdesivi, &status);
	if (rc) {
		printf("ERRORE join thread adesivi codice %d\n", rc);
	}
	rc = pthread_join(threadBorsa, &status);
	if (rc) {
		printf("ERRORE join thread borsa codice %d\n", rc);
	}
	rc = pthread_join(threadConfezionamento, &status);
	if (rc) {
		printf("ERRORE join thread confezionamento codice %d\n", rc);
	}

    return 0;
}
