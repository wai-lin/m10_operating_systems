#include <assert.h>
#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <pthread.h>
#include <unistd.h>
#include <inttypes.h>

// NOTE: The minners were stucking because the manager exit when one thread
// found the solution, but the other thread were still waiting for the manager.
// 
// Fixe:
// Here I add a `closed` flag so that manager can signal the miners to exit.
// I also added `head` and `tail` to make the queue run in circular way.
// Also added condition variables to avoid busy waiting.

#define QUEUE_SIZE 10
uint64_t solution = 0;

void *manager(void*);
void *miner(void*);

typedef struct {
   uint64_t items[QUEUE_SIZE];
   size_t length;
   size_t head;
   size_t tail;
   int closed;
   pthread_mutex_t mutex;
   pthread_cond_t not_empty;
   pthread_cond_t not_full;
} queue_t;

pthread_mutex_t solution_mutex = PTHREAD_MUTEX_INITIALIZER;

void queue_init(queue_t *this) {
   this->length = 0;
   this->head = 0;
   this->tail = 0;
   this->closed = 0;
   pthread_mutex_init(&this->mutex, NULL);
   pthread_cond_init(&this->not_empty, NULL);
   pthread_cond_init(&this->not_full, NULL);
}

void queue_close(queue_t *this) {
   pthread_mutex_lock(&this->mutex);
   this->closed = 1;
   pthread_cond_broadcast(&this->not_empty);
   pthread_cond_broadcast(&this->not_full);
   pthread_mutex_unlock(&this->mutex);
}

int queue_add(queue_t *this, uint64_t item) {
   pthread_mutex_lock(&this->mutex);
   while(this->length == QUEUE_SIZE && !this->closed) {
        pthread_cond_wait(&this->not_full, &this->mutex);
   }
   if (this->closed) {
        pthread_mutex_unlock(&this->mutex);
        return 0;
   }

   this->items[this->tail] = item;
   this->tail = (this->tail + 1) % QUEUE_SIZE;
   this->length++;
   pthread_cond_signal(&this->not_empty);
   pthread_mutex_unlock(&this->mutex);
   return 1;
}

uint64_t queue_pop(queue_t *this, uint64_t *item) {
   pthread_mutex_lock(&this->mutex);
   while (this->length == 0 && !this->closed) {
        pthread_cond_wait(&this->not_empty, &this->mutex);
   }
   if (this->length == 0 && this->closed) {
        pthread_mutex_unlock(&this->mutex);
        return 0;
   }

   *item = this->items[this->head];
   this->head = (this->head + 1) % QUEUE_SIZE;
   this->length--;
   pthread_cond_signal(&this->not_full);
   pthread_mutex_unlock(&this->mutex);
   return 1;
}

uint64_t hash(uint64_t x) {
    x = (x ^ (x >> 30)) * UINT64_C(0xbf58476d1ce4e5b9);
    x = (x ^ (x >> 27)) * UINT64_C(0x94d049bb133111eb);
    x = x ^ (x >> 31);
    return x;
}

const int N_MINERS = 100;
const uint64_t SLICE_SIZE = 1000000;
const uint64_t LOWER_BITS_MASK = 0xfffffff;

uint64_t seed;
queue_t queue;

int main() {
   pthread_t thread_manager;
   pthread_t threads_miners[N_MINERS];

   srandom(time(NULL));
   seed = random();

   queue_init(&queue);

   assert(!pthread_create(&thread_manager, NULL, &manager, NULL));
   for (int i = 0; i < N_MINERS; i++) {
      assert(!pthread_create(&threads_miners[i], NULL, &miner, NULL));
   }
   assert(!pthread_join(thread_manager, NULL));
   int sum = 0;
   for (int i = 0; i < N_MINERS; i++) {
      assert(!pthread_join(threads_miners[i], NULL));
   }
   printf("%d\n", sum);

   exit(0);
}

void *manager(void* _param) {
   uint64_t slice_base = SLICE_SIZE;

   while (1) {
        pthread_mutex_lock(&solution_mutex);
        int done = (solution != 0);
        pthread_mutex_unlock(&solution_mutex);
        if (done) {
            break;
        }
        if (!queue_add(&queue, slice_base)) {
            break;
        }

        printf("sent %" PRIu64 "\n", slice_base);
        slice_base += SLICE_SIZE;
   }
   return NULL;
}

void *miner(void* _param) {
   while (1) {
      pthread_mutex_lock(&solution_mutex);
      if (solution != 0) {
         pthread_mutex_unlock(&solution_mutex);
         return NULL;
      }
      pthread_mutex_unlock(&solution_mutex);

      uint64_t slice_base;
      if (!queue_pop(&queue, &slice_base)) {
         return NULL;
      }

      for (uint64_t i = slice_base; i < slice_base + SLICE_SIZE; i++) {
         uint64_t hashed = i ^ seed;
         for (int j = 0; j < 10; j++) {
            hashed = hash(hashed);
         }

         if ((hashed & LOWER_BITS_MASK) == 0) {
            pthread_mutex_lock(&solution_mutex);
            if (solution == 0) {
               solution = i;
               queue_close(&queue);
               printf("miner found solution %" PRIu64 "\n", i);
            }
            pthread_mutex_unlock(&solution_mutex);
            return NULL;
         }
      }
   }
}