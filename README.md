# win_screenshots_sorter
A little tool I made for myself in GO to sort screenshots created by snipping tool.

When you take a screenshot using snipping tool, it saves it to \Pictures\Screenshots and names it "Screenshot YYYY-MM-DD HHMMSS.png".

I primarily created this so that I can get used to GO (and sort my screenshots out for once ;)), it's not perfect by any means...

It doesn't support languages other than English, the screenshots can't be opened in another program while they're being moved, etc.

Maybe i'll improve it at some point, does the job for me though.

# how it works
It first checks for the operating system, if it's not Windows, it exits.

It then grabs the users home directory and makes the path to the default screenshots folder to be used by ReadDir.

While looping through each entry:

it checks if it's a PNG file, then it gets the month and the year from the filename using time.Parse, then a folder is created based on the year and month, then the file is moved using os.Rename
