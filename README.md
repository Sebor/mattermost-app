# Mattermost App

## Local development

### Start MM server

- Run the server

  ```bash
  docker-compose -f docker-compose-mm-server.yml up -d
  ```

- Add entry in `/etc/hosts`

  ```txt
  127.0.0.1 mattermost
  ```

- Create admin user and team in web UI (<http://mattermost:8065>)
- Create a new channel(s) (or use existed) and get its id(s)

### Start bot

- Copy env config and place the channel ID(s), team name and API URL inside

  ```bash
  cp .env.bot-example .env
  ```

- Run the bot

  ```bash
  docker-compose up
  ```
