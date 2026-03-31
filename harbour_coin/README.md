# Summary

- C: `pthread_mutex_t`, `pthread_cond_t`, `pthread_join`
- Go: `sync.Mutex`, `sync.Cond`, `sync.WaitGroup`

Main change:
- same idea, less manual thread handling
- no busy-wait loops
- queue close + broadcast used for clean shutdown
- `solution` guarded by its own mutex
