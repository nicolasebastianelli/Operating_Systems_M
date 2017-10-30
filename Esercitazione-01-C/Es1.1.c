#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <math.h>
#define V_LENGTH 100
#define SUBV_LEN 5

int vector[V_LENGTH];

void *Max(void *t) //codice worker
{
	int num_threads=V_LENGTH/SUBV_LEN;
		if (V_LENGTH%SUBV_LEN)
			num_threads++;
	int i,max;
	long tid;
	tid= (int)t;
	int start = tid*SUBV_LEN;
	printf("Thread %ld è partito...\n",tid);
	max= vector[start];
	for (i=start+1;i<start+SUBV_LEN;i++)
		if(vector[i]>max)
			max=vector[i];
	printf("Thread %ld ha finito.Valore massimo= %ld\n",tid, max);
	pthread_exit((void*) max);
}

int main (int argc , char * argv[])
{
	int num_threads=V_LENGTH/SUBV_LEN;
	if (V_LENGTH%SUBV_LEN)
		num_threads++;
	pthread_t thread[num_threads];
	int rc,max=-1;
	long t,i;
	long status;
	printf("Main: Creazione vettore [");
	for(i=0;i<V_LENGTH;i++){
		vector[i]= rand()%100 +1;
		printf("%d, " ,vector[i]);
	}
	printf("]\n\n");

	for(t=0; t<num_threads; t++) {
		printf("Main: creazione thread %ld\n", t);
		rc=pthread_create(&thread[t], NULL,Max, (void *)t);
		if (rc) {
			printf("ERRORE: %d\n",rc);
			exit(-1); }
	}
	for(t=0; t<num_threads; t++) {
		rc=pthread_join(thread[t], (void *)&status);
		if (rc)
			printf("ERRORE join thread %ld	codice	%d\n", t,rc);
		else{
			printf("Finito thread %ld con ris. %ld\n",t,status);
			if(status>max)
				max=status;
		}
	}
	printf("Main: valore massimo trovato %d\n",max);
}
