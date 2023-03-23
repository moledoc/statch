## Summary

trak is a program that tracks time.
**NOTE:** only one label tracking (trak) can be opened at any given time.

## Usage

trak ACTION [LABEL] [COMMENT]

## ACTION

start	

```
Starts new trak (time tracking).
By default label 'all' is used.
If any trak is opened at the time of starting a new trak, then the previous trak is closed.
After starting a new trak, the last 5 (including started) traks are printed.
```

end	

```
Ends the open trak and prints the last 5 traks.
```

show

```
Prints all logged traks.
```

summary	

```
Calculates monthly, weekly and daily summaries of traks, grouped by labels.
```

from %Y-%m-%dT%H:%M:%S
```
Starts new trak (time tracking) from given timestamp.
Recognized format is yyyy-mm-ddTHH:MM:SS. 
By default label 'all' is used.
If any trak is opened at the time of starting a new trak,
then the previous trak is closed.
After starting a new trak, the last 5 (including started) traks are printed.
```

to %Y-%m-%dT%H:%M:%S
```
Ends the open trak at given timetamp and prints the last 5 traks.
Recognized format is yyyy-mm-ddTHH:MM:SS.
```

## LABEL

By default label 'all' is used.
However, user can specify custom label after ACTION.
Only the first given label is used. Character '|' in label is not allowed.

## COMMENT

Every argument after label is considered to be a part of the comment for corresponding trak.
**NB!** to add comment, label must be provided! Character '|' in comment is not allowed.

## Author

Meelis Utt
