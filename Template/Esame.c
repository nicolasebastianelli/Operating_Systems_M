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

typedef struct{
	sem_t semaforo1;
	sem_t semaforo2;
}Struttura;

Struttura struttura[3];
time_t  t;

void *funcThread1(void *t) {
    int tipo = (intptr_t)t;
    long result=0;

	while (1) {
		printf("[P %i]: Inizio \n",tipo);
		sleep(rand()%5+1);
		printf("[P %i]: Fine \n",tipo);
		sem_post(&struttura[tipo].semaforo1);
		sem_wait(&struttura[tipo].semaforo2);
	}

    pthread_exit((void*) result);
}

void *funcThread2(void *t) {
    int tipo = (intptr_t)t;
    long result=0;

	while (1) {
		sem_wait(&struttura[0].semaforo1);
		sem_post(&struttura[0].semaforo2);
		printf("[P]: Faccio cose\n");
		sem_wait(&struttura[1].semaforo1);
		sem_post(&struttura[1].semaforo2);

		printf("[Confezionamento]: Nuova confezione realizzata\n");
	}

    pthread_exit((void*) result);
}


int main (int argc, char *argv[]) {
    pthread_t thread1[3],thread2;
    srand((unsigned) time(&t));

    int rc,i;

    void *status;

    // inizializzazione semafori
    for(i=0;i<=2;i++){
    	Struttura _struttura;

    	struttura[i] = _struttura;

    	sem_init(&struttura[i].semaforo1, 0, 0);
		sem_init(&struttura[i].semaforo2, 0, 0);
    }

    //create processi

    for(i=0;i<=2;i++){
		rc = pthread_create(&thread1[i], NULL, funcThread1, (void *)(intptr_t)i);
			if (rc) {
				printf("ERRORE: %d\n", rc);
				exit(-1);
			}
    }
    rc = pthread_create(&thread2, NULL, funcThread2, (void *)(intptr_t)0);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}

	// join processi

	for(i=0;i<=2;i++){
			rc = pthread_join(thread1[i], &status);
				if (rc) {
					printf("ERRORE: %d\n", rc);
					exit(-1);
				}
		}

	rc =  pthread_join(thread2, &status);
		if (rc) {
			printf("ERRORE: %d\n", rc);
			exit(-1);
		}

    return 0;
}
