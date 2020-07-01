### A web crawler in Go
- Concurrent fetching
- Follow TDD methodology


### TODO
- re-factor time sleep, with a time limit to wait for new request, if no new request and out channel has been drained then cancel in channel
    - how to check if a channel has been drained or not?
    - the closing decision depends on that
- Copy tests from fetchbot, re-write a few so i know what it's doing 