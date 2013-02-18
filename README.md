wave2c
======

status: unfinished

about
-----
wave2c a simple converter to dump wave audio into a C-array. This can be used
to make audio recordings part of microconroller firmware. To do so the audio
file has to be an 8bit mono wavefile with a samplerate of 8KHz (for now).


usage
-----
* simply run: wave2c inputfile.wav

FAQ
---
* How can i convert an arbitary audio file so wave2c can eat it?
  Just use sox. Here is a commandline: sox INPUTFILE.wav --bits 8 --channels 1 --encoding unsigned-integer --rate 8k OUTPUT.wav
