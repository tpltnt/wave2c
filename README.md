wave2c
======

status: unfinished
todo:
* flushing array data to disk
* write tests

about
-----
wave2c a simple converter to dump [wave audio](http://en.wikipedia.org/wiki/WAV) into a C-array. This can be used to make audio recordings part of a
microconroller firmware. To do so the audio file has to be an 8bit mono
wavefile with a samplerate of 8KHz (for now).


usage
-----
* simply run: ```wave2c INPUTFILE.wav```

FAQ
---
* How do i build it?
  As long as you have [go](http://golang.org/) and [GNU make](http://www.gnu.org/software/make/) installed, simply type ```make```

* How can i convert an arbitary audio file so wave2c can eat it?
  Just use [SoX](http://sox.sourceforge.net/). Here is a commandline: ```sox INPUTFILE.wav --bits 8 --channels 1 --encoding unsigned-integer --rate 8k OUTPUTFILE.wav```

* How do i test the program?
  Simply run: ```make test```
