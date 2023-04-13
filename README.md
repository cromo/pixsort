## pixsort

Applying the traveling salesman problem to pixel art. This fork adds conversion of that path to SVG paths suitable for pen plotting.

Goal: Find the shortest path to visit all black pixels in an image.

Algorithm: Simulated annealing.

![Frog](http://i.imgur.com/2xiwTVE.gif)

![Peach](http://i.imgur.com/sCBhROn.gif)

## Usage

    go get github.com/fogleman/pixsort
    pixsort image.png

This will generate a file named `image.png.gif` with the result.

You can also pass in a `quality` parameter to make it try harder.

    pixsort image.png 28

The algorithm will run `2 ^ quality` iterations.
