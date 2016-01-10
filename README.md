# TyrantBot

An ambitious project to provide a free bot to every Twitch channel

## Architecture
 - A Meteor app for the web frontend
 - DynamoDB as the primary data store
 - A go app to receive incoming messages from TMI
 - A go app to parse incoming messages & handle them appropriately
 - A go app to send outbound messages to TMI
 - RabbitMQ to transfer messages between the go apps

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
 - Accept donations?
