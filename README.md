# Gravemind

Controlling the flood of Twitch chat

## Architecture
 - React for the web frontend
 - A go app for the API
 - A go app to receive, parse, handle and send messages to & from TMI
 - DynamoDB as the primary data store

## Roadmap
 - Set up continuous integration & deployment via Travis CI and AWS CodeDeploy
 - Set up a RabbitMQ cluster
 - Make the most minimal go apps necessary to handle simple command & response
 - Make a minimal web UI to allow sign ups and adding commands
 - Add audits of command changes
 - Allow channel moderators to add & edit commands, with owner approval
 - Implement command regexes
 - Allow responding to commands with timeouts & bans
 - Add stats for how often commands are called & how long they take
 - Add scripting
 - Implement proper auto-scaling
 - Allow renaming the bot for a given channel
