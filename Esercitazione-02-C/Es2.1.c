/*
 * main.c
 *
 *  Created on: Oct 23, 2017
 *      Author: s00007228894
 */

#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#include <semaphore.h>

#define NUMBER_OF_FILM 10
#define NUMBER_OF_PEOPLE 6
#define MAX_DOWNLOAD 3

typedef struct Film {
	int somma_voti;
	int votanti;
} Film;

struct Film arrayFilm[NUMBER_OF_FILM];
int arrayUtenti[NUMBER_OF_PEOPLE];
sem_t semUtenti[10];
int listaUtenti[10];
sem_t barriera;
sem_t download;
sem_t evento;
pthread_mutex_t mutex_threads[NUMBER_OF_FILM];
int completati=0, bestfilm, userdownloading=0;
float votomediofilm  = 0, bestvoto =0;

void *ThreadCode(void *t) {

	long tid;
	long result = 1;
	int i;
	int voto;
	tid = (int)t;
	arrayUtenti[tid]=0;
	for (i = 0; i < NUMBER_OF_FILM; i++) {
		pthread_mutex_lock(&mutex_threads[i]);
		voto = rand() % 10 + 1;
		arrayFilm[i].somma_voti += voto;
		arrayUtenti[tid] += voto;
		arrayFilm[i].votanti++;
		completati++;
		printf("Utente %ld. Film %d - Voto %d\n", tid, i, voto);
		pthread_mutex_unlock(&mutex_threads[i]);
	}
	arrayUtenti[tid] = arrayUtenti[tid]/NUMBER_OF_FILM;
	if(completati == NUMBER_OF_FILM*NUMBER_OF_PEOPLE){
		sem_post(&evento);
	}
	sem_wait(&barriera);
	sem_post(&barriera);
	if(userdownloading >= MAX_DOWNLOAD){
		listaUtenti[arrayUtenti[tid]]++;
		sem_wait(&semUtenti[arrayUtenti[tid]]);
	}
	userdownloading++;
	printf("Utente %ld. Inizio download\n", tid);
	sleep(rand() % 5 + 5);
	printf("Utente %ld. Termina download\n", tid);
	userdownloading--;
	if(userdownloading < MAX_DOWNLOAD){
		for(i=9;i>0;i--){
			if(listaUtenti[i]!=0){
				listaUtenti[i]--;
				sem_post(&semUtenti[i]);
				break;
			}
		}
	}
	printf("Utente %ld: guarda film %d\n", tid, bestfilm);
	sleep(rand() % 5 + 5);
	printf("Utente %ld: finisce di guardare film %d\n", tid, bestfilm);

	pthread_exit((void*) result);

}

int main (int argc, char *argv[]) {

	int i, rc, votomedioutente;
	long t;
	long result;
	float voto;
	pthread_t threads[NUMBER_OF_PEOPLE];
	sem_init (&barriera, 0, 0);
	sem_init (&evento, 0, 0);
	for (i = 0; i < NUMBER_OF_PEOPLE; i++) {
		sem_init (&semUtenti[i], 0, 0);
		listaUtenti[i]=0;
	}
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
	sem_wait(&evento);
	printf("\n\nMain: Risultati finali.\n");
	for (i = 0; i < NUMBER_OF_FILM; i++) {
		votomediofilm = arrayFilm[i].somma_voti/((double)arrayFilm[i].votanti);
		printf("Film %d: %f\n", i, votomediofilm);
		if(votomediofilm> bestvoto)
		{
			bestvoto = votomediofilm;
			bestfilm=i;
		}
	}
	for (i = 0; i < NUMBER_OF_PEOPLE; i++) {
			printf("Utente %d: priorita: %d\n", i, arrayUtenti[i]);
		}
	printf("\nMiglior film %d con voto %f\n\n",bestfilm,bestvoto);
	sem_post(&barriera);

	for (t = 0; t < NUMBER_OF_PEOPLE; t++) {

		rc = pthread_join(threads[t], (void *)&result);

		if (rc) {
			printf ("ERRORE: join thread %ld codice %d.\n", t, rc);
		} else {
			printf("Main: Finito Utente %ld.\n", t);
		}
	}
	printf("Main: Termino...\n\n");

}
