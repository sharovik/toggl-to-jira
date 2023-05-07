# How to create Toggl.Track API token

1. Go to your account profile settings. It is important you select the Profile settings, but not a Workspace settings!
2. In the bottom of your account profile page you will find the block with API Token. Reveal your API token.
3. In the `.env` file please add the following attributes: 
    ``` 
    TOGGL_API_TOKEN={YOUR_TOKEN_PUT_HERE}
    TOGGL_API_URL=https://api.track.toggl.com
    TOGGL_DEFAULT_WORKSPACE_ID={YOUR_WORKSPACE_ID_PUT_HERE}
    ```
4. In the `TOGGL_API_TOKEN` please put your generated API token
5. Please go to [toggl.track timer screen](https://track.toggl.com) and click in the menu sidebar to `Reports` link to see the reports. In the URL, where you was redirected, you will find something like that `https://track.toggl.com/reports/summary/{workspace-id}`. Please take the `{workspace-id}` from that URL and put it in `TOGGL_DEFAULT_WORKSPACE_ID`.

Done.