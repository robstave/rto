

# RTO App

An app that I used to really just kinda see how easy it was to use ChatGPT or Claude to build applications in one fell swoop.  Sometimes its hard to just make up something to code.  You need a problem first.

## Return to work app overview

### Problem

Our workplace mandates four days a week, but trying to square that up with everybodies personal situation and vacations is a bit complicated.  To address this, the metric is simply days a week in the office, regardless of situation.  So if you are coming in four days a week and not counting days that are missing to do vacation, illness, holidays, off site meetings and so on, then your actual target is still in the high 2.x to 3.0 right?  

I backload my vactions to the end of the year and several folks are not really sure how that math works.  So, this just tracks it for you.  You set a target value where you want to be and can work around that.

Its also handy for me to track vacation.  You can export the calendar to markdown. ( I use obsidian too ).

Its not really written for the world to use, its a project to play with ChatGPT in coding.
But if you think its applicable to your situation, go fork it!

### Features

- Golang
- Golang Templates / Echo
- SQLite
- Docker
- D3 burnup chart
- Export to Markdown
- Import ( needs work )

## Application
 
 This is set up to work with a quarter at a time.  If I have time, Ill fix this,
 but for now, its just q4 2024

![cal](/docs/cal1.png)

![cal](/docs/cal2.png)

You can set all the days in q4 as remote, in office or on vaction ( out ).
It counts everything for the quarter.

### Prefs 

When starting, the calendar is empty.  From the prefs you can fill out all the days in the quarter.
Pick your days.  You can also set your target date.

![cal](/docs/cal3.png)

There are some bulk adds where you can add a batch of days using json.  It works, but I cant really say I use it anymore.


## Deployment

There is a helper file that does building, docker, mocks and everything

` ./helper help ` to see commands

### Docker

```
docker compose build
# not docker-compose

docker compose up -d
docker compose down -d
```
 


### Testing

Generate the mocks with mockery.  See Helper File

## Working with Chatgpt

I run ` .helper concat` to get all the concatinated files into the /output dir
There is a file, file_list.yaml that sets all of this up

## Further Features

I could see if I can dial this up a notch and add
- Cloud Hosted
- Passwords/Accounts
- SSO login with OKTA

We will see.  Its nice to have an app with some meat on it to try that out with.

There is a password.  its aaa,aaa for now

