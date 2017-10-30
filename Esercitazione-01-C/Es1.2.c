/*
 * main.c
 *
 *  Created on: Oct 23, 2017
 *      Author: s0000718463
 */

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#define NUMBER_OF_FILM 10
#define NUMBER_OF_PEOPLE 3

typedef struct Film {
	int somma_voti;
	int votanti;
} Film;

struct Film arrayFilm[NUMBER_OF_FILM];

pthread_mutex_t mutex_threads[NUMBER_OF_FILM];

void *ThreadCode(void *t) {

	long tid;
	long result = 1;
	int i;
	int voto;

	tid = (int)t;

	for (i = 0; i < NUMBER_OF_FILM; i++) {
		pthread_mutex_lock(&mutex_threads[i]);
		voto = rand() % 10 + 1;
		arrayFilm[i].somma_voti += voto;
		arrayFilm[i].votanti++;
		printf("Utente %ld. Film %d - Voto %d\n", tid, i, voto);
		pthread_mutex_unlock(&mutex_threads[i]);
	}

	pthread_exit((void*) result);

}

int main (int argc, char *argv[]) {

	int i, rc, bestfilm;
	long t;
	long result;
	float voto, bestvoto = 0;
	pthread_t threads[NUMBER_OF_PEOPLE];

	srand(time(NULL));
	for (i = 0; i < NUMBER_OF_FILM; i++) {
		pthread_mutex_init(&mutex_threads[i], NULL);
		arrayFilm[i] = (Film){.somma_voti = 0, .votanti = 0};
	}

	for (t = 0; t < NUMBER_OF_PEOPLE; t++) {

		printf("Main: Intervistato %ld.\n", t);
		rc = pthread_create(&threads[t], NULL, ThreadCode, (void*)t);

		if (rc) {
			printf ("ERRORE: %d\n", rc);
			exit(-1);
		}
	}

	for (t = 0; t < NUMBER_OF_PEOPLE; t++) {

		rc = pthread_join(threads[t], (void *)&result);

		if (rc) {
			printf ("ERRORE: join thread %ld codice %d.\n", t, rc);
		} else {
			printf("Main: Finito intervistato %ld.\n", t);
		}
	}

	printf("\n\nMain: Risultati finali.\n");
	for (i = 0; i < NUMBER_OF_FILM; i++) {
		voto = arrayFilm[i].somma_voti/((double)arrayFilm[i].votanti);
		printf("Film %d: %f\n", i, voto);
		if(voto> bestvoto)
		{
			bestvoto = voto;
			bestfilm=i;
		}
	}
	printf("\nMiglior film %d con voto %f\n\n",bestfilm,bestvoto);
	printf("Main: Termino...\n\n");

}
