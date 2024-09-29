# Observations

It REALLY does css well

I usually use a Test Driven Design approach.  Is start with a function, and as I add functionality to a function, work it with the test.

With ChatGpt, this does not really work as all the functionality is just thrown out all at once.  I can ask for a test case, but really its more to validate the function after the fact.

so Retroactive Test Catchup ?   I dunno.  Its almost more like when you just right testcases after the fact to make SonarQube 




#  FEatures

Added testcases...but wondering why

Added persistance.  This took about an hour really.  Some issues with timestamp marshalling/unmarshalling.
Chat wants to take 09-12-24 and write it as an iso timestamp.  Fair enough. But it got down a rabbit hole with
unmarshalling overrides.  Changed it so the data.json is always written and timestamp.  Seems to have fixed it.

Added CSS
After lots of hints...chat gives a lotta hints...I caved in and went with suggestions.
Worked Greate

Calculate Average days and display
Pretty much nailed it

Bar Graph
Nailed it

Refactoring for better testing
Some functions were large and could be broken down to pure functions.  Thats me intervening.

Make days dynamic update
