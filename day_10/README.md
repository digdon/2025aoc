# Day 10: Factory

Part 1 isn't terribly difficult. I set up a simple BFS for check subsequent button presses. What's interesting here is that because the buttons toggle the lights, we're essentially doing an XOR. What this means is that if we press a particular button a second time, we're undoing the first press. This means that when queuing up next presses for the BFS, we ultimately only need to press each button at most a single time.

My original solution took about 500ms to complete. This seemed a bit long to me, so while I was thinking about it, I realized that when I was processing button presses, the code was treating "button 1, then button 2" as being different than "button 2, then button 1". But that's not actually true - despite the different order, the end result is still the same. Seemed like there was an opportunity to change how I was tracking light states and button presses.

Gave it a go, tracking light states and button presses needed to get there "globally", rather than for each path as I was originally doing. And boom, down to 9ms.
