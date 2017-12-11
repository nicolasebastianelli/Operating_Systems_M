//  Created by Nicola Sebastianelli
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
#define PADRE 0
#define MADRE 1
#define FIGLIA 2

typedef struct{
	sem_t prodotti;
	sem_t daprodurre;
}Contenitore;

typedef struct{
	sem_t produzione;
	sem_t vestizione;
}Realizzazione;

Contenitore contenitore[3];
Realizzazione realizzazione[3];

void *funcRealizzazione(void *t) {
    long result = 0;
    int tipo = (intptr_t)t;
    int res;
	while (1) {
		printf("[Realizzazione %i]: Inizio realizzazione\n",tipo);
		sleep(tipo+5);
		printf("[Realizzazione %i]: Fine realizzazione\n",tipo);
		sem_post(&realizzazione[tipo].produzione);
		sem_wait(&realizzazione[tipo].vestizione);
	}

    pthread_exit((void*) result);
}

void *funcVestizione(void *t) {
    long result = 0;
	int tipo = (intptr_t)t;
	int res;
	while(1){

		sem_wait(&realizzazione[tipo].produzione);
		sem_post(&realizzazione[tipo].vestizione);
		printf("[Vestizione %i]: Inizio vestizione\n",tipo);
		sleep(tipo+1);
		printf("[Vestizione %i]: Fine vestizione\n",tipo);
		sem_wait(&contenitore[tipo].daprodurre);
		sem_post(&contenitore[tipo].prodotti);
		printf("[Vestizione %i]: Bambola inserita nel contenitore\n",tipo);
	}

    pthread_exit((void*) result);
}


void *funcConfeziona(void *t) {
    long result = 0;
    int tipo = (intptr_t)t;

	while (1) {
		sem_wait(&contenitore[PADRE].prodotti);
		sem_post(&contenitore[PADRE].daprodurre);
		printf("[Confezionamento]: Raccolta bambola padre\n");
		sem_wait(&contenitore[MADRE].prodotti);
		sem_post(&contenitore[MADRE].daprodurre);
		printf("[Confezionamento]: Raccolta bambola madre\n");
		sem_wait(&contenitore[FIGLIA].prodotti);
		sem_post(&contenitore[FIGLIA].daprodurre);
		printf("[Confezionamento]: Raccolta bambola figlia\n");

		printf("[Confezionamento]: Nuova confezione realizzata\n");
	}

    pthread_exit((void*) result);
}


int main (int argc, char *argv[]) {
    pthread_t threadRealizzazione[3],threadVestizione[3],threadConfezionamento;
    int rc,i;

    void *status;

    for(i=PADRE;i<=FIGLIA;i++){
    	Contenitore _contenitore;
    	Realizzazione _realizzazione;

    	contenitore[i] = _contenitore;
    	realizzazione[i] = _realizzazione;

    	sem_init(&realizzazione[i].produzione, 0, 0);
		sem_init(&realizzazione[i].vestizione, 0, 0);
		sem_init(&contenitore[i].prodotti, 0, 0);
		sem_init(&contenitore[i].daprodurre, 0, N);
    }

    for(i=PADRE;i<=FIGLIA;i++){
		rc = pthread_create(&threadRealizzazione[i], NULL, funcRealizzazione, (void *)(intptr_t)i);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
    }

    for(i=PADRE;i<=FIGLIA;i++){
		rc = pthread_create(&threadVestizione[i], NULL, funcVestizione, (void *)(intptr_t)i);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
	}

    rc = pthread_create(&threadConfezionamento, NULL, funcConfeziona, (void *)(intptr_t)i);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}

	for(i=PADRE;i<=FIGLIA;i++){
			rc = pthread_join(threadRealizzazione[i], &status);
				if (rc) {
					printf("ERRORE: %d\n", rc);
					exit(-1);
				}
		}

	for(i=PADRE;i<=FIGLIA;i++){
		rc =  pthread_join(threadVestizione[i], &status);
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
