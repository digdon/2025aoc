# Day 10: Factory

Part 1 isn't terribly difficult. I set up a simple BFS for check subsequent button presses. What's interesting here is that because the buttons toggle the lights, we're essentially doing an XOR. What this means is that if we press a particular button a second time, we're undoing the first press. This means that when queuing up next presses for the BFS, we ultimately only need to press each button at most a single time.

My original solution took about 500ms to complete. This seemed a bit long to me, so while I was thinking about it, I realized that when I was processing button presses, the code was treating "button 1, then button 2" as being different than "button 2, then button 1". But that's not actually true - despite the different order, the end result is still the same. Seemed like there was an opportunity to change how I was tracking light states and button presses.

Gave it a go, tracking light states and button presses needed to get there "globally", rather than for each path as I was originally doing. And boom, down to 9ms.

_Some time later..._

Digging deeply into part 2 stuff. There was a lot of discussion around using linear algebra techniques to solve this - particularly Gaussian-Jordan elimination. I spent a _lot_ of time digging into this, working with various matrix forms (row echelon, reduced row echelon, and, finally, integer-only row echelon). However, that's only part of it. From here, you need to come up with a range of values for the free variables so that you can then brute-force through all the possibilities, solve all of the equations, and pick the one with the lowest answer.

At this point, however, I had nearly 500 lines of code and I was surely going to need hundreds more to work out reasonable ranges for the free variables. It was time to see if there was another eway....

Digging into Reddit, I found this [conversation](https://www.reddit.com/r/adventofcode/comments/1pk87hl/2025_day_10_part_2_bifurcate_your_way_to_victory/), and it changed everything.

It took me a long time to understand the algorithm, but I think I finally figured it out and managed to get a working solution. I ran the OP's code (Python) to compare and got the same answer.... eventually. Part of the understanding hurdle for me was remembering that rather than toggling lights, we're essentially counting the number of times a particular light is being hit by a sequence of button pushes. Once I understood the need of keeping track of how many times each light is touched in a set of pushes, it all fell into place.

There's lot of comments in the code that explain the complicated bits.

While doing all of this, I also tested a brute force method for part 1, similar to how the pattern generation code for part 2 works, but actually toggling light bits. A much simpler solution, requiring only half the lines of code as the BFS solution, and finishes in similar time. No brainer.
