# Artemis-Bot

### Roadmap
- [x] Fundamental bot structure
- [x] Redis integration + tests
- [x] Firebase + Firestore integration + tests
- [x] Add suggest and report commands
- [x] Implement first layer of cache flushing
- [x] Implement status command for changing submission statusses
- [x] Add config command
- [ ] Implement user notes
- [ ] Full support for polls and tickets
- [ ] User verification command
- [ ] Giveaways
- [ ] Server statistics
- [ ] Advanced user and server information
- [ ] Meeting (reminder) system
- [ ] Invite tracking + rewards

### Setup
This bot requires credentials to firebase (firebase-credentials.json) to operate. They must be in the same directory as the executable. The configuration file with be created automatically when the binary gets executed for the first time.

### Why encrypt/hash?
Well I must admit that it is not exactly as effective as people might think. The purpose is also not to stop the hoster from the bot from reading things like your notes. Therefore never put any sensitive data in there. Tracing back a note to a specific user for example would be a bit harder since the userid is hashed. However given the fact that one can easily retrieve all the users in all the guilds the bot is in, it becomes a relatively short bruteforce list. However for people that do not have access to the keys used to encrypt and hash the data it is rather hard to retrieve the actual information that is being secured. This is useful in case of a firebase-credentials leak or something alike. Now there might be features where users can use their own keys/passwords. However still keep in mind that this is a Discord bot. Every interaction or data you give it will be collected by Discord and this data will not be removed by Discord, ergo never save any confidential information using the bot or using Discord whatsoever.