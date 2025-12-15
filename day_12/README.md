# Day 12: Christmas Tree Farm

Saw this and immediately noped out. Went to Reddit to see how people were faring and it seemed that the puzzle was not quite what it at first seemed to be. Indeed, it seemed to be a bit of a joke on Eric's part. That is, despite what the description and sample data indicated, the actual puzzle data revealed something far simpler...

I started by looking to see if maybe we could reduce the number of manipulations we could do. ie, we're supposed to rotate and/or flip, but maybe, if the presents were easily mirrored, we could skip flipping. That didn't appear to be the case.

I mentioned this problem to my wife and she said, "well, the first thing I'd do is try to eliminate regions that we know can't possibily hold the presents. Calculate the area required, and if it's larger than the region area, you know they won't fit. Then you can do the complicated checks."

So I started with that.... I output a bunch of debug data while working out the area estimates, and noticed that when the region is too small, it's too small by a lot. Likewise, when the area is sufficient, it's sufficient for holding all of the presents side-by-side (ie, no interleaving).

Well, that's pretty interesting. I submitted my initial answer and presto, it worked. Sly, Eric, very sly. Dug deeper into Reddit and confirmed that everyone else experienced the same thing. Nice.
