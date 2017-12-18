# Operating System M

This repository contains laboratory exercise, exams and templates for the [Operating Systems M](http://lia.deis.unibo.it/Courses/som1718/) course of the Master's Degree in [Computer Engineering](http://corsi.unibo.it/ingegneriainformaticam/Pagine/default.aspx) at the [University of Bologna](http://www.unibo.it/it).

---

**To compile C programs:**

`gcc -D_REENTRANT -o fileName fileName.c -lpthread`

**To run C programs:**

`./fileName`

**To compile and run Go programs:**

`go run fileName.go`

:warning: **WARNING:**

Programs correctly compile also in OS X systems, but the semaphores do NOT behave the way you would expect because they are not initialized by the deprecated `sem_init()`. If you want to use OS X systems [here](https://heldercorreia.com/semaphores-in-mac-os-x-fd7a7418e13b) you can find a solution.
