#toggl-jira

It is a small application which can be executed manually or via crontab. This application will export your [toggl.track](https://toggl.com/track/) report for selected period to the work-logs of your Jira project. I wrote this app because I don't like the put my work-log time in Jira - it is annoying, and it can be automated.

## Requirements
As each application, my app also requires some things to be done, before the execution.
1. You must have an account in Jira, so you can create the Jira auth token. [Click here to see how.](documentation/create-jira-auth-token.md)
2. Create the toggl.track API token. [Click here to see how.](documentation/create-toggl-api-token.md)
3. You should have installed the sqlite on your system.
    You can use this command for ubuntu
    ```
    sudo apt-get install sqlite3 libsqlite3-dev
    ```
    Or by using brew
    ```
    brew install sqlite
    ```
    Or for centos
    ```
    sudo yum install sqlite
    ```
   
## How to use
1. Download the binary. I would recommend downloading the binary supported by your system from `bin/` to the `~/toggl-to-jira` folder. Of course, it is up to you.
2. Make sure all needed environment variables has been set in the `.env`.
3. Using Command line tool, please go to `~/toggl-to-jira`(or your path) and run the following command:
   ``` 
    ./{your-binary}
   ```
   Please note, that after run the app will try to fetch the time entries from toggl API and send it to the Jira API.
   Ideally, you should see something like that. In my case there was no time to put in Jira. 
   ![app-run](documentation/images/app-dry-run.png)
   
   If you see that, then it's done.
   


## Hint
I would recommend use this application with combination of toggl.track browser extension, [which you can find here](https://toggl.com/track/jira-time-tracking/). 



