# How to create Jira auth token

1. Go to [Account settings > Security > API Token](https://id.atlassian.com/manage-profile/security/api-tokens)
2. Generate the API token
3. In `.env` file please add the following attributes:
    ``` 
    JIRA_APP_TOKEN={YOUR_TOKEN}
    JIRA_EMAIL={YOUR_JIRA_EMAIL}
    JIRA_BASE_URL=https://your-company.atlassian.net
    ```
4. In the `JIRA_APP_TOKEN` please put your generated API token
5. In the `JIRA_EMAIL` please put your Jira account email, which will be used as login
6. In the `JIRA_BASE_URL` please put your company Jira URL

Voilà
